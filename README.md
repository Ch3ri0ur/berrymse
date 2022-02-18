# 🍓 BerryMSE

Simple low-latency live video streaming from a Raspberry Pi&trade; via the [Media Source Extensions API](//developer.mozilla.org/en-US/docs/Web/API/Media_Source_Extensions_API).

Note: As of March 2020, Safari on iOS devices still does not support this API (excluding iOS 13 on iPad devices, which do support the API).

## Overview

H.264 Network Abstraction Layer (NAL) units are read from `/dev/video0`, a
Video4Linux2 compatible camera interface. Each unit corresponds to one frame.
Frames are packaged into MPEG-4 ISO BMFF (ISO/IEC 14496-12) compliant
fragments and sent via a websocket to the browser client. The client appends
each received buffer to the media source for playback.

## Demo

The Demo executable can be downloaded from the release page and run on a Raspberry Pi 32-bit (Buster) with a Raspberry Pi camera as `/dev/video0`.

The demo files can be found in ``/cmd/berrymse`` with a build and usage instruction.

To run, copy the appropriate `berrymse` executable to the Raspberry Pi and run:

	./berrymse -l <raspberry pi ip address>:2020 -d /dev/video<X>

For example:

    ./berrymse -l 192.168.2.1:2020 -d /dev/video0

The Raspberry Pi Zero uses the `armv6l/berrymse` executable. Other models use
the `armv7l/berrymse` executable.

The webpage will show a live video stream with approximately 200ms of latency.
The browser will buffer frames, providing a lookback window.

## Settings

The currently implemented configurations are:

-  H264 bitrate:
    Changes the bitrate of the video and the encoder will now try to only produce a [H264 Stream](Theory/Video/h264.md).
- Video height resolution:
    Changes the video height resolution.
- Video width resolution:
    Changes the video width resolution.
- Video rotation:
    Changes the video rotation resolution. It can only be changed in 90 degree steps and rotates the picture clockwise.
- Video source:
    Changes the source Device Node of the video.
- Server URL address:
    Changes the URL address were the video and website gets published to.
- Server websocket name:
    Changes the websocket name were the video packages can get received. This will break the demo page!

### Flags

Flags for ./berrymse:

    -c, -- string                  Use config Path/Name.yml
                                    Default Path is current directory! (default "config.yml")
    -b, --Camera.Bitrate int       Bitrate in bit/s!
                                    Only supported for RPI Camera
                                    Other Cameras need to use -1 (default 1500000)
    -h, --Camera.Height int        Height Resolution (default 720)
    -r, --Camera.Rotation int      Rotation in 90degree Step
                                    Only supported for RPI Camera
                                    Other Cameras need to use -1
    -d, --Camera.SourceFD string   Use camera /dev/videoX (default "/dev/video0")
    -w, --Camera.Width int         Width Resolution (default 1280)
    -l, --Server.URL string        listen on host:port (default "localhost:2020")


### Config File

Default config file name is `config.yml` and default path for it is local directory.

If these configurations don't work/match your camera this can freeze the camera stack. e.g. using resolutions above 1920 times 1080 created crashes.

USB cameras don't support the advanced settings rotation and bitrate and need a -1 as parameter.

``` yaml title="config.yml"
camera:
  sourceFD: "/dev/video0"
  width: 1280
  height: 720
  bitrate: 1500000
  rotation: 0

server:
  url: "0.0.0.0:80"
```

## Project folder structure

- berrymse/                 : Project folder
    - cmd/                  : CMD Applications
        - berryMSE          : Demo
        - berryMSEmulti     : Demo for multiple streams
    - src/                  : Contains library
        - berryMSE          : Main Class
        - ...
    - configs/              : ConfigTemplates

