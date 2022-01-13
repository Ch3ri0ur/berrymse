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
	//Byte Arrays
	"bytes"
	//CMD Prints
	"fmt"
	//Log
	"log"

	//Time stuff
	"time"

	//OS Function to access Path and files https://pkg.go.dev/os
	"os"
	//Error handling https://pkg.go.dev/errors
	"errors"
	//Easy to use Path and File utilities https://pkg.go.dev/path/filepath
	"path/filepath"

	//Ip and Net stuff https://pkg.go.dev/net https://pkg.go.dev/net/http
	"net"
	"net/http"

	//Websocket https://github.com/gorilla/websocket
	"github.com/gorilla/websocket"

	//For static website stuff in binary github.com/markbates/pkger
	"github.com/markbates/pkger"

	//Old CMD Flag Handle
	//"flag"
	//Improved Flag Handle compatible with Viper Configuration Manger https://github.com/spf13/pflag
	//Default Values set for the Flags are used to set the Configuration
	//Imported flag to replace the old "flag" handle

	flag "github.com/spf13/pflag"

	//Configuration Manger https://github.com/spf13/viper
	//This Configuration Manager allows the use of a Configfile and Flags to set the Configuration
	//Priority is UsedFlag>Configfile>DefaultFlag!
	"github.com/spf13/viper"

	//v4l2 Lib to Access Camera
	"github.com/ch3ri0ur/go-v4l2"
)

//Command line flag parameters
//Tmp Flag Save Location DO NOT USE IN CODE!!
//USE configuration (Configurations) to access the config values
var flagServerURL string
var flagCameraFD string
var flagConfig string

//Init methode
//Defining Flags and Default values
//	PFlagtype(&Variable, ConfigID,
//		FlagID,
//		Defaultvalue,
//		Info Text,
//	)
func init() {
	//Config Path Flag
	//No config Name, only needed to load selected config file
	flag.StringVarP(&flagConfig, "",
		"c",
		"config.yml",
		"Use config Path/Name.yml"+
			"\nDefault Path is current directory!",
	)

	//Flag to selected an URL
	flag.StringVarP(&flagServerURL, "Server.URL",
		"l",
		"localhost:2020",
		"listen on host:port",
	)

	//Flag to change the Device input file
	flag.StringVarP(&flagCameraFD, "Camera.SourceFD",
		"d",
		"/dev/video0",
		"Use camera /dev/videoX",
	)
}

//All Configurations Stored in this. Look config.go for structure
var configuration Configurations

//Reads Flags and Configfile to set and overwrite the Config
func setupConfigFlags() {

	//Get all Flags and Parse them in Variables
	flag.Parse()
	//Bind Flags to Config
	viper.BindPFlags(flag.CommandLine)
	// Not bound variables
	//viper.SetDefault("Camera.FD", "/dev/video1")

	//Checks for and Loads Configfile
	LoadConfigs()

	//Loads the Config into the Struct for easier usage
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v\n", err)
	}

	//Showcase of usage of Config and some test code
	//fmt.Printf("Camera FD Configuration: %s \n", configuration.Camera.SourceFD)
	//fmt.Printf("Server URL Configuration: %s \n", configuration.Server.URL)
}

//Checks if Configfile exists and read it
//When flagConfig only contains a Path it will use the default config name "config.yml"
func LoadConfigs() {
	//Checks if the Configfile exists. uses the flag Value that contains a given or the default ConfigfilePath. Skipping the Loading of Configfile if not exists
	if _, err := os.Stat(flagConfig); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Error: Configfile not found: %s \n", flagConfig)
		fmt.Printf("Skip loading Configfile! Using default Settings!\n")
		return
	}

	//Path and Filename with extention gets split up
	dir, file := filepath.Split(flagConfig)

	//If the Path is empty (current directory) a "." needs to be used
	if dir == "" {
		dir = "."
	}

	// If only a Path was given use the default Configfilename "config.yml"
	if file == "" {
		file = "config.yml"
		fmt.Printf("Warning. Flag -c contained a Filepath without filename, %s \n"+
			"Using default \"config.yml\" as config file.\n", flagConfig)
	}
	//Extract the Filename by removing the Extention
	fileName := file[:len(file)-len(filepath.Ext(file))]

	//Extract the Extention
	fileExtention := filepath.Ext(file)
	//The "." of the extention needs to be removed
	if fileExtention != "" {
		fileExtention = fileExtention[1:]
	}

	//Check if Extention of the is Supported. Skip Loading the Configfile, when not supported!
	if fileExtention != "yml" {
		fmt.Printf("Error. Flag -c contained an File with wrong file extention, %s \n", flagConfig)
		fmt.Printf("Skip loading Configfile! Using default Settings!\n")
		return
	}

	// Set the path, name and extention to look for the configurations file
	viper.AddConfigPath(dir)
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileExtention)

	//Read the Configuration File
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}

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
	dev, err := v4l2.Open(flagCameraFD)
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

	setupConfigFlags()

	// Parse host:port into host and port
	host, port, err := net.SplitHostPort(flagServerURL)
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
	log.Fatal(http.ListenAndServe(flagServerURL, nil))

}
