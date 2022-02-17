# Usage Instruction



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