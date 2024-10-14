package internal

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func SetConfiguration(path string) {
	if path == "" {
		log.Println("No configuration path has been passed, use default " +
			"configuration.")
		config.ResponsesMap = getDefaultResponses()
		config.BaseDocument = getDefaultBaseDocument()
		config.BaseDocument.Paths = make(map[string]PathDocument)
		return
	}

	file, err := os.Open(path)

	if err != nil {
		log.Panic("Error opening the configuration file: ", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Panic("Error reading the configuration file: ", err)
	}

	if err = json.Unmarshal(data, &config); err != nil {
		log.Panic("Error when unmarshal the configuration file: ", err)
	}
	config.BaseDocument.Paths = make(map[string]PathDocument)
}
