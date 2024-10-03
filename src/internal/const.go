package internal

func GetFileTypeOptions() []string {
	return []string{"js", "ts"}
}

func GetFrameworks(fileType string) []string {
	frameworksByFileType := map[string][]string{
		"js": {"Express"},
		"ts": {"Express"},
	}
	return frameworksByFileType[fileType]
}
