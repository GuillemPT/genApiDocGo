package fileslogic

import (
	"strings"
)

type MethodDoc struct {
	Route               string
	Type                string
	Description         string
	Response            int
	ResponseDescription string
}

func FormatMethods(methods []string) []MethodDoc {
	var formatMethodsStruct []MethodDoc
	for _, method := range methods {
		lines := strings.Split(method, "\n")
		var methodDocIter MethodDoc
		inDescription := false
		for _, line := range lines {
			if strings.Contains(line, "*/") {
				inDescription = false
			}
			if inDescription {
				description := line
				description = strings.Trim(description, " ")
				description = strings.Trim(description, "*")
				methodDocIter.Description += description
			}
			if strings.Contains(line, "/**") {
				inDescription = true
			}
			switch {
			case strings.Contains(line, "get"):
				methodDocIter.Type = "get"
				methodDocIter.Route = _getName(line)
				methodDocIter.Response = 200
				methodDocIter.ResponseDescription = "OK"
			case strings.Contains(line, "post"):
				methodDocIter.Type = "post"
				methodDocIter.Route = _getName(line)
				methodDocIter.Response = 200
				methodDocIter.ResponseDescription = "OK"
			case strings.Contains(line, "put"):
				methodDocIter.Type = "put"
				methodDocIter.Route = _getName(line)
				methodDocIter.Response = 200
				methodDocIter.ResponseDescription = "OK"
			case strings.Contains(line, "delete"):
				methodDocIter.Type = "delete"
				methodDocIter.Route = _getName(line)
				methodDocIter.Response = 200
				methodDocIter.ResponseDescription = "OK"
			case strings.Contains(line, "patch"):
				methodDocIter.Type = "patch"
				methodDocIter.Route = _getName(line)
				methodDocIter.Response = 200
				methodDocIter.ResponseDescription = "OK"
			}
		}

		formatMethodsStruct = append(formatMethodsStruct, methodDocIter)
	}
	return formatMethodsStruct
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
