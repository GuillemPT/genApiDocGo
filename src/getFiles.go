package main

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func getFiles(targetDirectoryPath string, filesType string, excludeDirectories []string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(targetDirectoryPath, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(file.Name(), ".") && !_pathContainsAny(path, excludeDirectories){
			if strings.Split(file.Name(), ".")[1] == filesType {
				files = append(files, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}

func _pathContainsAny(path string, directories []string) bool {
	for _, dir := range directories {
		if strings.Contains(path, dir) {
			return true
		}
	}
	return false
}
