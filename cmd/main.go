package main

import (
	"fmt"
	"genApiDocGo/internal"
	"log"
	"os"

	"github.com/common-nighthawk/go-figure"
)

type item struct {
	ID         string
	IsSelected bool
}

func main() {
	/*
	 * Get arguments by console, there are four optional parameters
	 * 1 - type of files to generate the documentation
	 * 2 - exclude directories
	 */

	figure.NewFigure("GenApiDocGo", "", true).Print()
	fmt.Println("Version: ", internal.Version)

	var targetDirectoryPath string

	if len(os.Args) > 1 {
		targetDirectoryPath = os.Args[1]
	} else {

		targetDirectoryPath, _ = os.Getwd()
	}

	files, directories, err := getFiles(targetDirectoryPath,
		getSelectLanguage())

	if err != nil {
		log.Fatal("Error when extracting files", err)
	}

	excludeDirectories, err := getExcludedDirectories(0, directories)

	if err != nil {
		log.Fatal("Error obtaining exclude directories", err)
	}

	files = excludeFilesInBanDirectories(excludeDirectories, files)
	fmt.Println(files)
	methodsToDoc := getContent(files)
	structMethods := formatMethods(methodsToDoc)
	writeDocument(structMethods, targetDirectoryPath)
}
