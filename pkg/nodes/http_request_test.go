package nodes

import (
	"context"
	"encoding/json"
	"go-workflow/pkg/framework"
	"net/http"
	"net/http/httptest"
	"testing"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

func TestHTTPRequest_Execute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test" {
			t.Errorf("expected path /test, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		body, _ := json.Marshal(map[string]string{"status": "ok"})
		w.Write(body)
	}))
	defer server.Close()

	node := NewHTTPRequest("url", "method", "headers", "body")

	ctx := &framework.Context{
		Ctx:        context.Background(),
		HTTPClient: retryablehttp.NewClient(),
		Env:        map[string]string{},
	}

	inputs := []map[string]interface{}{
		{
			"url":    server.URL + "/test",
			"method": "POST",
			"headers": map[string]interface{}{"X-Test-Header": "test-value"},
			"body":    map[string]string{"key": "testValue"},
		},
	}
	out, err := node.Execute(ctx, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("expected 1 output, got %d", len(out))
	}
	if out[0]["status"] != "ok" {
		t.Errorf("expected status ok, got %v", out[0]["status"])
	}
}

func TestHTTPRequest_Execute_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	node := NewHTTPRequest("url", "method", "", "")

	ctx := &framework.Context{
		Ctx:        context.Background(),
		HTTPClient: retryablehttp.NewClient(),
		Env:        map[string]string{},
	}

	inputs := []map[string]interface{}{
		{
			"url":    server.URL,
			"method": "GET",
		},
	}
	_, err := node.Execute(ctx, inputs)
	if err == nil {
		t.Fatal("expected an error, but got nil")
	}
}

func TestHTTPRequest_Execute_ItemsResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := json.Marshal(map[string]interface{}{"items": []map[string]string{{"item1": "value1"}, {"item2": "value2"}}})
		w.Write(body)
	}))
	defer server.Close()

	node := NewHTTPRequest("url", "method", "", "")

	ctx := &framework.Context{
		Ctx:        context.Background(),
		HTTPClient: retryablehttp.NewClient(),
		Env:        map[string]string{},
	}

	inputs := []map[string]interface{}{
		{
			"url":    server.URL,
			"method": "GET",
		},
	}
	out, err := node.Execute(ctx, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out) != 2 {
		t.Fatalf("expected 2 outputs, got %d", len(out))
	}
	if out[0]["item1"] != "value1" || out[1]["item2"] != "value2" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestHTTPRequest_Execute_MissingURLKey(t *testing.T) {
	node := NewHTTPRequest("url", "method", "", "")
	ctx := &framework.Context{Ctx: context.Background(), HTTPClient: retryablehttp.NewClient()}
	inputs := []map[string]interface{}{{"method": "GET"}}
	_, err := node.Execute(ctx, inputs)
	if err == nil || err.Error() != "URL not found or not a string in input record for key url" {
		t.Errorf("expected error for missing URLKey, got %v", err)
	}
}

func TestHTTPRequest_Execute_InvalidURLType(t *testing.T) {
	node := NewHTTPRequest("url", "method", "", "")
	ctx := &framework.Context{Ctx: context.Background(), HTTPClient: retryablehttp.NewClient()}
	inputs := []map[string]interface{}{{"url": 123, "method": "GET"}}
	_, err := node.Execute(ctx, inputs)
	if err == nil || err.Error() != "URL not found or not a string in input record for key url" {
		t.Errorf("expected error for invalid URL type, got %v", err)
	}
}

func TestHTTPRequest_Execute_MissingMethodKey(t *testing.T) {
	node := NewHTTPRequest("url", "method", "", "")
	ctx := &framework.Context{Ctx: context.Background(), HTTPClient: retryablehttp.NewClient()}
	inputs := []map[string]interface{}{{"url": "http://example.com"}}
	_, err := node.Execute(ctx, inputs)
	if err == nil || err.Error() != "Method not found or not a string in input record for key method" {
		t.Errorf("expected error for missing MethodKey, got %v", err)
	}
}

func TestHTTPRequest_Execute_InvalidMethodType(t *testing.T) {
	node := NewHTTPRequest("url", "method", "", "")
	ctx := &framework.Context{Ctx: context.Background(), HTTPClient: retryablehttp.NewClient()}
	inputs := []map[string]interface{}{{"url": "http://example.com", "method": 123}}
	_, err := node.Execute(ctx, inputs)
	if err == nil || err.Error() != "Method not found or not a string in input record for key method" {
		t.Errorf("expected error for invalid Method type, got %v", err)
	}
}

func TestHTTPRequest_Execute_InvalidHeadersType(t *testing.T) {
	node := NewHTTPRequest("url", "method", "headers", "")
	ctx := &framework.Context{Ctx: context.Background(), HTTPClient: retryablehttp.NewClient()}
	inputs := []map[string]interface{}{{"url": "http://example.com", "method": "GET", "headers": "invalid"}}
	_, err := node.Execute(ctx, inputs)
	if err == nil || err.Error() != "Headers not found or not a map in input record for key headers" {
		t.Errorf("expected error for invalid Headers type, got %v", err)
	}
}

func TestHTTPRequest_Execute_MissingBodyKey(t *testing.T) {
	node := NewHTTPRequest("url", "method", "", "body")
	ctx := &framework.Context{Ctx: context.Background(), HTTPClient: retryablehttp.NewClient()}
	inputs := []map[string]interface{}{{"url": "http://example.com", "method": "POST"}}
	_, err := node.Execute(ctx, inputs)
	if err == nil || err.Error() != "Body not found in input record for key body" {
		t.Errorf("expected error for missing BodyKey, got %v", err)
	}
}

func TestHTTPRequest_Execute_UnmarshallableBody(t *testing.T) {
	node := NewHTTPRequest("url", "method", "", "body")
	ctx := &framework.Context{Ctx: context.Background(), HTTPClient: retryablehttp.NewClient()}
	inputs := []map[string]interface{}{
		{
			"url":    "http://example.com",
			"method": "POST",
			"body":    make(chan int), // Channels cannot be marshalled to JSON
		},
	}
	_, err := node.Execute(ctx, inputs)
	if err == nil || err.Error() != "failed to marshal body content: json: unsupported type: chan int" {
		t.Errorf("expected error for unmarshallable body, got %v", err)
	}
}