# Run

```
sudo ./berrymse
```

This runs a server that provides a demo website. Without a config file or flags the port is `:2020` and the url is `localhost`.

The Raspberry Pi camera needs to be activated first on the Pi and on Bullseye (RPi OS 11) the old camera driver needs to be reactivated. This can be done in the ``raspi-config`` under "Interfacing options".

## Flags

By using flags some parameter can be changed, as shown below for website url and input device:

``` bash
sudo ./berrymse -l <raspberry pi ip address>:<port> -d /dev/video<X>
```

More flags for ./berrymse:

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


## Config File

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


## Register Service

To register the executable as an autostart service:

Ensure that the paths in `berrymse.service` are correct. 

Default is to clone the directory into the `/home/pi/` directory

```
cd for_autostart
sudo ./register.sh
```

Restart to test the service.
