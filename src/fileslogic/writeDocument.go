package fileslogic

import (
	"fmt"
	"io"
	"log"
	"os"
)

// TODO: change a json, i think is easier to manage
const baseDocument = "openapi: 3.0.3\n" +
	"info:\n" +
	"  title: Sample API\n" +
	"  description: API description.\n" +
	"  version: 1.0.0\n" +
	"paths:\n"

func WriteDocument(structuredMethods []MethodDoc, path string) {
	swagger, err := os.Create(path + "/swagger.yaml")
	if err != nil {
		log.Fatal(err)
	}

	defer swagger.Close()
	_, err = io.WriteString(swagger, baseDocument)

	if err != nil {
		log.Panic("Error when write swagger: ", err)
	}

	for _, structuredMethod := range structuredMethods {
		structToString := fmt.Sprintf(
			" %s:\n  %s:\n    description: %s\n    responses:\n     "+
				"%d:\n      description: %s\n",
			structuredMethod.Route, structuredMethod.Type,
			structuredMethod.Description, structuredMethod.Response,
			structuredMethod.ResponseDescription)
		_, err = io.WriteString(swagger, structToString)

		if err != nil {
			log.Panic("Error when write swagger: ", err)
		}
	}
}
