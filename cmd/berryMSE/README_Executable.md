# 🍓 BerryMSE

Simple low-latency live video streaming from a Raspberry Pi&trade; via the [Media Source Extensions API](//developer.mozilla.org/en-US/docs/Web/API/Media_Source_Extensions_API).

Note: As of March 2020, Safari on iOS devices still does not support this API (excluding iOS 13 on iPad devices, which do support the API).

## Overview

H.264 Network Abstraction Layer (NAL) units are read from `/dev/video0`, a
Video4Linux2 compatible camera interface. Each unit corresponds to one frame.
Frames are packaged into MPEG-4 ISO BMFF (ISO/IEC 14496-12) compliant
fragments and sent via a websocket to the browser client. The client appends
each received buffer to the media source for playback.

## Run

```
sudo ./berrymse
```
This runs a server that provides a demo website. Without a config file port 2020 is standard. To open it visit `localhost`.

The server provides a webpage (`index.html`), a websocket stream of the camera (`/video_websocket`) and the javascript (`/msevideo.js`) to run it. For more information on how to integrate this into another project please see the chapter below.

For HTTP port 80 use sudo and specify port (0.0.0.0:80). 

``` bash
sudo ./berrymse -l <raspberry pi ip address>:<port> -d /dev/video<X>
```

Usage of ./berrymse:

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

or configure the executable by placing a `config.yml`  with the following content in the same folder as the executable. The possible parameters can be seen under `berrymse -h`.

!!! info inline end
    If these configurations don't work/match your camera this can freeze the camera stack. e.g. using resolutions above 1920 times 1080 created crashes.

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

Run with sudo and visit website under ```localhost```.


## Register Service
To register the executable as an autostart service:

Ensure that the paths in `berrymse.service` are correct. 

Default is to clone the directory into the `/home/pi/` directory

```
cd for_autostart
sudo ./register.sh
```

Restart to test the service.
