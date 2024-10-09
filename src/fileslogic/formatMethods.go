package fileslogic

import (
	"regexp"
	"strings"
)

const regexMinLength = 4

func FormatMethods(methods []string) map[string]pathDocument {
	pathDocMap := make(map[string]pathDocument)
	for _, method := range methods {
		pathName, pathDoc := formatMethod(method)
		addElementInPathDoc(pathDocMap, pathName, pathDoc)
	}
	return pathDocMap
}

func formatMethod(method string) (string, pathDocument) {
	var pathName string
	lines := strings.Split(method, "\n")
	pathDoc := make(map[string]operationDocument)
	var optDoc operationDocument
	optDoc.Responses = make(map[string]responsesDocument)
	var responsesDoc responsesDocument
	inDescription := false
	// Regular expression to match HTTP method calls and extract path name
	regex := regexp.MustCompile(
		`\w+\.(get|post|put|delete)\((["'])([^"']+)(["'])`)

	for _, line := range lines {
		line = strings.Trim(line, " ")
		if strings.Contains(line, "/*") {
			inDescription = true
		}
		if inDescription {
			description := line
			description = strings.Trim(description, "/")
			description = strings.Trim(description, "*")
			optDoc.Description += description
		}
		if strings.Contains(line, "*/") {
			inDescription = false
		}

		if !inDescription {
			match := regex.FindStringSubmatch(line)
			if len(match) >= regexMinLength {
				// route name
				pathName = match[3]
				responsesDoc.Description = "Ok"
				optDoc.Responses["200"] = responsesDoc
				// method name
				pathDoc[match[1]] = optDoc
			}
		}
	}
	return pathName, pathDoc
}

func addElementInPathDoc(pathMap map[string]pathDocument, key string,
	pathDoc pathDocument) {
	if pathMap[key] != nil {
		pathCurrent := pathMap[key]
		newKey, newValue := getOperation(pathDoc)
		pathCurrent[newKey] = newValue
		pathMap[key] = pathCurrent
	} else {
		pathMap[key] = pathDoc
	}
}

func getOperation(pathDoc pathDocument) (string, operationDocument) {
	var found bool
	var pathDocValue operationDocument
	pathDocKey := -1
	operations := []string{"get", "post", "delete", "update", "put"}
	for !found {
		pathDocKey++
		pathDocValue, found = pathDoc[operations[pathDocKey]]
	}
	return operations[pathDocKey], pathDocValue
}
