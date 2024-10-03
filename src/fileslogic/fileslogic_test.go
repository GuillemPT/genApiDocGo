package fileslogic_test

import (
	"genApiDocGo/src/fileslogic"
	"testing"
)

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
			result := fileslogic.ExcludeFilesInBanDirectories(test.directories,
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
