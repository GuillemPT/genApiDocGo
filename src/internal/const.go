package internal

const (
	// App Name.
	AppName = "GenApiDocGo"
	// Type of letters.
	LettersType = "slant"
	// View heigh for the selection view.
	ViewHeigh = 5
	// Tag to identify the methods to documentation.
	Tag = "@api_generate_doc"
	// The minium length for regex to identify the name and type of http call.
	RegexHeaderLength = 4
	// The exactly length for regex to identify the response number.
	RegexStatusCodeLength = 2
	// Default value for fileType.
	DefaultFileType = ".js"
	// Default value for Frameworks, first position of the array.
	DefaultFramework = 0
)

var config configFile //nolint: gochecknoglobals //config variable

// Get file type options.
func GetFileTypeOptions() []string {
	return []string{".js", ".ts"}
}

// Frameworks for each file type.
func GetFrameworks(fileType string) []string {
	frameworksByFileType := map[string][]string{
		".js": {"Express"},
		".ts": {"Express"},
	}
	return frameworksByFileType[fileType]
}

func getDefaultResponses() map[string]string {
	return map[string]string{
		"100": "100",
		"101": "101",
		"102": "102",
		"200": "200",
		"201": "201",
		"202": "202",
		"203": "203",
		"204": "204",
		"205": "205",
		"206": "206",
		"207": "207",
		"208": "208",
		"226": "226",
		"300": "300",
		"301": "301",
		"302": "302",
		"303": "303",
		"304": "304",
		"305": "305",
		"307": "307",
		"308": "308",
		"400": "400",
		"401": "401",
		"402": "402",
		"403": "403",
		"404": "404",
		"405": "405",
		"406": "406",
		"407": "407",
		"408": "408",
		"409": "409",
		"410": "410",
		"411": "411",
		"412": "412",
		"413": "413",
		"414": "414",
		"415": "415",
		"416": "416",
		"417": "417",
		"418": "418",
		"421": "421",
		"422": "422",
		"423": "423",
		"424": "424",
		"425": "425",
		"426": "426",
		"428": "428",
		"429": "429",
		"431": "431",
		"451": "451",
		"500": "500",
		"501": "501",
		"502": "502",
		"503": "503",
		"504": "504",
		"505": "505",
		"506": "506",
		"507": "507",
		"508": "508",
		"510": "510",
		"511": "511",
	}
}

func getDefaultBaseDocument() SwaggerDocument {
	return SwaggerDocument{
		Openapi: "3.0.0",
		Info: InfoDocument{
			Title:       "Title",
			Description: "Description",
			Version:     "1.0.0",
		},
	}
}

func GetResponsesConfig() map[string]string {
	return config.ResponsesMap
}

func GetBaseDocumentConfig() SwaggerDocument {
	return config.BaseDocument
}
