package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-workflow/pkg/framework"
	_ "go-workflow/internal/noderegistry" // Import for side effect of registering nodes
	"go-workflow/pkg/store"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// WorkflowRequest represents the request body for uploading a workflow.
type WorkflowRequest struct {
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

// WorkflowResponse represents the response body for a workflow.
type WorkflowResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// APIError represents a generic API error response.
type APIError struct {
	Message string `json:"message"`
}

var workflowStore store.WorkflowStore

func main() {
	// Initialize workflow store
	dbPath := os.Getenv("WORKFLOW_DB_PATH")
	if dbPath == "" {
		dbPath = "workflows.db"
	}
	workflowStore = store.NewSQLiteStore(dbPath)
	if err := workflowStore.Init(); err != nil {
		log.Fatalf("Failed to initialize workflow store: %v", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/workflows", createWorkflowHandler).Methods("POST")
	router.HandleFunc("/api/v1/workflows/{id}", getWorkflowHandler).Methods("GET")
	router.HandleFunc("/api/v1/workflows/{id}/run", runWorkflowHandler).Methods("POST")

	log.Printf("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func createWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	var req WorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, jsonError("Invalid request body"), http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Definition == "" {
		http.Error(w, jsonError("Name and definition are required"), http.StatusBadRequest)
		return
	}

	workflow := &store.Workflow{
		ID:         uuid.New().String(),
		Name:       req.Name,
		Definition: req.Definition,
	}

	if err := workflowStore.SaveWorkflow(workflow); err != nil {
		http.Error(w, jsonError(fmt.Sprintf("Failed to save workflow: %v", err)), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(WorkflowResponse{ID: workflow.ID, Name: workflow.Name})
}

func getWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	workflow, err := workflowStore.GetWorkflow(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, jsonError("Workflow not found"), http.StatusNotFound)
		} else {
			http.Error(w, jsonError(fmt.Sprintf("Failed to retrieve workflow: %v", err)), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":         workflow.ID,
		"name":       workflow.Name,
		"definition": workflow.Definition,
	})
}

func runWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	storedWorkflow, err := workflowStore.GetWorkflow(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, jsonError("Workflow not found"), http.StatusNotFound)
		} else {
			http.Error(w, jsonError(fmt.Sprintf("Failed to retrieve workflow: %v", err)), http.StatusInternalServerError)
		}
		return
	}

	var initialInput []map[string]interface{}
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&initialInput); err != nil {
			http.Error(w, jsonError("Invalid input payload"), http.StatusBadRequest)
			return
		}
	}

	// Parse the YAML workflow definition
	workflowDef, err := framework.LoadWorkflowDefFromYAMLString(storedWorkflow.Definition)
	if err != nil {
		http.Error(w, jsonError(fmt.Sprintf("Failed to parse workflow definition: %v", err)), http.StatusInternalServerError)
		return
	}

	log.Printf("DEBUG: Stored Workflow Definition:\n%s\n", storedWorkflow.Definition)
	log.Printf("DEBUG: Parsed WorkflowDef Nodes: %+v\n", workflowDef.Nodes)

	// Create a framework.Workflow instance
	wf := &framework.Workflow{
		Definition:  workflowDef,
		ErrorConnections: map[string]string{}, // TODO: Get from workflowDef
	}

	// Create a framework.Context (simplified for now)
	logger, err := framework.NewLogger()
	if err != nil {
		log.Printf("Failed to create logger: %v", err)
		http.Error(w, jsonError("Internal server error"), http.StatusInternalServerError)
		return
	}
	ctx := &framework.Context{
		Ctx:    r.Context(),
		Logger: logger,
		// Other context fields (HTTPClient, LangChain, DynamoDBClient, Metrics) would be initialized here
	}

	// Run the workflow in a goroutine to avoid blocking the API response
	go func() {
		// For now, hardcode the start node. This should ideally be part of the workflowDef.
		startNode := "manualTrigger"
		if err := wf.Run(ctx, startNode, initialInput); err != nil {
			log.Printf("Workflow %s run failed: %v", storedWorkflow.ID, err)
		}
		log.Printf("Workflow %s run completed.", storedWorkflow.ID)
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Workflow triggered successfully",
		"workflow_run_id": uuid.New().String(), // Placeholder for a real run ID
	})
}

func jsonError(message string) string {
	b, _ := json.Marshal(APIError{Message: message})
	return string(b)
}

