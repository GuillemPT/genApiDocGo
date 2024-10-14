package main

import (
	"log"
	"os"

	"genApiDocGo/src/fileslogic"
	"genApiDocGo/src/internal"

	"github.com/common-nighthawk/go-figure"
)

const (
	OneArgument int = iota + 1
	TwoArguments
)

func main() {
	// Draw the name and log the version
	figure.NewFigure(internal.AppName, internal.LettersType, true).Print()
	log.Print("Start "+internal.AppName+" with version: ", internal.Version)

	var targetDirectoryPath string
	var configPath string

	switch len(os.Args) {
	case OneArgument:
		// Use this path as the default.
		targetDirectoryPath, _ = os.Getwd()
	case TwoArguments:
		targetDirectoryPath = os.Args[1]
	default:
		targetDirectoryPath = os.Args[1]
		configPath = os.Args[2]
	}

	internal.SetConfiguration(configPath)

	// Selector to chose the file type to search.
	fileType, _ := RunSelector(internal.GetFileTypeOptions(), SingleSelection,
		"Enter type of files to generate the documentation "+
			"(by default .js):", internal.ViewHeigh)

	if fileType == "" {
		fileType = internal.DefaultFileType
	}

	log.Println("File type set: ", fileType)

	// Selector to chose the framework (currently only active express).
	frameworks := internal.GetFrameworks(fileType)
	framework, _ := RunSelector(frameworks, SingleSelection,
		"Enter the framework used (by default the first element of the list):",
		internal.ViewHeigh)

	if framework == "" {
		framework = frameworks[internal.DefaultFramework]
	}

	log.Println("Framework set: ", framework)
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
