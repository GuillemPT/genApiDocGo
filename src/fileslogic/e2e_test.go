package fileslogic

import (
	"encoding/json"
	"genApiDocGo/src/internal"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/quick"
)

// TestE2E_FullPipeline exercises the complete pipeline:
// GetFiles → GetContent → FormatMethods → WriteDocument.
// It creates a temporary JS file with tagged routes and verifies
// that a valid swagger.json is produced.
func TestE2E_FullPipeline(t *testing.T) {
	internal.SetConfiguration("")

	dir := t.TempDir()
	jsContent := `
// @api_generate_doc
/* Get all users */
router.get('/users', async (req, res) => {
  res.status(200).json(users);
});

// @api_generate_doc
/* Create a user */
router.post('/users', async (req, res) => {
  res.status(201).json(newUser);
});

// @api_generate_doc
/* Get user by ID */
router.get('/users/:id', async (req, res) => {
  res.status(200).json(user);
});
`
	jsFile := filepath.Join(dir, "routes.js")
	if err := os.WriteFile(jsFile, []byte(jsContent), 0o600); err != nil {
		t.Fatalf("failed to write test JS file: %v", err)
	}

	files, _, err := GetFiles(dir, ".js")
	if err != nil {
		t.Fatalf("GetFiles failed: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	methods := GetContent(files)
	if len(methods) != 3 {
		t.Fatalf("expected 3 methods, got %d", len(methods))
	}

	structured := FormatMethods(methods)
	if len(structured) != 2 {
		t.Fatalf("expected 2 paths (/users and /users/{id}), got %d", len(structured))
	}

	outDir := t.TempDir()
	WriteDocument(structured, outDir)

	data, err := os.ReadFile(filepath.Join(outDir, "swagger.json"))
	if err != nil {
		t.Fatalf("swagger.json not created: %v", err)
	}

	var doc internal.SwaggerDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("swagger.json is not valid JSON: %v", err)
	}

	if doc.Openapi == "" {
		t.Error("swagger.json missing 'openapi' field")
	}
	if _, ok := doc.Paths["/users"]; !ok {
		t.Error("expected path /users in swagger.json")
	}
	if _, ok := doc.Paths["/users/{id}"]; !ok {
		t.Error("expected path /users/{id} in swagger.json")
	}

	getUsersOp, ok := doc.Paths["/users"]["get"]
	if !ok {
		t.Error("expected GET /users operation")
	}
	if getUsersOp.Responses == nil {
		t.Error("GET /users responses must not be nil")
	}
	if _, ok := getUsersOp.Responses["200"]; !ok {
		t.Error("expected response 200 for GET /users")
	}

	postUsersOp, ok := doc.Paths["/users"]["post"]
	if !ok {
		t.Error("expected POST /users operation")
	}
	if _, ok := postUsersOp.Responses["201"]; !ok {
		t.Error("expected response 201 for POST /users")
	}
}

// TestProperty_ExcludeFilesInBanDirectories_SubsetOfInput verifies that
// the result of ExcludeFilesInBanDirectories is always a subset of the
// original file list.
func TestProperty_ExcludeFilesInBanDirectories_SubsetOfInput(t *testing.T) {
	property := func(directories []string, files []string) bool {
		result := ExcludeFilesInBanDirectories(directories, files)
		fileSet := make(map[string]bool, len(files))
		for _, f := range files {
			fileSet[f] = true
		}
		for _, r := range result {
			if !fileSet[r] {
				return false
			}
		}
		return true
	}
	if err := quick.Check(property, nil); err != nil {
		t.Error(err)
	}
}

// TestProperty_ExcludeFilesInBanDirectories_NoBannedPaths verifies that
// no file in the result contains any of the banned directory paths.
func TestProperty_ExcludeFilesInBanDirectories_NoBannedPaths(t *testing.T) {
	property := func(directories []string, files []string) bool {
		result := ExcludeFilesInBanDirectories(directories, files)
		for _, f := range result {
			for _, dir := range directories {
				if dir != "" && strings.Contains(f, dir) {
					return false
				}
			}
		}
		return true
	}
	if err := quick.Check(property, nil); err != nil {
		t.Error(err)
	}
}

// TestProperty_ExcludeFilesInBanDirectories_LengthMonotone verifies that
// adding more banned directories never increases the result length.
func TestProperty_ExcludeFilesInBanDirectories_LengthMonotone(t *testing.T) {
	property := func(directories []string, extra string, files []string) bool {
		baseline := ExcludeFilesInBanDirectories(directories, files)
		extended := ExcludeFilesInBanDirectories(
			append(directories, extra), files)
		return len(extended) <= len(baseline)
	}
	if err := quick.Check(property, nil); err != nil {
		t.Error(err)
	}
}

// TestProperty_FormatMethods_ResponsesNeverNil verifies that every operation
// produced by FormatMethods always has a non-nil Responses map.
func TestProperty_FormatMethods_ResponsesNeverNil(t *testing.T) {
	httpMethods := []string{"get", "post", "put", "delete", "patch"}
	paths := []string{"/users", "/items", "/orders", "/health"}
	statusCodes := []string{"200", "201", "400", "404", "500"}

	property := func(methodIdx uint8, pathIdx uint8, statusIdx uint8) bool {
		httpMethod := httpMethods[int(methodIdx)%len(httpMethods)]
		path := paths[int(pathIdx)%len(paths)]
		status := statusCodes[int(statusIdx)%len(statusCodes)]

		snippet := "/* Test route */\n" +
			"router." + httpMethod + "('" + path + "', async (req, res) => {\n" +
			"  res.status(" + status + ").json({});\n" +
			"});"

		result := FormatMethods([]string{snippet})
		for _, pathDoc := range result {
			for _, op := range pathDoc {
				if op.Responses == nil {
					return false
				}
			}
		}
		return true
	}
	if err := quick.Check(property, nil); err != nil {
		t.Error(err)
	}
}
