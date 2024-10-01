package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func getFiles(targetDirectoryPath string,
	filesType string) ([]string, []*item, error) {
	var files []string
	var directories []*item

	err := filepath.WalkDir(targetDirectoryPath,
		func(path string, file fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if file.IsDir() {
				directory := item{ID: file.Name()}
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

func excludeFilesInBanDirectories(directories []string,
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

	fmt.Println(tmpFiles)
	return tmpFiles
}
