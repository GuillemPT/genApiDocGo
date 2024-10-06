package fileslogic

import (
	"strings"
)

// TODO: change all method, don't work well, do the test before please
func FormatMethods(methods []string) map[string]pathDocument {
	pathDocMap := make(map[string]pathDocument)
	for _, method := range methods {
		lines := strings.Split(method, "\n")
		pathDoc := make(map[string]operationDocument)
		var optDoc operationDocument
		optDoc.Responses = make(map[string]responsesDocument)
		var responsesDoc responsesDocument
		inDescription := false
		for _, line := range lines {
			if strings.Contains(line, "*/") {
				inDescription = false
			}
			if inDescription {
				description := line
				description = strings.Trim(description, " ")
				description = strings.Trim(description, "*")
				optDoc.Description += description
			}
			if strings.Contains(line, "/**") {
				inDescription = true
			}
			switch {
			case strings.Contains(line, ".get"):
				responsesDoc.Description = "Ok"
				optDoc.Responses["200"] = responsesDoc
				pathDoc["get"] = optDoc
				pathDocMap[_getName(line)] = pathDoc
			case strings.Contains(line, ".post"):
				responsesDoc.Description = "Ok"
				optDoc.Responses["200"] = responsesDoc
				pathDoc["post"] = optDoc
				pathDocMap[_getName(line)] = pathDoc
			case strings.Contains(line, ".put"):
				responsesDoc.Description = "Ok"
				optDoc.Responses["200"] = responsesDoc
				pathDoc["put"] = optDoc
				pathDocMap[_getName(line)] = pathDoc
			case strings.Contains(line, ".delete"):
				responsesDoc.Description = "Ok"
				optDoc.Responses["200"] = responsesDoc
				pathDoc["delete"] = optDoc
				pathDocMap[_getName(line)] = pathDoc
			case strings.Contains(line, ".patch"):
				responsesDoc.Description = "Ok"
				optDoc.Responses["200"] = responsesDoc
				pathDoc["patch"] = optDoc
				pathDocMap[_getName(line)] = pathDoc
			}
		}
	}
	return pathDocMap
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
