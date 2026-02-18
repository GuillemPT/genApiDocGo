package fileslogic

import (
	"genApiDocGo/src/internal"
	"strings"
	"testing"
)

func init() {
	internal.SetConfiguration("")
}

func TestFormatMethods_BasicGet(t *testing.T) {
	methods := []string{
		`/* Description */
router.get('/users', async (req, res) => {
  res.status(200).json(users);
});`,
	}
	result := FormatMethods(methods)

	path, ok := result["/users"]
	if !ok {
		t.Fatal("expected path /users in result")
	}
	op, ok := path["get"]
	if !ok {
		t.Fatal("expected 'get' operation for /users")
	}
	if op.Description != "Description" {
		t.Errorf("expected description 'Description', got '%s'", op.Description)
	}
	if _, ok := op.Responses["200"]; !ok {
		t.Error("expected response code 200")
	}
}

func TestFormatMethods_PatchMethod(t *testing.T) {
	methods := []string{
		`/* Update a user */
router.patch('/users/:id', async (req, res) => {
  res.status(200).json({});
});`,
	}
	result := FormatMethods(methods)

	path, ok := result["/users/{id}"]
	if !ok {
		t.Fatal("expected path /users/{id} in result")
	}
	if _, ok := path["patch"]; !ok {
		t.Error("expected 'patch' operation")
	}
}

func TestFormatMethods_RouteParameters(t *testing.T) {
	methods := []string{
		`/* Get user by ID */
router.get('/users/:id', async (req, res) => {
  res.status(200).json(user);
});`,
	}
	result := FormatMethods(methods)

	op := result["/users/{id}"]["get"]
	if len(op.Parameters) != 1 {
		t.Fatalf("expected 1 path parameter, got %d", len(op.Parameters))
	}
	param := op.Parameters[0]
	if param.Name != "id" {
		t.Errorf("expected parameter name 'id', got '%s'", param.Name)
	}
	if param.In != "path" {
		t.Errorf("expected parameter in 'path', got '%s'", param.In)
	}
	if !param.Required {
		t.Error("expected path parameter to be required")
	}
	if param.Schema.Type != "string" {
		t.Errorf("expected schema type 'string', got '%s'", param.Schema.Type)
	}
}

func TestFormatMethods_MultipleRouteParameters(t *testing.T) {
	methods := []string{
		`/* Get item in store */
router.get('/stores/:storeId/items/:itemId', async (req, res) => {
  res.status(200).json(item);
});`,
	}
	result := FormatMethods(methods)

	op := result["/stores/{storeId}/items/{itemId}"]["get"]
	if len(op.Parameters) != 2 {
		t.Fatalf("expected 2 path parameters, got %d", len(op.Parameters))
	}
	if op.Parameters[0].Name != "storeId" {
		t.Errorf("expected first param 'storeId', got '%s'", op.Parameters[0].Name)
	}
	if op.Parameters[1].Name != "itemId" {
		t.Errorf("expected second param 'itemId', got '%s'", op.Parameters[1].Name)
	}
}

func TestFormatMethods_SummaryAnnotation(t *testing.T) {
	methods := []string{
		`/*
 * @summary List all users
 * @tags users
 */
router.get('/users', async (req, res) => {
  res.status(200).json(users);
});`,
	}
	result := FormatMethods(methods)

	op := result["/users"]["get"]
	if op.Summary != "List all users" {
		t.Errorf("expected summary 'List all users', got '%s'", op.Summary)
	}
}

func TestFormatMethods_TagsAnnotation(t *testing.T) {
	methods := []string{
		`/*
 * @tags users, admin
 */
router.get('/users', async (req, res) => {
  res.status(200).json(users);
});`,
	}
	result := FormatMethods(methods)

	op := result["/users"]["get"]
	if len(op.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(op.Tags))
	}
	if op.Tags[0] != "users" {
		t.Errorf("expected first tag 'users', got '%s'", op.Tags[0])
	}
	if op.Tags[1] != "admin" {
		t.Errorf("expected second tag 'admin', got '%s'", op.Tags[1])
	}
}

func TestFormatMethods_DefaultResponse(t *testing.T) {
	methods := []string{
		`/* No status call */
router.get('/health', async (req, res) => {
  res.json({ status: 'ok' });
});`,
	}
	result := FormatMethods(methods)

	op := result["/health"]["get"]
	if _, ok := op.Responses["200"]; !ok {
		t.Error("expected default 200 response when no res.status() call found")
	}
}

func TestFormatMethods_MultipleStatusCodes(t *testing.T) {
	methods := []string{
		`/* Create user */
router.post('/users', async (req, res) => {
  if (!req.body.name) {
    res.status(400).json({ error: 'Name required' });
    return;
  }
  res.status(201).json(newUser);
});`,
	}
	result := FormatMethods(methods)

	op := result["/users"]["post"]
	if _, ok := op.Responses["201"]; !ok {
		t.Error("expected response code 201")
	}
	if _, ok := op.Responses["400"]; !ok {
		t.Error("expected response code 400")
	}
}

func TestFormatMethods_SamePathMultipleMethods(t *testing.T) {
	methods := []string{
		`/* List users */
router.get('/users', async (req, res) => {
  res.status(200).json(users);
});`,
		`/* Create user */
router.post('/users', async (req, res) => {
  res.status(201).json(newUser);
});`,
	}
	result := FormatMethods(methods)

	path, ok := result["/users"]
	if !ok {
		t.Fatal("expected path /users in result")
	}
	if _, ok := path["get"]; !ok {
		t.Error("expected 'get' operation for /users")
	}
	if _, ok := path["post"]; !ok {
		t.Error("expected 'post' operation for /users")
	}
}

func TestFormatMethods_AnnotationsNotInDescription(t *testing.T) {
	methods := []string{
		`/*
 * @summary Create a user
 * @tags users
 */
router.post('/users', async (req, res) => {
  res.status(201).json(newUser);
});`,
	}
	result := FormatMethods(methods)

	op := result["/users"]["post"]
	if strings.Contains(op.Description, "@summary") {
		t.Errorf("description should not contain @summary annotation, got: %q", op.Description)
	}
	if strings.Contains(op.Description, "@tags") {
		t.Errorf("description should not contain @tags annotation, got: %q", op.Description)
	}
	if op.Summary != "Create a user" {
		t.Errorf("expected summary 'Create a user', got %q", op.Summary)
	}
}

func TestFormatMethods_PathNormalization(t *testing.T) {
	methods := []string{
		`/* Get item */
router.get('/stores/:storeId/items/:itemId', async (req, res) => {
  res.status(200).json(item);
});`,
	}
	result := FormatMethods(methods)

	if _, ok := result["/stores/{storeId}/items/{itemId}"]; !ok {
		t.Error("expected OpenAPI-normalized path /stores/{storeId}/items/{itemId}")
	}
	if _, ok := result["/stores/:storeId/items/:itemId"]; ok {
		t.Error("Express-style path /stores/:storeId/items/:itemId should not appear in output")
	}
}

func TestFormatMethods_DoubleQuotedPath(t *testing.T) {
	methods := []string{
		`/* Health check */
router.get("/health", async (req, res) => {
  res.status(200).json({ status: 'ok' });
});`,
	}
	result := FormatMethods(methods)

	if _, ok := result["/health"]; !ok {
		t.Error("expected path /health when using double quotes in route")
	}
}
