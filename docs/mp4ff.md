ffmpeg -f video4linux2 -input_format h264 -video_size 1280x720 -framerate 30 -i /dev/video0 ~/out1.mp4

mp4ff cmd
info go run main.go ~/out1.mp4

v4l2-ctl --list-formats -d /dev/video3

-> h264 (libx264)


ffmpeg uses own codec convertion so all nulheader are the same just frame rate (often much lower) and sometimes profile gets adjusted