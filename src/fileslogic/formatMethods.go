package fileslogic

import (
	"regexp"
	"strings"
)

const regexMinLength = 4

var regex = regexp.MustCompile(
	`\w+\.(get|post|put|delete)\((["'])([^"']+)(["'])`)

func FormatMethods(methods []string) map[string]pathDocument {
	pathDocMap := make(map[string]pathDocument, len(methods))
	for _, method := range methods {
		result := formatMethod(method)
		addElementInPathDoc(
			pathDocMap,
			result.pathName,
			result.pathDoc,
			result.operationName)
	}
	return pathDocMap
}

type formatResult struct {
	pathName      string
	pathDoc       pathDocument
	operationName string
}

func formatMethod(method string) formatResult {
	var result formatResult
	lines := strings.Split(method, "\n")
	pathDoc := make(map[string]operationDocument)
	optDoc := operationDocument{
		Responses: map[string]responsesDocument{
			"200": {Description: "Ok"},
		},
	}
	inDescription := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "/*") {
			inDescription = true
		}
		if inDescription {
			optDoc.Description += strings.Trim(line, "/*")
		}
		if strings.Contains(line, "*/") {
			inDescription = false
		}
		if !inDescription {
			match := regex.FindStringSubmatch(line)
			if len(match) >= regexMinLength {
				result.pathName = match[3]
				result.operationName = match[1]
				pathDoc[result.operationName] = optDoc
			}
		}
	}
	result.pathDoc = pathDoc
	return result
}

func addElementInPathDoc(pathMap map[string]pathDocument,
	key string, pathDoc pathDocument, operationName string) {
	if currentPath, exists := pathMap[key]; exists {
		currentPath[operationName] = pathDoc[operationName]
	} else {
		pathMap[key] = pathDoc
	}
}
