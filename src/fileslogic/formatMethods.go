package fileslogic

import (
	"strings"
)

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
			switch {
			case strings.Contains(line, ".get("):
				pathName = _getName(line)
				responsesDoc.Description = "Ok"
				optDoc.Responses["200"] = responsesDoc
				pathDoc["get"] = optDoc
			case strings.Contains(line, ".post("):
				pathName = _getName(line)
				responsesDoc.Description = "Ok"
				optDoc.Responses["200"] = responsesDoc
				pathDoc["post"] = optDoc
			case strings.Contains(line, ".put("):
				pathName = _getName(line)
				responsesDoc.Description = "Ok"
				optDoc.Responses["200"] = responsesDoc
				pathDoc["put"] = optDoc
			case strings.Contains(line, ".delete("):
				pathName = _getName(line)
				responsesDoc.Description = "Ok"
				optDoc.Responses["200"] = responsesDoc
				pathDoc["delete"] = optDoc
			}
		}
	}
	return pathName, pathDoc
}

func _getName(line string) string {
	start := strings.Index(line, "(")
	end := strings.Index(line, ")")

	partialLine := line[start+1 : end]
	parts := strings.Split(partialLine, ",")

	route := strings.TrimSpace(parts[0])
	route = strings.Trim(route, "'")
	return route
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
