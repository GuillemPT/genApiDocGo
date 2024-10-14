package fileslogic

import (
	"fmt"
	"genApiDocGo/src/internal"
	"regexp"
	"strings"
)

// regex to extract the route name and the type of http call.
var regex = regexp.MustCompile(
	`\w+\.(get|post|put|delete)\((["'])([^"']+)(["'])`)

var responsesConfig map[string]string //nolint: gochecknoglobals // i need
// regex to extract the request status.
var statusRegex *regexp.Regexp //nolint: gochecknoglobals // i need

// Process all the methods passed by argument and returns a map of type
// key: route type, value: PathDocument.
func FormatMethods(methods []string) map[string]internal.PathDocument {
	pathDocMap := make(map[string]internal.PathDocument, len(methods))
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
	pathDoc       internal.PathDocument
	operationName string
}

// Process each method individually and formats it.
func formatMethod(method string) formatResult {
	responsesConfig = internal.GetResponsesConfig()
	statusRegex = regexp.MustCompile(fmt.Sprintf(`res\.status\((.{%s})\)`,
		strings.Join(getKeys(responsesConfig), "|")))
	var result formatResult
	lines := strings.Split(method, "\n")
	pathDoc := make(map[string]internal.OperationDocument)
	optDoc := internal.OperationDocument{Responses: make(
		map[string]internal.ResponsesDocument)}
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
			optDoc.Description = strings.Trim(optDoc.Description, " ")
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
				code := responsesConfig[statusCode]
				optDoc.Responses[code] = internal.ResponsesDocument{
					Description: "Response status: " + code,
				}
			}
		}
	}
	if len(optDoc.Responses) == 0 {
		optDoc.Responses["200"] = internal.ResponsesDocument{
			Description: "Ok",
		}
	}

	pathDoc[result.operationName] = optDoc
	result.pathDoc = pathDoc
	return result
}

// Handle if exist a value in the map for the current key,
// and adds the new value.
func addElementInPathDoc(pathMap map[string]internal.PathDocument,
	key string, pathDoc internal.PathDocument, operationName string) {
	if currentPath, exists := pathMap[key]; exists {
		currentPath[operationName] = pathDoc[operationName]
	} else {
		pathMap[key] = pathDoc
	}
}

// Function to get keys of the map, may be in the future can go to some internal
// file, but for now it is only used here.
func getKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
