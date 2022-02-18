# BerryMSE Demo

This is a demo application that can build and run without many preparation. The release download contains this executable for the Raspberry Pi with armv7l architecture.

## Build

To fetch dependencies:

Linux:

    GOOS=linux go get -v ./...
    go install github.com/markbates/pkger/cmd/pkger


Windows:

    set GOOS=linux
    go get -v ./...
    go install github.com/markbates/pkger/cmd/pkger

To build:

    make

or:

    make armv6l
    make armv7l

A folder will be created that contains the executable

## Usage

To run, copy the appropriate `berrymse` executable to the Raspberry Pi and perform ``chmod +x berrymse``. Now you can run it with:

	./berrymse -l <raspberry pi ip address>:2020 -d /dev/video<X>

For example:

    ./berrymse -l 192.168.2.1:2020 -d /dev/video0

A more detailed usage instruction can be found in the [README_Executable](README_Executable.md).