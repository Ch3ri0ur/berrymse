package main

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
}

// ServerConfigurations Struct exported
type ServerConfigurations struct {
	URL string
}
