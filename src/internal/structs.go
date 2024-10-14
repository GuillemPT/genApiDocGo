package internal

type SwaggerDocument struct {
	Openapi string       `json:"openapi"`
	Info    InfoDocument `json:"info"`
	// The same path can have different operations.
	Paths map[string]PathDocument `json:"paths"` // key is the route (/home)
}

type InfoDocument struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// key is the type of operation (get).
type PathDocument map[string]OperationDocument

type OperationDocument struct {
	Description string `json:"description"`
	// key is the code id of response (200).
	Responses map[string]ResponsesDocument `json:"responses"`
}

type ResponsesDocument struct {
	Description string `json:"description"`
}

type configFile struct {
	ResponsesMap map[string]string `json:"responsesMap"`
	BaseDocument SwaggerDocument   `json:"baseDocument"`
}
