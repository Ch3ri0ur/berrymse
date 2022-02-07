window.onload = function () {
    let videoElement = document.querySelector("video");
    // select button element

    let playButton = document.getElementById("play-button");
    let autoSkipButton = document.getElementById("auto-skip-button");

    // select log button element
    let resetButton = document.getElementById("reset-button");
    let ws = null;
    let auto_skip = true;
    playButton.addEventListener("click", function () {
        if (videoElement.paused) {
            videoElement.play();
            playButton.textContent = "Pause";
        } else {
            videoElement.pause();
            playButton.textContent = "Play";
        }
    });

    autoSkipButton.addEventListener("click", function () {
        if (auto_skip) {
            auto_skip = false;
            autoSkipButton.textContent = "Auto Skip Disabled";
        } else {
            auto_skip = true;
            autoSkipButton.textContent = "Auto Skip Enabled";
        }
    });

    resetButton.addEventListener("click", function () {
        if (window.MediaSource) {
            if (ws !== null) {
                ws.close();
            }
            let mediaSource = new MediaSource();
            videoElement.loop = false;
            videoElement.src = URL.createObjectURL(mediaSource);
            mediaSource.addEventListener("sourceopen", sourceOpen);
        } else {
            console.log("Media Source Extensions API is NOT supported");
        }
    });

    // Check that browser supports Media Source Extensions API
    if (window.MediaSource) {
        let mediaSource = new MediaSource();
        videoElement.loop = false;
        videoElement.src = URL.createObjectURL(mediaSource);
        mediaSource.addEventListener("sourceopen", sourceOpen);
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

        // start websocket connection function
        function startConnection(params) {
            // remote pushes media segments via websocket
            ws = new WebSocket("ws://" + location.hostname + (location.port ? ":" + location.port : "") + "/websocket");
            // ws = new WebSocket("ws://" + "79.208.31.62" + "/websocket");

            ws.binaryType = "arraybuffer";
            let queue = [];

            // received file or media segment
            ws.onmessage = function (event) {
                // console.log("message")
                // check if frame is iframe or dframe
                let [is_mdat, is_iframe] = checkFrameType(event.data);
                if (!is_mdat) {
                    console.log("whoknowsframe");

                    sourceBuffer.appendBuffer(event.data);
                    console.log("not mdat");
                    return;
                }
                if (is_iframe === null) {
                    console.log("not iframe or dframe");

                    return;
                }
                if (is_iframe === true) {
                    // iframe
                    queue.length = 0;
                    // // sourceBuffer.remove(0, sourceBuffer.buffered.end(sourceBuffer.buffered.length));
                    if (sourceBuffer.updating) {
                        queue.push(event.data);
                    } else {
                        // console.log("iframe");
                        sourceBuffer.appendBuffer(event.data);
                    }

                    let buffered = videoElement.buffered;
                    if (buffered.length > 0) {
                        if (auto_skip && videoElement.currentTime < videoElement.buffered.end(videoElement.buffered.length - 1) - 1) {
                            videoElement.currentTime = videoElement.buffered.end(videoElement.buffered.length - 1) - 0.05;
                            console.log("document hidden:", document.hidden);
                            if (!document.hidden) {
                                videoElement.play();
                            }
                            console.log("triggered skip ahead");
                        }
                    }
                } else {
                    // dframe

                    if (sourceBuffer.updating || queue.length > 0) {
                        queue.push(event.data);
                    } else {
                        // console.log("dframe");
                        sourceBuffer.appendBuffer(event.data);
                    }
                }
            };
            // remote closed websocket. end-of-stream.
            ws.onclose = function (event) {
                console.log("remote closed websocket");
                // mediaSource.endOfStream();
                // sourceBuffer.removeEventListener("updateend", sourceBufferUpdateEnd);
                // mediaSource.removeSourceBuffer(sourceBuffer);
                // setTimeout(function() {
                //     startConnection();
                //   }, 5000);
            };
            ws.onerror = function (e) {
                console.log("Error: " + e.data);
            };
            onupdate = function () {
                if (queue.length > 0 && !sourceBuffer.updating) {
                    // console.log("updateend");
                    sourceBuffer.appendBuffer(queue.shift());
                }
            };
            sourceBuffer.addEventListener("updateend", onupdate, false);
        }

        // start websocket connection
        startConnection();
    }
};

function checkFrameType(data) {
    let mdat = new Uint8Array([116, 114, 117, 110]);
    let is_iframe = null;
    let is_mdat = false;
    if (arraybufferEqual(data.slice(76, 80), mdat.buffer)) {
        // console.log("mdat");
        is_mdat = true;

        let view = new Uint8Array(data.slice(92, 93));
        // check if iframe or dframe
        // console.log(view);
        if (view[0] === 1) {
            // console.log("dframe");
            is_iframe = false;
        } else if (view[0] === 2) {
            // console.log("iframe");
            is_iframe = true;
        }
    }
    return [is_mdat, is_iframe];
}

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
}
