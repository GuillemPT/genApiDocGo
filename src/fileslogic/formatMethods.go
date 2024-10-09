package fileslogic

import (
	"regexp"
	"strings"
)

const (
	regexHeaderLength     = 4
	regexStatusCodeLength = 2
)

var regex = regexp.MustCompile(
	`\w+\.(get|post|put|delete)\((["'])([^"']+)(["'])`)

var statusRegex = regexp.MustCompile(`res\.status\((\d{3})\)`)

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
	optDoc := operationDocument{Responses: make(map[string]responsesDocument)}
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
			if len(match) >= regexHeaderLength {
				result.pathName = match[3]
				result.operationName = match[1]
				pathDoc[result.operationName] = optDoc
			}

			statusMatch := statusRegex.FindStringSubmatch(line)
			if len(statusMatch) == regexStatusCodeLength {
				statusCode := statusMatch[1]
				optDoc.Responses[statusCode] = responsesDocument{
					Description: "Response status: " + statusCode,
				}
			}
		}
	}
	if len(optDoc.Responses) == 0 {
		optDoc.Responses["200"] = responsesDocument{
			Description: "Ok",
		}
	}

	pathDoc[result.operationName] = optDoc
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
