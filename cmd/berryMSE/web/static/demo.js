window.onload = function() {

  let videoElement = document.querySelector('video');
  
  // Check that browser supports Media Source Extensions API
  if (window.MediaSource) {
    let mediaSource = new MediaSource();
    videoElement.loop = false;
    videoElement.src = URL.createObjectURL(mediaSource);
    mediaSource.addEventListener('sourceopen', sourceOpen);
    videoElement.onpause = function() {
      console.log("buffered:", videoElement.buffered);
      let buffered = videoElement.buffered;
      videoElement.currentTime = buffered.end(buffered.length-1) - 0.2;
      videoElement.play();
    }
  } else {
    console.log("Media Source Extensions API is NOT supported");
  }
  
  function sourceOpen(e) {
    URL.revokeObjectURL(videoElement.src);

    let mediaSource = e.target;

    // remote pushes media segments via websocket
    ws = new WebSocket("ws://" + location.hostname + (location.port ? ":"+location.port : "" ) + "/websocket");
    ws.binaryType = "arraybuffer";

    // The six hexadecimal digit suffix after avc1 is the H.264
    // profile, flags, and level (respectively, one byte each). See
    // ITU-T H.264 specification for details.
    let mime = 'video/mp4; codecs="avc1.640028"';
    let sourceBuffer = mediaSource.addSourceBuffer(mime);
    let queue = [null];
    // received file or media segment
    ws.onmessage = function(event) {
      // console.log("message")
        if (sourceBuffer.updating) {
          queue[0]=event.data;
        } else {
          sourceBuffer.appendBuffer(event.data);
        }

        // if (sourceBuffer.updating || queue.length > 0) {
        //   queue=[event.data];
        // } else {
        //   sourceBuffer.appendBuffer(event.data);
        // }
      
    }
    sourceBuffer.addEventListener('updateend', function() {
      // console.log(queue.length)
      if (queue[0]!=null && !sourceBuffer.updating) {
        sourceBuffer.appendBuffer(queue[0]);
        queue[0]=null;
      }
      // if (queue.length > 0 && !sourceBuffer.updating) {
      //   sourceBuffer.appendBuffer(queue.shift());
      // }
    }
  , false);
    // remote closed websocket. end-of-stream.
    ws.onclose = function(event) {
      mediaSource.endOfStream();
    }
  }
};
