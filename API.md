# Workflow API Specification

This document describes the RESTful API for managing and triggering workflows.

## Base URL

`/api/v1`

## Endpoints

### 1. Upload Workflow Definition

`POST /workflows`

Uploads a new workflow definition to the system. The workflow definition should be provided in YAML format.

*   **Request Body:**
    ```json
    {
        "name": "my_workflow",
        "definition": "<YAML_WORKFLOW_DEFINITION_AS_STRING>"
    }
    ```
*   **Responses:**
    *   `201 Created`: Workflow uploaded successfully.
        ```json
        {
            "id": "<workflow_id>",
            "name": "my_workflow"
        }
        ```
    *   `400 Bad Request`: Invalid request body or workflow definition.
    *   `500 Internal Server Error`: Server error.

### 2. Get Workflow Definition

`GET /workflows/{id}`

Retrieves a stored workflow definition by its ID.

*   **Path Parameters:**
    *   `id` (string, required): The ID of the workflow.
*   **Responses:**
    *   `200 OK`: Workflow definition retrieved successfully.
        ```json
        {
            "id": "<workflow_id>",
            "name": "my_workflow",
            "definition": "<YAML_WORKFLOW_DEFINITION_AS_STRING>"
        }
        ```
    *   `404 Not Found`: Workflow with the specified ID not found.
    *   `500 Internal Server Error`: Server error.

### 3. Trigger Workflow

`POST /workflows/{id}/run`

Triggers a stored workflow by its ID with an optional input payload.

*   **Path Parameters:**
    *   `id` (string, required): The ID of the workflow to trigger.
*   **Request Body (Optional):**
    ```json
    [
        {
            "key1": "value1",
            "key2": "value2"
        }
    ]
    ```
    An array of maps representing the initial input records for the workflow.
*   **Responses:**
    *   `202 Accepted`: Workflow trigger request accepted. The workflow will be executed asynchronously.
        ```json
        {
            "message": "Workflow triggered successfully",
            "workflow_run_id": "<unique_run_id>"
        }
        ```
    *   `404 Not Found`: Workflow with the specified ID not found.
    *   `400 Bad Request`: Invalid request body.
    *   `500 Internal Server Error`: Server error.

### 4. Get Workflow Run Status (Placeholder)

`GET /workflows/{id}/status`

Retrieves the status of a specific workflow run. This endpoint is a placeholder and will require a more robust state management system for asynchronous workflow executions.

*   **Path Parameters:**
    *   `id` (string, required): The ID of the workflow run.
*   **Responses:**
    *   `200 OK`: Workflow run status retrieved successfully.
        ```json
        {
            "workflow_run_id": "<unique_run_id>",
            "status": "running" | "completed" | "failed",
            "start_time": "<timestamp>",
            "end_time": "<timestamp>" (if completed/failed),
            "error_message": "<error_details>" (if failed)
        }
        ```
    *   `404 Not Found`: Workflow run with the specified ID not found.
    *   `500 Internal Server Error`: Server error.
