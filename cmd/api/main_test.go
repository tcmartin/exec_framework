package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-workflow/pkg/store"

	"github.com/gorilla/mux"
)

// initTestStore initializes a new in-memory SQLite store for testing.
func initTestStore() store.WorkflowStore {
	dbPath := "file::memory:?cache=shared" // In-memory SQLite
	store := store.NewSQLiteStore(dbPath)
	if err := store.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize test store: %v", err))
	}
	return store
}

func TestCreateWorkflowHandler(t *testing.T) {
	workflowStore = initTestStore()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/workflows", createWorkflowHandler).Methods("POST")

	// Test case 1: Valid workflow creation
	validWorkflow := WorkflowRequest{
		Name:       "test_workflow_1",
		Definition: "workflow_def_1",
	}
	body, _ := json.Marshal(validWorkflow)
	req := httptest.NewRequest("POST", "/api/v1/workflows", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var res WorkflowResponse
	json.NewDecoder(rr.Body).Decode(&res)
	if res.Name != validWorkflow.Name || res.ID == "" {
		t.Errorf("handler returned unexpected body: %v", res)
	}

	// Verify workflow is stored
	_, err := workflowStore.GetWorkflow(res.ID)
	if err != nil {
		t.Errorf("workflow not found in store after creation: %v", err)
	}

	// Test case 2: Missing name
	missingNameWorkflow := WorkflowRequest{
		Definition: "workflow_def_2",
	}
	body, _ = json.Marshal(missingNameWorkflow)
	req = httptest.NewRequest("POST", "/api/v1/workflows", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for missing name: got %v want %v", status, http.StatusBadRequest)
	}

	// Test case 3: Duplicate name
	duplicateNameWorkflow := WorkflowRequest{
		Name:       "test_workflow_1",
		Definition: "workflow_def_duplicate",
	}
	body, _ = json.Marshal(duplicateNameWorkflow)
	req = httptest.NewRequest("POST", "/api/v1/workflows", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code for duplicate name: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestGetWorkflowHandler(t *testing.T) {
	workflowStore = initTestStore()

	// Save a workflow first
	wf := &store.Workflow{
		ID:         "test_id_123",
		Name:       "get_test_workflow",
		Definition: "get_workflow_def",
	}
	if err := workflowStore.SaveWorkflow(wf); err != nil {
		t.Fatalf("Failed to save workflow for get test: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/workflows/{id}", getWorkflowHandler).Methods("GET")

	// Test case 1: Valid workflow ID
	req := httptest.NewRequest("GET", "/api/v1/workflows/test_id_123", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var res map[string]string
	json.NewDecoder(rr.Body).Decode(&res)
	if res["id"] != wf.ID || res["name"] != wf.Name || res["definition"] != wf.Definition {
		t.Errorf("handler returned unexpected body: got %v want %v", res, wf)
	}

	// Test case 2: Non-existent workflow ID
	req = httptest.NewRequest("GET", "/api/v1/workflows/non_existent_id", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent ID: got %v want %v", status, http.StatusNotFound)
	}
}

func TestRunWorkflowHandler(t *testing.T) {
	workflowStore = initTestStore()

	// Save a workflow first
	wf := &store.Workflow{
		ID:   "run_test_id",
		Name: "run_test_workflow",
		// A simple workflow that uses ManualTrigger and SetNode
		Definition: `
nodes:
  manualTrigger:
    type: manualTrigger
    payload:
      - message: "hello"
  setNode:
    type: setNode
    setValues:
      status: "processed"
connections:
  manualTrigger: [setNode]
`,
	}
	if err := workflowStore.SaveWorkflow(wf); err != nil {
		t.Fatalf("Failed to save workflow for run test: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/workflows/{id}/run", runWorkflowHandler).Methods("POST")

	// Test case 1: Valid workflow run with input
	inputPayload := []map[string]interface{}{
		{"initial": "data"},
	}
	body, _ := json.Marshal(inputPayload)
	req := httptest.NewRequest("POST", "/api/v1/workflows/run_test_id/run", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}

	var res map[string]string
	json.NewDecoder(rr.Body).Decode(&res)
	if res["message"] != "Workflow triggered successfully" || res["workflow_run_id"] == "" {
		t.Errorf("handler returned unexpected body: %v", res)
	}

	// Test case 2: Non-existent workflow ID
	req = httptest.NewRequest("POST", "/api/v1/workflows/non_existent_id/run", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent ID: got %v want %v", status, http.StatusNotFound)
	}

	// Test case 3: Invalid input payload
	req = httptest.NewRequest("POST", "/api/v1/workflows/run_test_id/run", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for invalid payload: got %v want %v", status, http.StatusBadRequest)
	}
}