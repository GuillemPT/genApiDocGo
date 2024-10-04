package main

import (
	"log"
	"os"

	"genApiDocGo/src/fileslogic"
	"genApiDocGo/src/internal"

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

	fileType, _ := RunSelector(internal.GetFileTypeOptions(), SingleSelection,
		"Enter type of files to generate the documentation:", internal.ViewHeigh)

	_, _ = RunSelector(internal.GetFrameworks(fileType), SingleSelection,
		"Enter the framework used:", internal.ViewHeigh)

	files, directories, err := fileslogic.GetFiles(targetDirectoryPath,
		fileType)

	if err != nil {
		log.Fatal("Error when extracting files", err)
	}

	_, excludeDirectories := RunSelector(directories, MultiSelection,
		"Select the directories that want exclude:", internal.ViewHeigh)

	if err != nil {
		log.Fatal("Error obtaining exclude directories", err)
	}

	files = fileslogic.ExcludeFilesInBanDirectories(excludeDirectories, files)

	methodsToDoc := fileslogic.GetContent(files)
	structMethods := fileslogic.FormatMethods(methodsToDoc)
	fileslogic.WriteDocument(structMethods, targetDirectoryPath)
}
