package fileslogic

import (
	"encoding/json"
	"genApiDocGo/src/internal"
	"io"
	"log"
	"os"
)

// Take the formatted methods and write the swagger.json file.
func WriteDocument(structuredMethods map[string]internal.PathDocument,
	path string) {
	baseDocument := internal.GetBaseDocumentConfig()
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
