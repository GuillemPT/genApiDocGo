package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

/*
 * return array with each method to doc
 */
func getContent(files []string) []string {
	var extractedFunctions []string

	for _, file := range files {
		fileData, err := os.Open(file)

		if err != nil {
			log.Fatal(err)
		}

		defer fileData.Close()
		scanner := bufio.NewScanner(fileData)
		var methodsToDoc strings.Builder
		inMethod := false
		braceCounter := -1

		for scanner.Scan() {
			line := scanner.Text()

			if inMethod {
				methodsToDoc.WriteString(line + "\n")
				if strings.Contains(line, "{") {
					if braceCounter == -1 {
						braceCounter = 0
					}
					braceCounter++
				}

				if strings.Contains(line, "}") {
					braceCounter--
				}

				if braceCounter == 0 {
					inMethod = false
					braceCounter = -1
					extractedFunctions = append(extractedFunctions,
						methodsToDoc.String())
					methodsToDoc.Reset()

				}
			}

			if strings.Contains(line, "@api_generate_doc") {
				inMethod = true
			}
		}

	}

	return extractedFunctions
}
