package fileslogic

import (
	"genApiDocGo/src/internal"
	"regexp"
	"strings"
)

// regex to extract the route name and the type of http call.
var regex = regexp.MustCompile(
	`\w+\.(get|post|put|delete)\((["'])([^"']+)(["'])`)

// regex to extract the request status.
var statusRegex = regexp.MustCompile(`res\.status\((\d{3})\)`)

// Process all the methods passed by argument and returns a map of type
// key: route type, value: PathDocument.
func FormatMethods(methods []string) map[string]PathDocument {
	pathDocMap := make(map[string]PathDocument, len(methods))
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

// Return one type instead of three elements.
type formatResult struct {
	pathName      string
	pathDoc       PathDocument
	operationName string
}

// Process each method individually and formats it.
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
			if len(match) >= internal.RegexHeaderLength {
				result.pathName = match[3]
				result.operationName = match[1]
				pathDoc[result.operationName] = optDoc
			}

			statusMatch := statusRegex.FindStringSubmatch(line)
			if len(statusMatch) == internal.RegexStatusCodeLength {
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

// Handle if exist a value in the map for the current key,
// and adds the new value.
func addElementInPathDoc(pathMap map[string]PathDocument,
	key string, pathDoc PathDocument, operationName string) {
	if currentPath, exists := pathMap[key]; exists {
		currentPath[operationName] = pathDoc[operationName]
	} else {
		pathMap[key] = pathDoc
	}
}
