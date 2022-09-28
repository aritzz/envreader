package main

import (
	"fmt"
	"os"

	"github.com/aritzz/envreader"
)

// Define some data structure with environment values
// and defaults
type ConfigCustom struct {
	ListenHost string `envvalue:"LISTEN_HOST" default:"127.0.0.1"`
	ListenPort string `envvalue:"LISTEN_PORT" default:"5000"`
	Debug      bool   `envvalue:"ENABLE_DEBUG"`
}

func main() {
	var config ConfigCustom

	// Initialize module
	reader := envreader.EnvReader{}
	reader.Init()
	reader.SetTagName("envvalue")

	// Read environment variables to structure
	if err := reader.Read(&config); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Print values
	fmt.Printf("%+v\n", config)
}
