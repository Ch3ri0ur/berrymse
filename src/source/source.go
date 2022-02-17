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
	v4l2 "github.com/ch3ri0ur/go-v4l2"
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
		configuration.Camera.Width,
		configuration.Camera.Height,
		v4l2.V4L2_PIX_FMT_H264,
	); nil != err {
		log.Fatal(err)
	}

	// Set bitrate
	if configuration.Camera.Bitrate > 0 {
		fmt.Println("Set Bitrate")
		if err := dev.SetBitrate(int32(configuration.Camera.Bitrate)); nil != err {
			log.Fatal(err)
		}
	}

	if configuration.Camera.Rotation >= 0 {
		fmt.Println("Set Rotation")
		if err := dev.SetRotation(int32(configuration.Camera.Rotation)); nil != err {
			log.Fatal(err)
		}
	}

	// // Custom Configuration possible with
	// dev.SetCustomUserControl(id uint32, value int32)
	// dev.SetCustomCodecControl(id uint32, value int32)
	// // Check device with "v4l2-ctl --all -d /dev/videoX" for IDs
	// // User stuff = 0x00980000 - 0x0098ffff
	// // Codec stuff = 0x00990000 - 0x0099ffff
	// // e.g. user control vertical flip = 0x00980915

	return &Source{
		device: dev,
		hub:    h,
	}
}

func (s *Source) Run() {
	// Start stream
	fmt.Println("Soruce Run")
	if err := s.device.Start(); nil != err {
		log.Fatal(err)
		fmt.Println(err)
	}
	defer s.device.Stop()
	fmt.Println("Source started")
	for {
		select {
		case b := <-s.device.C:
			s.hub.Nals <- b.Data
			b.Release()
		}
	}
}
