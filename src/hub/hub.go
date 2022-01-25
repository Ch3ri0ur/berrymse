package hub

import (
	//Byte Arrays
	"bytes"

	//CMD Prints
	"fmt"

	//Log
	"log"

	//Time stuff
	"time"

	// WEB STUFF			################################################################

	//Ip and Net stuff https://pkg.go.dev/net https://pkg.go.dev/net/http
	"net/http"

	//Websocket https://github.com/gorilla/websocket
	"github.com/gorilla/websocket"


	// OWN STUFF ##########################################################################

	//Configuration
	config "github.com/ch3ri0ur/berrymse/src/config"

	bmff "github.com/ch3ri0ur/berrymse/src/bmff"
)

const (
	nalTypeNonIDRCodedSlice = 1
	nalTypeIDRCodedSlice    = 5

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

// client structure
type Client struct {
	hub     *Hub
	conn    *websocket.Conn // Websocket connection
	frags   chan []byte     // Buffered channel of outbound MP4 fragments
	n       int             // Frame number
	haveIDR bool            // Received i-frame?
}

// hub maintains a set of active clients and broadcasts video to clients
type Hub struct {
	clients    map[*Client]bool // registered clients
	Nals       chan []byte      // NAL units from camera source
	register   chan *Client     // register requests from clients
	unregister chan *Client     // unregister requests from clients
}

// newHub instantiates a new hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Nals:       make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// run processes register and unregister requests, and nal units
func (h *Hub) Run(configuration config.Configurations) {
	for {
		select {
		// Register request
		case c := <-h.register:
			h.clients[c] = true

			var frag bytes.Buffer
			bmff.WriteFTYP(&frag)
			bmff.WriteMOOV(&frag, uint16(configuration.Camera.Width), uint16(configuration.Camera.Height))
			c.frags <- frag.Bytes()

		// Unregister request
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.frags)
			}

		// New NAL from source
		case nal := <-h.Nals:
			nal = bytes.TrimPrefix(nal, []byte{0, 0, 0, 1})
			if len(nal) == 0 {
				break
			}
			nalType := (nal[0] & 0x1F)

			//TODO OPTIMISING for multiple clients by extrakting msgbuilding out the loop of every client problematic is the client frame indexs that is different for each client

			//Send new frag to all Clients that are registered.
			for c := range h.clients {
				//Buffer to fill with Header, info, frag
				var frag bytes.Buffer

				//Check nalType
				//- When it is a "nalTypeIDRCodedSlice" (frag with IDR) it will set the flag of the client for having received a IDR (client.haveIDR), than fallthrough into the case nalTypeNonIDRCodedSlice" and send the frag
				//- When it is a "nalTypeNonIDRCodedSlice" check if the client has ever received a "nalTypeIDRCodedSlice" (client.haveIDR), if yes than send the frag, if not just skip
				//This will cause that the client will receive its first frag when it is a frag with idr. After that it will all send all slices to the client
				switch nalType {

				//frag contains IDR. Initial frag for the client
				case nalTypeIDRCodedSlice:
					//Set Flag for client has received a
					c.haveIDR = true
					//Jump into the case "nalTypeNonIDRCodedSlice" to send the Data to the client
					fallthrough

				//frag contains no IDR
				case nalTypeNonIDRCodedSlice:
					if c.haveIDR {
						bmff.WriteMOOF(&frag, c.n, nal)
						bmff.WriteMDAT(&frag, nal)
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

				//If naltype doesnt fit just do nothing
				default:
					// noop
				}
			}
		}
	}
}

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

// Handle websocket client connections
func ServeWs(h *Hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade websocket connection from HTTP to TCP
	conn, err := upgrader.Upgrade(w, r, nil)
	if nil != err {
		log.Println(err)
		return
	}

	// Instantiate client
	c := &Client{hub: h, conn: conn, frags: make(chan []byte, 30), n: 1}
	c.hub.register <- c

	// Go routine writes new MP4 fragment to client websocket
	go func(c *Client) {
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