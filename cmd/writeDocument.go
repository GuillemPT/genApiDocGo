package main

import (
	"fmt"
	"log"
	"os"
)

const BASE_DOCUMENT = "openapi: 3.0.3\n" +
	"info:\n" +
	"  title: Sample API\n" +
	"  description: API description.\n" +
	"  version: 1.0.0\n" +
	"paths:\n"

func writeDocument(structuredMethods []methodDoc, path string) {
	swagger, err := os.Create(path + "/swagger.yaml")
	if err != nil {
		log.Fatal(err)
	}

	defer swagger.Close()
	swagger.Write([]byte(BASE_DOCUMENT))

	for _, structuredMethod := range structuredMethods {
		structToString := fmt.Sprintf(
			" %s:\n  %s:\n    description: %s\n    responses:\n     "+
				"%d:\n      description: %s\n",
			structuredMethod.Route, structuredMethod.Type,
			structuredMethod.Description, structuredMethod.Response,
			structuredMethod.ResponseDescription)
		swagger.Write([]byte(structToString))
	}
}
