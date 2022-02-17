package berryMSE

import (
	//CMD Prints
	"fmt"

	//Log
	"log"

	// WEB STUFF			################################################################

	//Ip and Net stuff https://pkg.go.dev/net https://pkg.go.dev/net/http
	"net"
	"net/http"

	// OWN STUFF ##########################################################################

	//FLAG & Configuration
	Config "github.com/ch3ri0ur/berrymse/src/config"

	Source "github.com/ch3ri0ur/berrymse/src/source"

	Hub "github.com/ch3ri0ur/berrymse/src/hub"
)

func BerryMSE(configuration Config.Configurations, WebPageFiles http.FileSystem) {
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
	src := Source.NewSource(hub, configuration.Camera)
	go src.Run()

	fmt.Println("Starting Server")

	http.HandleFunc("/"+configuration.Server.WebSocket, func(w http.ResponseWriter, r *http.Request) {
		Hub.ServeWs(hub, w, r)
	})
	http.Handle("/", http.FileServer(WebPageFiles))

	// Start server
	fmt.Printf("Listening on http://%v:%v\n", host, port)
	log.Fatal(http.ListenAndServe(configuration.Server.URL, nil))
}
