package fileslogic

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type swaggerDocument struct {
	Openapi string                  `json:"openapi"`
	Info    infoDocument            `json:"info"`
	Paths   map[string]pathDocument `json:"paths"` // key is the route (/home)
}

type infoDocument struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// key is the type of operation (get)
type pathDocument map[string]operationDocument

type operationDocument struct {
	Description string `json:"description"`
	// key is the code id of response (200)
	Responses map[string]responsesDocument `json:"responses"`
}

type responsesDocument struct {
	Description string `json:"description"`
}

// TODO: change a json, i think is easier to manage
var baseDocument = swaggerDocument{Openapi: "3.0.3",
	Info: infoDocument{Title: "API name", Description: "API description",
		Version: "1.0.0"}}

func WriteDocument(structuredMethods map[string]pathDocument, path string) {
	swagger, err := os.Create(path + "/swagger.json")
	if err != nil {
		log.Fatal(err)
	}

	defer swagger.Close()
	baseDocument.Paths = structuredMethods
	jsonData, err := json.MarshalIndent(baseDocument, "", " ")

	if err != nil {
		log.Panic("Error marshalling to JSON:", err)
		return
	}

	_, err = io.Writer.Write(swagger, jsonData)

	if err != nil {
		log.Panic("Error when write swagger: ", err)
	}
}
