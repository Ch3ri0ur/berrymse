package main

// Configurations exported
type Configurations struct {
	Camera CameraConfigurations
	Server ServerConfigurations
}

// CameraConfigurations Struct exported
type CameraConfigurations struct {
	SourceFD string
}

// ServerConfigurations Struct exported
type ServerConfigurations struct {
	URL string
}
