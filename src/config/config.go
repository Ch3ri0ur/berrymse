package config

import (

	//CMD Prints
	"fmt"

	//OS Function to access Path and files https://pkg.go.dev/os
	"os"

	//Error handling https://pkg.go.dev/errors
	"errors"

	//Easy to use Path and File utilities https://pkg.go.dev/path/filepath
	"path/filepath"

	// FLAG & Configuration Manager	################################################################

	//Improved Flag Handle compatible with Viper Configuration Manger https://github.com/spf13/pflag
	//Old CMD Flag Handle
	//"flag"
	//Default Values set for the Flags are used to set the Configuration
	//Imported flag to replace the old "flag" handle
	flag "github.com/spf13/pflag"

	//Configuration Manger https://github.com/spf13/viper
	//This Configuration Manager allows the use of a Configfile and Flags to set the Configuration
	//Priority is UsedFlag>Configfile>DefaultFlag!
	"github.com/spf13/viper"
)

// ConfigurationsStruct ##########################################################################################

// Configurations exported
type Configurations struct {
	Camera CameraConfigurations
	Server ServerConfigurations
}

// CameraConfigurations Struct exported
type CameraConfigurations struct {
	SourceFD string
	Width    int
	Height   int
	Bitrate  int
	Rotation int
}

// ServerConfigurations Struct exported
type ServerConfigurations struct {
	URL       string
	WebSocket string
}

// Flag & Configuration Loading ##################################################################################

//Command line flag parameters
//Tmp Flag Save Location for Strings DO NOT USE IN CODE!!
//USE configuration (Configurations) to access the config values
var flagConfig string

//Init methode
//Default Defining Flags and Default values
func DefaultFlagInit() {

	//Basic Type Flags
	//PFlagtype(ConfigID,
	//	FlagID,
	//	Defaultvalue,
	//	Info Text,
	//)
	//Stored in a extra Variable (Strings)
	//PFlagVartype(&Variable, ConfigID,
	//	FlagID,
	//	Defaultvalue,
	//	Info Text,
	//)

	//Config Path Flag
	//No config Name, only needed to load selected config file
	flag.StringVarP(&flagConfig, "",
		"c",
		"config.yml",
		"Use config Path/Name.yml"+
			"\nDefault Path is current directory!",
	)

	//Flag to selected an URL
	flag.StringP("Server.URL",
		"l",
		"localhost:2020",
		"listen on host:port",
	)

	//Flag to selected an WebsocketName
	flag.StringP("Server.WebSocket",
		"s",
		"video_websocket",
		"Name of Websocket for Video Stream",
	)

	//Flag to change the Device input file / device nodes
	flag.StringP("Camera.SourceFD",
		"d",
		"/dev/video0",
		"Use camera /dev/videoX",
	)

	//Flag to change the width of the camera and video
	flag.IntP("Camera.Width",
		"w",
		1280,
		"Width Resolution",
	)

	//Flag to change the height of the camera and video
	flag.IntP("Camera.Height",
		"h",
		720,
		"Height Resolution",
	)

	//Flag to change the bitrate video
	flag.IntP("Camera.Bitrate",
		"b",
		1500000,
		"Bitrate in bit/s!\nOnly supported for RPI Camera\nOther Cameras need to use -1",
	)

	//Flag to change the rotation video
	flag.IntP("Camera.Rotation",
		"r",
		0,
		"Rotation in 90degree Step\nOnly supported for RPI Camera\nOther Cameras need to use -1",
	)
}

//All Configurations Stored in this. Look config.go for structure
var configuration Configurations

//Reads Flags and Configfile to set and overwrite the Config
func SetupConfigFlags() Configurations {

	//Get all Flags and Parse them in Variables
	flag.Parse()
	//Bind Flags to Config
	viper.BindPFlags(flag.CommandLine)
	//Not bound variables
	//viper.SetDefault("Camera.FD", "/dev/video1")

	//Checks for and Loads Configfile
	loadConfigs()

	//Loads the Config into the Struct for easier usage
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v\n", err)
	}

	//Showcase of usage of Config and some test code
	//fmt.Printf("Camera FD Configuration: %s \n", configuration.Camera.SourceFD)
	//fmt.Printf("Server URL Configuration: %s \n", configuration.Server.URL)
	return configuration
}

//Checks if Configfile exists and read it
//When flagConfig only contains a Path it will use the default config name "config.yml"
func loadConfigs() {
	//Checks if the Configfile exists. uses the flag Value that contains a given or the default ConfigfilePath. Skipping the Loading of Configfile if not exists
	if _, err := os.Stat(flagConfig); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Error: Configfile not found: %s \n", flagConfig)
		fmt.Printf("Skip loading Configfile! Using default Settings!\n")
		return
	}

	//Path and Filename with extention gets split up
	dir, file := filepath.Split(flagConfig)

	//If the Path is empty (current directory) a "." needs to be used
	if dir == "" {
		dir = "."
	}

	// If only a Path was given use the default Configfilename "config.yml"
	if file == "" {
		file = "config.yml"
		fmt.Printf("Warning. Flag -c contained a Filepath without filename, %s \n"+
			"Using default \"config.yml\" as config file.\n", flagConfig)
	}
	//Extract the Filename by removing the Extention
	fileName := file[:len(file)-len(filepath.Ext(file))]

	//Extract the Extention
	fileExtention := filepath.Ext(file)
	//The "." of the extention needs to be removed
	if fileExtention != "" {
		fileExtention = fileExtention[1:]
	}

	//Check if Extention of the is Supported. Skip Loading the Configfile, when not supported!
	if fileExtention != "yml" {
		fmt.Printf("Error. Flag -c contained an File with wrong file extention, %s \n", flagConfig)
		fmt.Printf("Skip loading Configfile! Using default Settings!\n")
		return
	}

	// Set the path, name and extention to look for the configurations file
	viper.AddConfigPath(dir)
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileExtention)

	//Read the Configuration File
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}
}
