(function (window, undefined) {
    function initialize(videoElementObj, wsTargetURL, autoSkipButtonObj, restButtonObj) {
        // Check that browser supports Media Source Extensions API
        if (!window.MediaSource) {
            console.log("Media Source Extensions API is NOT supported");
            return;
        }
        if (wsTargetURL == null) {
            wsTargetURL = "ws://" + location.hostname + (location.port ? ":" + location.port : "") + "/video_websocket";
        }
        let targetURL = wsTargetURL;

        let ws = null;
        let videoElement = videoElementObj;
        let mediaSource = new MediaSource();
        videoElement.loop = false;
        videoElement.src = URL.createObjectURL(mediaSource);
        mediaSource.addEventListener("sourceopen", sourceOpen);

        //Auto Skip Button
        let auto_skip = true;
        let autoSkipButton = autoSkipButtonObj;
        if (autoSkipButton) {
            autoSkipButton.textContent = "Disable Auto Skip";
            autoSkipButton.addEventListener("click", function () {
                if (auto_skip) {
                    auto_skip = false;
                    autoSkipButton.textContent = "Enable Auto Skip";
                } else {
                    auto_skip = true;
                    autoSkipButton.textContent = "Disable Auto Skip";
                }
            });
        }

        //Reset Button
        let resetButton = restButtonObj;
        if (resetButton) {
            resetButton.textContent = "Reset";
            resetButton.addEventListener("click", function () {
                if (ws !== null) {
                    ws.close();
                }
                let mediaSource = new MediaSource();
                videoElement.loop = false;
                videoElement.src = URL.createObjectURL(mediaSource);
                mediaSource.addEventListener("sourceopen", sourceOpen);
            });
        }

        function sourceOpen(e) {
            URL.revokeObjectURL(videoElement.src);

            // The six hexadecimal digit suffix after avc1 is the H.264
            // profile, flags, and level (respectively, one byte each). See
            // ITU-T H.264 specification for details.
            let mediaSource = e.target;
            let sourceBuffer = null;

            // start websocket connection function
            function startConnection() {
                // remote pushes media segments via websocket
                ws = new WebSocket(wsTargetURL);

                ws.binaryType = "arraybuffer";
                // queue for saving frames to be played
                let queue = [];
                let is_first = true;

                // received media segment
                ws.onmessage = function (event) {
                    // check data is a frame (has mdat) and if it is an inter(i)-frame
                    let [is_mdat, is_iframe] = checkFrameType(event.data);
                    // if it is not a frame it is probably a setup packet and still goes into the Buffer
                    if (!is_mdat) {
                        // Get the first websocket message ( it should have ftyp moov with information about the stream)
                        if (is_first) {
                            is_first = false;
                            // dynamically set the codec string (only works for h.264 in mp4 with avcC box)
                            let codec_string = findCodecString(event.data); // profile, profile compatibility and level as hex string
                            console.log("dynamically set codec string: " + codec_string);
                            let mime = 'video/mp4; codecs="avc1.' + codec_string + '"';
                            // added sourcebuffer here, because information about the codec is necessary
                            sourceBuffer = mediaSource.addSourceBuffer(mime);
                            // add update function for sourcebuffer
                            onupdate = function () {
                                if (queue.length > 0 && !sourceBuffer.updating) {
                                    sourceBuffer.appendBuffer(queue.shift());
                                }
                            };
                            sourceBuffer.addEventListener("updateend", onupdate, false);
                            sourceBuffer.appendBuffer(event.data);
                            return;
                        }
                        // this should not be reached^
                        console.log("Got another ftyp and moov packet ");
                    }
                    // data is a frame but frame type does not match -> skip
                    if (is_iframe === null) {
                        console.log("not iframe or pframe");
                        return;
                    }
                    // data is an iframe
                    if (is_iframe === true) {
                        // reset queue
                        queue.length = 0;
                        // if the sourceBuffer is still updating use queue otherwise append buffer
                        if (sourceBuffer.updating) {
                            queue.push(event.data);
                        } else {
                            sourceBuffer.appendBuffer(event.data);
                        }

                        // if auto skip is enabled and to much is buffered skip forward, except if the video is hidden (tab is in background)
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
                        // data is not an iframe
                        // if the sourceBuffer is still updating or there are frames in the queue, use the queue. Otherwise append to buffer.
                        if (sourceBuffer.updating || queue.length > 0) {
                            queue.push(event.data);
                        } else {
                            sourceBuffer.appendBuffer(event.data);
                        }
                    }
                };

                // remote closed websocket. end-of-stream.
                // Looked into automatically restarting the stream, but that creates a recursive problem right now.
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

                //update function for sourceBuffer, it checks if there are any frames in the queue and if so appends them
            }

            // start websocket connection
            startConnection();
        }
    };

    function findCodecString(data) {
        // check in trun box if data is an iframe or pframe
        let data_array = new Uint8Array(data);

        // check function to search for index of avcC box
        let avcC = new Uint8Array([97, 118, 99, 67]); // spells avcC
        function isavcC(element, index, array) {
            // if the first element matches, check the rest else return false
            if (element == 97) {
                let tmp = array.slice(index, index + 4);
                if (arraybufferEqual(tmp.buffer, avcC.buffer)) {
                    return true;
                }
            }
            return false;
        }

        index = data_array.findIndex(isavcC);
        if (index === -1) {
            throw Error("avcC not found");
        }
        // get condec data relative to found index
        let codec_data_buffer = data.slice(index + 4 + 1, index + 8);
        let bufferToHex = (buffer) => {
            return [...new Uint8Array(buffer)].map((b) => b.toString(16).padStart(2, "0")).join("");
        };
        return bufferToHex(codec_data_buffer);
    }

    function checkFrameType(data) {
        // check in trun box if data is an iframe or pframe
        let mdat = new Uint8Array([116, 114, 117, 110]); // spells trun
        let is_iframe = null;
        let is_mdat = false;
        if (arraybufferEqual(data.slice(76, 80), mdat.buffer)) {
            // check if mdat is in correct position if yes -> it is a frame (has mdat)
            is_mdat = true;

            let view = new Uint8Array(data.slice(92, 93));
            // if it is a frame we can check if it is an iframe
            if (view[0] === 1) {
                // console.log("not iframe");
                is_iframe = false;
            } else if (view[0] === 2) {
                // console.log("iframe");
                is_iframe = true;
            }
        }
        return [is_mdat, is_iframe];
    }

    // helper function to check if two arraybuffers are equal
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

    window.BerryMSE = {
        initialize: initialize
    };
})(window)