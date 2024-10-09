package main

import (
	"log"
	"os"

	"genApiDocGo/src/fileslogic"
	"genApiDocGo/src/internal"

	"github.com/common-nighthawk/go-figure"
)

func main() {
	// Draw the name and log the version
	figure.NewFigure(internal.AppName, internal.LettersType, true).Print()
	log.Print("Start "+internal.AppName+" with version: ", internal.Version)

	var targetDirectoryPath string

	if len(os.Args) > 1 {
		// Set the path target.
		targetDirectoryPath = os.Args[1]
	} else {
		// Use this path as the default.
		targetDirectoryPath, _ = os.Getwd()
	}

	// Selector to chose the file type to search.
	fileType, _ := RunSelector(internal.GetFileTypeOptions(), SingleSelection,
		"Enter type of files to generate the documentation:", internal.ViewHeigh)

	// Selector to chose the framework (currently only active express).
	_, _ = RunSelector(internal.GetFrameworks(fileType), SingleSelection,
		"Enter the framework used:", internal.ViewHeigh)

	files, directories, err := fileslogic.GetFiles(targetDirectoryPath,
		fileType)

	if err != nil {
		log.Fatal("Error when extracting files", err)
	}

	// Selector to choose which directories you don't want to look at, this has
	// two functionalities:
	// 1- read less files and make the process faster.
	// 2- in languages like TS the output address may cause problems.
	_, excludeDirectories := RunSelector(directories, MultiSelection,
		"Select the directories that want exclude:", internal.ViewHeigh)

	if err != nil {
		log.Fatal("Error obtaining exclude directories", err)
	}

	// Exclude files that are inside the ""banned" directories.
	files = fileslogic.ExcludeFilesInBanDirectories(excludeDirectories, files)

	// Extract the methods to documentation.
	methodsToDoc := fileslogic.GetContent(files)
	// Format the methods to make them ready to be written.
	structMethods := fileslogic.FormatMethods(methodsToDoc)
	// Write the JSON documentation
	fileslogic.WriteDocument(structMethods, targetDirectoryPath)
}
