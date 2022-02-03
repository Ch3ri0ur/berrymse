window.onload = function () {

  let videoElement = document.querySelector('video');
  // select button element

  let playButton = document.getElementById("play-button");
  // select log button element
  let logButton = document.getElementById("log-button");


  playButton.addEventListener('click', function () {
    if (videoElement.paused) {
      videoElement.play();
      playButton.textContent = "Pause";
    } else {
      videoElement.pause();
      playButton.textContent = "Play";
    }
  });

  // Check that browser supports Media Source Extensions API
  if (window.MediaSource) {
    let mediaSource = new MediaSource();
    videoElement.loop = false;
    videoElement.src = URL.createObjectURL(mediaSource);
    mediaSource.addEventListener('sourceopen', sourceOpen);
    // videoElement.onpause = function () {
    //   console.log("buffered:", videoElement.buffered);
    //   let buffered = videoElement.buffered;
    //   videoElement.currentTime = buffered.end(buffered.length - 1) - 0.2;
    //   videoElement.play();
    // }
  } else {
    console.log("Media Source Extensions API is NOT supported");
  }

  function sourceOpen(e) {
    URL.revokeObjectURL(videoElement.src);

    // The six hexadecimal digit suffix after avc1 is the H.264
    // profile, flags, and level (respectively, one byte each). See
    // ITU-T H.264 specification for details.
    let mediaSource = e.target;
    let mime = 'video/mp4; codecs="avc1.640028"';
    let sourceBuffer = mediaSource.addSourceBuffer(mime);


    // remote pushes media segments via websocket
    // ws = new WebSocket("ws://" + location.hostname + (location.port ? ":"+location.port : "" ) + "/websocket");
    ws = new WebSocket("ws://" + "wpplr.cc" + "/websocket");

    ws.binaryType = "arraybuffer";

    logButton.addEventListener('click', function () {
      // console.log(videoElement.)
      console.log(videoElement);
      console.log(mediaSource);
      console.log(sourceBuffer);
    });
    let queue = [];

    function arraybufferEqual(buf1, buf2) {
      if (buf1 === buf2) {
        return true;
      }

      if (buf1.byteLength !== buf2.byteLength) {
        return false;
      }

      var view1 = new DataView(buf1);
      var view2 = new DataView(buf2);

      var i = buf1.byteLength;
      while (i--) {
        if (view1.getUint8(i) !== view2.getUint8(i)) {
          return false;
        }
      }
      return true;
    };

    // received file or media segment
    ws.onmessage = function (event) {
      // console.log("message")
      // if (sourceBuffer.updating) {
      //   queue[0]=event.data;
      // } else {
      //   sourceBuffer.appendBuffer(event.data);
      
      // }
      let mdat = new Uint8Array([116, 114, 117, 110])
      let is_iframe = null;
      if (arraybufferEqual(event.data.slice(76, 80), mdat.buffer)) {
        // console.log("mdat");
        let view = new Uint8Array(event.data.slice(92, 93));
        // console.log(view);
        if (view[0] === 1) {

          // console.log("dframe");
          is_iframe = false;
        }

        else if (view[0] === 2) {
          // console.log("iframe");
          is_iframe = true;
        }
        if (is_iframe === null) {
          console.log("not iframe or dframe");
          return;
        }

        if (is_iframe === true) {
          queue.length = 0;
          // // sourceBuffer.remove(0, sourceBuffer.buffered.end(sourceBuffer.buffered.length));
          if (sourceBuffer.updating) {
            queue.push(event.data);
          } else {
            console.log("iframe");
            sourceBuffer.appendBuffer(event.data);
          }
          //   if (sourceBuffer.updating) {
          //     sourceBuffer.abort();
          //   } 
          //   if (sourceBuffer.buffered.length > 0) {
          //   sourceBuffer.remove( sourceBuffer.buffered.start(0), sourceBuffer.buffered.end(0));
          // }
          // sourceBuffer.appendBuffer(event.data);


          let buffered = videoElement.buffered;
          if (buffered.length > 0) {
            if (videoElement.currentTime < videoElement.buffered.end(videoElement.buffered.length - 1) - 1) {
              videoElement.currentTime = videoElement.buffered.end(videoElement.buffered.length - 1) - 0.05;
              console.log("document hidden:", document.hidden);
              if (!document.hidden) {
                videoElement.play();

              }
              console.log("triggered skip ahead")
            }
          }
        } else {


          if (sourceBuffer.updating || queue.length > 0) {
            queue.push(event.data);
          } else {
            console.log("dframe");

            sourceBuffer.appendBuffer(event.data);
          }
        }

      } else {
        console.log("whoknowsframe");

        sourceBuffer.appendBuffer(event.data);
        console.log("not mdat")
      }
    }
    sourceBuffer.addEventListener('updateend', function () {
      // console.log(queue.length)
      // if (queue[0]!=null && !sourceBuffer.updating) {
      //   sourceBuffer.appendBuffer(queue[0]);
      //   queue[0]=null;
      // }
      if (queue.length > 0 && !sourceBuffer.updating) {
        console.log("updateend");

        sourceBuffer.appendBuffer(queue.shift());
      }
    }
      , false);
    // remote closed websocket. end-of-stream.
    ws.onclose = function (event) {
      mediaSource.endOfStream();
    }
  }
};
