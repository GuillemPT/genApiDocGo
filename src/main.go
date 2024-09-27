package main

import (
	"fmt"
	"os"
)

func main() {
	/*
	 * Get arguments by console, there are four optional parameters
	 * 1 - type of files to generate the documentation
	 * (by default js)
	 * 2 - directory to browse the files to generate the documentation
	 * (by default current)
	 * 3 - destination to write the generated documentation
	 * (by default current directory with the name that doc.txt)
	 * 4 - exclude directories (all arguments based on it are considered to be
	 * part of the same)
	 */

    var targetDirectoryPath string
    var writeDocumentationDirectoryPath string
    var filesType string
	var excludeDirectories []string
	numberOfArgs := len(os.Args)

	switch numberOfArgs {
	case 1:
		filesType = "js"
		targetDirectoryPath, _ = os.Getwd()
		writeDocumentationDirectoryPath, _ = os.Getwd()
	case 2:
		filesType = os.Args[1]
		targetDirectoryPath, _ = os.Getwd()
		writeDocumentationDirectoryPath, _ = os.Getwd()
	case 3:
		filesType = os.Args[1]
		targetDirectoryPath = os.Args[2]
		writeDocumentationDirectoryPath, _ = os.Getwd()
	case 4:
		filesType = os.Args[1]
		targetDirectoryPath = os.Args[2]
		writeDocumentationDirectoryPath = os.Args[3]
	default:
		filesType = os.Args[1]
		targetDirectoryPath = os.Args[2]
		writeDocumentationDirectoryPath = os.Args[3]
		excludeDirectories = os.Args[4:]
	}

	files, err := getFiles(targetDirectoryPath, filesType, excludeDirectories)
	if err == nil {
		getContent(files)
	}

	//! Delete the next line, now is here to not throw errors
	fmt.Println(writeDocumentationDirectoryPath)
}