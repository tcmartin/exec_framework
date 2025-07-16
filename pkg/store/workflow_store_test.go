package store

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestSQLiteStore(t *testing.T) {
	dbPath := "test_workflows.db"
	defer os.Remove(dbPath) // Clean up after test

	store := NewSQLiteStore(dbPath)

	// Test Init
	if err := store.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Test SaveWorkflow
	workflow1 := &Workflow{
		ID:         uuid.New().String(),
		Name:       "workflow1",
		Definition: "workflow_def_1",
	}
	if err := store.SaveWorkflow(workflow1); err != nil {
		t.Fatalf("SaveWorkflow failed: %v", err)
	}

	workflow2 := &Workflow{
		ID:         uuid.New().String(),
		Name:       "workflow2",
		Definition: "workflow_def_2",
	}
	if err := store.SaveWorkflow(workflow2); err != nil {
		t.Fatalf("SaveWorkflow failed: %v", err)
	}

	// Test GetWorkflow
	retrievedWorkflow, err := store.GetWorkflow(workflow1.ID)
	if err != nil {
		t.Fatalf("GetWorkflow failed: %v", err)
	}
	if retrievedWorkflow.Name != workflow1.Name || retrievedWorkflow.Definition != workflow1.Definition {
		t.Errorf("GetWorkflow mismatch: expected %v, got %v", workflow1, retrievedWorkflow)
	}

	// Test GetWorkflowByName
	retrievedWorkflowByName, err := store.GetWorkflowByName(workflow2.Name)
	if err != nil {
		t.Fatalf("GetWorkflowByName failed: %v", err)
	}
	if retrievedWorkflowByName.ID != workflow2.ID || retrievedWorkflowByName.Definition != workflow2.Definition {
		t.Errorf("GetWorkflowByName mismatch: expected %v, got %v", workflow2, retrievedWorkflowByName)
	}

	// Test ListWorkflows
	workflows, err := store.ListWorkflows()
	if err != nil {
		t.Fatalf("ListWorkflows failed: %v", err)
	}
	if len(workflows) != 2 {
		t.Errorf("ListWorkflows mismatch: expected 2, got %d", len(workflows))
	}

	// Test SaveWorkflow with duplicate name (should fail due to UNIQUE constraint)
	duplicateWorkflow := &Workflow{
		ID:         uuid.New().String(),
		Name:       "workflow1", // Duplicate name
		Definition: "workflow_def_duplicate",
	}
	if err := store.SaveWorkflow(duplicateWorkflow); err == nil {
		t.Error("SaveWorkflow with duplicate name did not return an error")
	}
}
