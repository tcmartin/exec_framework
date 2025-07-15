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

	node := NewHTTPRequest("POST", server.URL+"/{{.path}}", `{"key":"{{.value}}"}`)

	ctx := &framework.Context{
		Ctx:        context.Background(),
		HTTPClient: retryablehttp.NewClient(),
		Env:        map[string]string{},
	}

	inputs := []map[string]interface{}{{"path": "test", "value": "testValue"}}
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

	node := NewHTTPRequest("GET", server.URL, "{}")

	ctx := &framework.Context{
		Ctx:        context.Background(),
		HTTPClient: retryablehttp.NewClient(),
		Env:        map[string]string{},
	}

	inputs := []map[string]interface{}{{}}
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

	node := NewHTTPRequest("GET", server.URL, "{}")

	ctx := &framework.Context{
		Ctx:        context.Background(),
		HTTPClient: retryablehttp.NewClient(),
		Env:        map[string]string{},
	}

	inputs := []map[string]interface{}{{}}
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