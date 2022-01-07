// BerryMSE: Low-latency live video via Media Source Extensions (MSE)
// Copyright (C) 2020 Chris Hiszpanski
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/markbates/pkger"
	"github.com/thinkski/go-v4l2"
)

// Command line flag parameters
var (
	flagListen string
	devListen  string
)

func init() {
	flag.StringVar(
		&flagListen, "l", "localhost:2020", "listen on host:port",
	)

	flag.StringVar(
		&devListen, "d", "/dev/video0", "/dev/video0",
	)
}

const (
	nalTypeNonIDRCodedSlice = 1
	nalTypeIDRCodedSlice    = 5

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

// Websocket parameters
var upgrader = websocket.Upgrader{
	// Tune read buffers for short acknowledgement messages
	ReadBufferSize: 256,

	// Tune write buffers to comfortably fit most all B and P frames.
	WriteBufferSize: 8192,

	// Allow any origin for demo purposes
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// client structure
type client struct {
	hub     *hub
	conn    *websocket.Conn // Websocket connection
	frags   chan []byte     // Buffered channel of outbound MP4 fragments
	n       int             // Frame number
	haveIDR bool            // Received i-frame?
}

// hub maintains a set of active clients and broadcasts video to clients
type hub struct {
	clients    map[*client]bool // registered clients
	nals       chan []byte      // NAL units from camera source
	register   chan *client     // register requests from clients
	unregister chan *client     // unregister requests from clients
}

// newHub instantiates a new hub
func newHub() *hub {
	return &hub{
		clients:    make(map[*client]bool),
		nals:       make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

// run processes register and unregister requests, and nal units
func (h *hub) run() {
	for {
		select {
		// Register request
		case c := <-h.register:
			h.clients[c] = true

			var frag bytes.Buffer
			writeFTYP(&frag)
			writeMOOV(&frag, 1280, 720)
			c.frags <- frag.Bytes()

		// Unregister request
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.frags)
			}

		// New NAL from source
		case nal := <-h.nals:
			nal = bytes.TrimPrefix(nal, []byte{0, 0, 0, 1})
			if len(nal) == 0 {
				break
			}
			nalType := (nal[0] & 0x1F)

			for c := range h.clients {
				var frag bytes.Buffer

				switch nalType {
				case nalTypeIDRCodedSlice:
					c.haveIDR = true
					fallthrough
				case nalTypeNonIDRCodedSlice:
					if c.haveIDR {
						writeMOOF(&frag, c.n, nal)
						writeMDAT(&frag, nal)
						c.n++

						select {
						// Write MP4 fragment
						case c.frags <- frag.Bytes():

						// Buffered channel full. Drop client.
						default:
							close(c.frags)
							delete(h.clients, c)
						}
					}
				default:
					// noop
				}
			}
		}
	}
}

type source struct {
	device *v4l2.Device
	hub    *hub
}

func newSource(h *hub) *source {
	// Open device
	dev, err := v4l2.Open(devListen)
	if nil != err {
		log.Fatal(err)
	}

	// Set pixel format
	if err := dev.SetPixelFormat(
		1280,
		720,
		v4l2.V4L2_PIX_FMT_H264,
	); nil != err {
		log.Fatal(err)
	}

	fmt.Println("before Set Bitrate")

	// Set bitrate
	if err := dev.SetBitrate(1500000); nil != err {
		log.Fatal(err)
	}
	fmt.Println("after Set Bitrate")
	return &source{
		device: dev,
		hub:    h,
	}
}

func (s *source) run() {
	// Start stream
	if err := s.device.Start(); nil != err {
		log.Fatal(err)
	}
	defer s.device.Stop()

	for {
		select {
		case b := <-s.device.C:
			s.hub.nals <- b.Data
			b.Release()
		}
	}
}

// Handle websocket client connections
func serveWs(h *hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade websocket connection from HTTP to TCP
	conn, err := upgrader.Upgrade(w, r, nil)
	if nil != err {
		log.Println(err)
		return
	}

	// Instantiate client
	c := &client{hub: h, conn: conn, frags: make(chan []byte, 30), n: 1}
	c.hub.register <- c

	// Go routine writes new MP4 fragment to client websocket
	go func(c *client) {
		defer func() {
			fmt.Println("closing socket")
			c.conn.Close()
		}()
		fmt.Println("starting socket")

		for {
			select {
			case frag, ok := <-c.frags:
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					// Hub closed the channel
					c.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				// Write next segment
				nw, err := c.conn.NextWriter(websocket.BinaryMessage)
				if nil != err {
					return
				}
				nw.Write(frag)

				// Close writer
				if err := nw.Close(); nil != err {
					return
				}
			}
		}
	}(c)
}

func main() {
	flag.Parse()

	// Parse host:port into host and port
	host, port, err := net.SplitHostPort(flagListen)
	if nil != err {
		log.Fatal(err)
	}
	fmt.Println("Starting Hub")
	// One-to-many hub broadcasts NAL units as MP4 fragments to clients
	hub := newHub()
	go hub.run()

	fmt.Println("Starting Source")

	// Open source
	src := newSource(hub)
	go src.run()

	fmt.Println("Starting Server")

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.Handle("/", http.FileServer(pkger.Dir("/web/static")))

	// Start server
	fmt.Printf("Listening on http://%v:%v\n", host, port)
	log.Fatal(http.ListenAndServe(flagListen, nil))
}
