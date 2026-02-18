package fileslogic

import (
	"fmt"
	"genApiDocGo/src/internal"
	"regexp"
	"sort"
	"strings"
)

// regex to extract the route name and the type of http call (includes patch).
var regex = regexp.MustCompile(
	`\w+\.(get|post|put|delete|patch)\((["'])([^"']+)(["'])`)

// regex to extract path parameters from a route string (e.g. :id, :userId).
var pathParamRegex = regexp.MustCompile(`:(\w+)`)

// regex to extract @summary annotation from doc comments.
var summaryAnnotationRegex = regexp.MustCompile(`@summary\s+([^@]+)`)

// regex to extract @tags annotation from doc comments.
var tagsAnnotationRegex = regexp.MustCompile(`@tags\s+([^@]+)`)

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

// normalizePathParams converts Express-style path parameters (e.g. :id) to
// OpenAPI-compliant template syntax (e.g. {id}).
func normalizePathParams(path string) string {
	return pathParamRegex.ReplaceAllString(path, `{$1}`)
}

// buildStatusRegex builds the status-code matching regexp from the provided
// map in a deterministic order, quoting each key so custom configs with
// special characters cannot break the regexp.
func buildStatusRegex(m map[string]string) *regexp.Regexp {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		keys[i] = regexp.QuoteMeta(k)
	}
	return regexp.MustCompile(fmt.Sprintf(`res\.status\((%s)\)`,
		strings.Join(keys, "|")))
}

// Process each method individually and formats it.
func formatMethod(method string) formatResult {
	responsesConfig = internal.GetResponsesConfig()
	statusRegex = buildStatusRegex(responsesConfig)
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
			isAnnotation := false
			if match := summaryAnnotationRegex.FindStringSubmatch(line); len(match) >= 2 {
				optDoc.Summary = strings.TrimSpace(match[1])
				isAnnotation = true
			}
			if match := tagsAnnotationRegex.FindStringSubmatch(line); len(match) >= 2 {
				tags := strings.Split(strings.TrimSpace(match[1]), ",")
				for i := range tags {
					tags[i] = strings.TrimSpace(tags[i])
				}
				optDoc.Tags = tags
				isAnnotation = true
			}
			if !isAnnotation {
				optDoc.Description += strings.Trim(line, "/*")
			}
		}
		if strings.Contains(line, "*/") {
			inDescription = false
			optDoc.Description = strings.Trim(optDoc.Description, " ")
		}
		if !inDescription {
			match := regex.FindStringSubmatch(line)
			if len(match) >= internal.RegexHeaderLength {
				expressPath := match[3]
				result.pathName = normalizePathParams(expressPath)
				result.operationName = match[1]
				optDoc.Parameters = extractPathParameters(expressPath)
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

// extractPathParameters parses a route path and returns a ParameterDocument
// for each Express path parameter (e.g. :id → {name:"id", in:"path", required:true}).
func extractPathParameters(path string) []internal.ParameterDocument {
	matches := pathParamRegex.FindAllStringSubmatch(path, -1)
	if len(matches) == 0 {
		return nil
	}
	params := make([]internal.ParameterDocument, 0, len(matches))
	for _, m := range matches {
		params = append(params, internal.ParameterDocument{
			Name:     m[1],
			In:       "path",
			Required: true,
			Schema:   internal.SchemaDocument{Type: "string"},
		})
	}
	return params
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
