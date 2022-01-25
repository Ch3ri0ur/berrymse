package source

import (
	//CMD Prints
	"fmt"

	//Log
	"log"

	// OWN STUFF ##########################################################################

	//Configuration
	config "github.com/ch3ri0ur/berrymse/src/config"

	hub "github.com/ch3ri0ur/berrymse/src/hub"

	//v4l2 go Lib to Access Camera
	"github.com/ch3ri0ur/go-v4l2"
)


type Source struct {
	device *v4l2.Device
	hub    *hub.Hub
}

func NewSource(h *hub.Hub, configuration config.Configurations) *Source {
	// Open device
	dev, err := v4l2.Open(configuration.Camera.SourceFD)
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
	if err := dev.SetBitrate(int32(configuration.Camera.Bitrate)); nil != err {
		log.Fatal(err)
	}
	fmt.Println("after Set Bitrate")
	return &Source{
		device: dev,
		hub:    h,
	}
}

func (s *Source) Run() {
	// Start stream
	if err := s.device.Start(); nil != err {
		log.Fatal(err)
	}
	defer s.device.Stop()

	for {
		select {
		case b := <-s.device.C:
			s.hub.Nals <- b.Data
			b.Release()
		}
	}
}