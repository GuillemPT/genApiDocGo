package fileslogic

import (
	"bufio"
	"genApiDocGo/src/internal"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetFiles returns files ([]string) that its extension match with fileType
// parameter, directories ([]string) and error.
func GetFiles(targetDirectoryPath string,
	filesType string) ([]string, []string, error) {
	var files []string
	var directories []string

	err := filepath.WalkDir(targetDirectoryPath,
		func(path string, file fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if file.IsDir() {
				if string(file.Name()[0]) == "." {
					return filepath.SkipDir
				}
				directories = append(directories, file.Name())
			} else if filepath.Ext(file.Name()) == filesType {
				files = append(files, path)
			}

			return nil
		})

	if err != nil {
		return nil, nil, err
	}

	return files, directories, nil
}

// Compare the current files and the excluded directories and if a file have not
// a ban directory in the path, tmpFiles variable. This variable will contain
// the files to extract the content.
func ExcludeFilesInBanDirectories(directories []string,
	files []string) []string {
	var tmpFiles []string

	for _, file := range files {
		inBanDirectory := false
		for _, directory := range directories {
			if strings.Contains(file, directory) {
				inBanDirectory = true
			}
		}
		if !inBanDirectory {
			tmpFiles = append(tmpFiles, file)
		}
	}

	return tmpFiles
}

// Return array with each method to create documentation.
func GetContent(files []string) []string {
	var extractedFunctions []string

	for _, file := range files {
		fileData, err := os.Open(file)

		if err != nil {
			log.Fatal(err)
		}

		defer fileData.Close()
		scanner := bufio.NewScanner(fileData)
		processFile(scanner, &extractedFunctions)
	}
	return extractedFunctions
}

// get a file and extract all methods mark with the tag.
func processFile(scanner *bufio.Scanner, extractedFunctions *[]string) {
	var inMethod bool
	var braceCounter int
	var methodsToDoc strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if inMethod {
			inMethod, braceCounter = processMethod(line, braceCounter,
				&methodsToDoc, extractedFunctions)
		} else if strings.Contains(line, internal.Tag) {
			inMethod = true
			braceCounter = 0
		}
	}
}

// return the bool value, true if have to process the next line or false if not
// and the current counter of brace.
func processMethod(line string, braceCounter int, methodsToDoc *strings.Builder,
	extractedFunctions *[]string) (bool, int) {
	methodsToDoc.WriteString(line + "\n")

	openBraces := strings.Count(line, "{")
	closesBraces := strings.Count(line, "}")
	braceCounter += openBraces - closesBraces
	if braceCounter == 0 && closesBraces > 0 {
		*extractedFunctions = append(*extractedFunctions,
			methodsToDoc.String())
		methodsToDoc.Reset()
		return false, 0
	}
	return true, braceCounter
}
