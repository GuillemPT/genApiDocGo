package fileslogic

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Item struct serves to make a multi selector in getExcludedDirectories
// (genApiDocGo/src/cli.go).
type Item struct {
	ID         string
	IsSelected bool
}

// GetFiles returns files ([]string) that its extension match with fileType
// parameter, directories ([]*item) and error.
func GetFiles(targetDirectoryPath string,
	filesType string) ([]string, []*Item, error) {
	var files []string
	var directories []*Item

	err := filepath.WalkDir(targetDirectoryPath,
		func(path string, file fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if file.IsDir() {
				directory := Item{ID: file.Name()}
				directories = append(directories, &directory)
			}

			if strings.Contains(file.Name(), ".") {
				if strings.Split(file.Name(), ".")[1] == filesType {
					files = append(files, path)
				}
			}

			return nil
		})

	if err != nil {
		return nil, nil, err
	}

	return files, directories, nil
}

func ExcludeFilesInBanDirectories(directories []string,
	files []string) []string {
	var tmpFiles []string
	var inBanDirectory bool

	for _, file := range files {
		for _, directory := range directories {
			if strings.Contains(file, directory) {
				inBanDirectory = true
			}
		}
		if !inBanDirectory {
			tmpFiles = append(tmpFiles, file)
		}
		inBanDirectory = false
	}

	return tmpFiles
}

// return array with each method to create documentation.
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

func processFile(scanner *bufio.Scanner, extractedFunctions *[]string) {
	inMethod := false
	braceCounter := -1
	var methodsToDoc strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if inMethod {
			inMethod, braceCounter = processMethod(line, braceCounter,
				&methodsToDoc, extractedFunctions)
		}
		if strings.Contains(line, "@api_generate_doc") {
			inMethod = true
		}
	}
}

// return the bool value, true if have to process the next line or false if not
// and the current counter of brace.
func processMethod(line string, braceCounter int, methodsToDoc *strings.Builder,
	extractedFunctions *[]string) (bool, int) {
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
		*extractedFunctions = append(*extractedFunctions,
			methodsToDoc.String())
		methodsToDoc.Reset()
		return false, -1
	}
	return true, braceCounter
}
