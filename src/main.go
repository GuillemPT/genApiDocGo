package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	/*
	 * Get arguments by console, there are four optional parameters
	 * 1 - type of files to generate the documentation
	 * (by default js)
	 * 2 - directory to browse the files to generate the documentation
	 * (by default current)
	 * 3 - exclude directories (all arguments based on it are considered to be
	 * part of the same)
	 */
	fmt.Println("Start GenApiDocGo version: ", Version)

	var filesType string
	var targetDirectoryPath string
	var excludeDirectories string

	defaultFilesType := "js"
	defaultTargetDirectoryPath, _ := os.Getwd()
	defaultExcludeDirectories := ""
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter type of files to generate the documentation " +
		"(default js): ")
	filesType, err := reader.ReadString('\n')
	filesType = strings.TrimSpace(filesType)

	if err != nil || filesType == "" {
		filesType = defaultFilesType
	}

	fmt.Println("Enter directory to browse the files to generate the " +
		"documentation (default current): ")
	targetDirectoryPath, err = reader.ReadString('\n')
	targetDirectoryPath = strings.TrimSpace(targetDirectoryPath)

	if err != nil || targetDirectoryPath == "" {
		targetDirectoryPath = defaultTargetDirectoryPath
	}

	fmt.Println("Enter exclude directories: ")
	excludeDirectories, err = reader.ReadString('\n')
	excludeDirectories = strings.TrimSpace(excludeDirectories)

	if err != nil {
		excludeDirectories = defaultExcludeDirectories
	}

	excludeDirectoriesArr := strings.Split(excludeDirectories, " ")
	files, err := getFiles(
		targetDirectoryPath, filesType, excludeDirectoriesArr)

	if err != nil {
		log.Fatal(err)
	}

	methodsToDoc := getContent(files)
	structedMethods := formatMethods(methodsToDoc)
	fmt.Println(structedMethods)
}
