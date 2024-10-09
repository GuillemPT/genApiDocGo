package fileslogic

import (
	"io/fs"
	"testing"
)

type MockWalkDir struct {
	MockFunc func(root string, fn fs.WalkDirFunc) error
}

func (m *MockWalkDir) WalkDir(root string, fn fs.WalkDirFunc) error {
	return m.MockFunc(root, fn)
}

func TestGetFiles(t *testing.T) { //nolint: gocognit// it's a test.
	tests := []struct {
		name                string
		targetDirectoryPath string
		filesType           string
		expectedFiles       []string
		expectedDirectories []string
		// expectedError error
	}{
		{
			name:                "Find files with the extension",
			targetDirectoryPath: "./mockDirectory",
			filesType:           "js",
			expectedFiles:       []string{"./mockDirectory/mockFile.js"},
			expectedDirectories: []string{},
		},
		{
			name:                "Directory does not exist",
			targetDirectoryPath: "./nonExistentDirectory",
			filesType:           "js",
			expectedFiles:       nil,
			expectedDirectories: nil,
		},
		{
			name:                "Empty directory",
			targetDirectoryPath: "./emptyDirectory",
			filesType:           "js",
			expectedFiles:       []string{},
			expectedDirectories: []string{},
		},
		{
			name:                "Find multiple files with the extension",
			targetDirectoryPath: "./mockDirectory",
			filesType:           "js",
			expectedFiles: []string{"./mockDirectory/mockFile1.js",
				"./mockDirectory/mockFile2.js"},
			expectedDirectories: []string{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filesResult, directoriesResult, err := GetFiles(
				test.targetDirectoryPath, test.filesType)

			if test.name == "Directory does not exist" {
				if err == nil {
					t.Errorf("expected error, got none")
				}
				return
			}

			for i := range filesResult {
				if filesResult[i] != test.expectedFiles[i] {
					t.Errorf("at index %d: expected %v, got %v", i,
						test.expectedFiles[i], filesResult[i])
				}
			}

			for i := range directoriesResult {
				if directoriesResult[i] != test.expectedDirectories[i] {
					t.Errorf("at index %d: expected %v, got %v", i,
						test.expectedDirectories[i], directoriesResult[i])
				}
			}
		})
	}
}
func TestExcludeFilesInBanDirectories(t *testing.T) {
	tests := []struct {
		name        string
		directories []string
		files       []string
		expected    []string
	}{
		{
			name:        "No banned directories",
			directories: []string{},
			files:       []string{"file1.txt", "file2.txt"},
			expected:    []string{"file1.txt", "file2.txt"},
		},
		{
			name:        "One banned directory, exclude matching files",
			directories: []string{"/ban/"},
			files: []string{"/ban/file1.txt", "/allowed/file2.txt",
				"/ban/file3.txt"},
			expected: []string{"/allowed/file2.txt"},
		},
		{
			name:        "Multiple banned directories",
			directories: []string{"/ban1/", "/ban2/"},
			files: []string{"/ban1/file1.txt", "/ban2/file2.txt",
				"/allowed/file3.txt"},
			expected: []string{"/allowed/file3.txt"},
		},
		{
			name:        "No files in banned directories",
			directories: []string{"/ban/"},
			files:       []string{"/allowed/file1.txt", "/allowed/file2.txt"},
			expected:    []string{"/allowed/file1.txt", "/allowed/file2.txt"},
		},
		{
			name:        "Files partially matching banned directories",
			directories: []string{"/ban/"},
			files: []string{"/ban_hammer/file1.txt", "/ban/file2.txt",
				"/banish/file3.txt"},
			expected: []string{"/ban_hammer/file1.txt", "/banish/file3.txt"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ExcludeFilesInBanDirectories(test.directories,
				test.files)
			if len(result) != len(test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
			for i := range result {
				if result[i] != test.expected[i] {
					t.Errorf("at index %d: expected %v, got %v", i,
						test.expected[i], result[i])
				}
			}
		})
	}
}
