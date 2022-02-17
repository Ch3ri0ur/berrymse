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

	//CMD Prints
	"fmt"

	//Log
	"log"

	// WEB STUFF			################################################################

	//Ip and Net stuff https://pkg.go.dev/net https://pkg.go.dev/net/http
	"net"
	"net/http"

	//For static website stuff in binary executable github.com/markbates/pkger
	"github.com/markbates/pkger"

	// OWN STUFF ##########################################################################

	//FLAG & Configuration
	Config "github.com/ch3ri0ur/berrymse/src/config"

	Source "github.com/ch3ri0ur/berrymse/src/source"

	Hub "github.com/ch3ri0ur/berrymse/src/hub"
)

//Init methode
//Defining Flags and Default values
func init() {
	Config.FlagInit()
}

func main() {

	var configuration = Config.SetupConfigFlags()

	// Parse host:port into host and port
	host, port, err := net.SplitHostPort(configuration.Server.URL)
	if nil != err {
		log.Fatal(err)
	}
	fmt.Println("Starting Hub")
	// One-to-many hub broadcasts NAL units as MP4 fragments to clients
	hub := Hub.NewHub()
	go hub.Run(configuration)

	fmt.Println("Starting Source")

	// Open source
	src := Source.NewSource(hub, configuration)
	go src.Run()

	fmt.Println("Starting Server")

	http.HandleFunc("/video_websocket", func(w http.ResponseWriter, r *http.Request) {
		Hub.ServeWs(hub, w, r)
	})
	http.Handle("/", http.FileServer(pkger.Dir("/cmd/berryMSE/web/static")))

	// Start server
	fmt.Printf("Listening on http://%v:%v\n", host, port)
	log.Fatal(http.ListenAndServe(configuration.Server.URL, nil))
}
