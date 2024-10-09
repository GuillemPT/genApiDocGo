package internal

const (
	// App Name.
	AppName = "GenApiDocGo"
	// Type of letters.
	LettersType = "slant"
	// View heigh for the selection view.
	ViewHeigh = 5
	// Tag to identify the methods to documentation.
	Tag = "@api_generate_doc"
	// The minium length for regex to identify the name and type of http call.
	RegexHeaderLength = 4
	// The exactly length for regex to identify the response number.
	RegexStatusCodeLength = 2
)

// Get file type options.
func GetFileTypeOptions() []string {
	return []string{".js", ".ts"}
}

// Frameworks for each file type.
func GetFrameworks(fileType string) []string {
	frameworksByFileType := map[string][]string{
		".js": {"Express"},
		".ts": {"Express"},
	}
	return frameworksByFileType[fileType]
}
