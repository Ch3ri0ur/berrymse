package main

import (
	//CMD Prints
	"fmt"
	config "github.com/ch3ri0ur/berrymse/src/config"
)

func init() {
	config.FlagInit()
}

func main() {
	fmt.Println("Load Config")
	var configuration = config.SetupConfigFlags()
	fmt.Println(configuration)
}
