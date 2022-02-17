package config

import (
	//Testing
	"testing"

	//CMD Prints
	"fmt"
)

func init() {
	FlagInit()
}

func Test_setupConfigFlags(t *testing.T) {
	fmt.Println("Load Config")
	var configuration = SetupConfigFlags()
	fmt.Println(configuration)
}
