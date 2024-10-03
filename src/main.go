package main

import (
	"genApiDocGo/src/fileslogic"
	"genApiDocGo/src/internal"
	"log"
	"os"

	"github.com/common-nighthawk/go-figure"
)

func main() {
	figure.NewFigure("GenApiDocGo", "slant", true).Print()
	log.Print("Start GenApiDocGo with version: ", internal.Version)

	var targetDirectoryPath string

	if len(os.Args) > 1 {
		targetDirectoryPath = os.Args[1]
	} else {
		targetDirectoryPath, _ = os.Getwd()
	}
	fileType := getUniqueSelect("Enter type of files to generate "+
		"the documentation", internal.GetFileTypeOptions())

	_ = getUniqueSelect("Enter the framework used",
		internal.GetFrameworks(fileType))

	files, directories, err := fileslogic.GetFiles(targetDirectoryPath,
		fileType)

	if err != nil {
		log.Fatal("Error when extracting files", err)
	}

	excludeDirectories, err := getExcludedDirectories(0, directories)

	if err != nil {
		log.Fatal("Error obtaining exclude directories", err)
	}

	files = fileslogic.ExcludeFilesInBanDirectories(excludeDirectories, files)

	methodsToDoc := fileslogic.GetContent(files)
	structMethods := fileslogic.FormatMethods(methodsToDoc)
	fileslogic.WriteDocument(structMethods, targetDirectoryPath)
}
