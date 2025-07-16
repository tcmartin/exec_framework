# Workflow Guide

This guide explains how to define, create, and run workflows within this project.

## What is a Workflow?

A workflow in this system is a sequence of interconnected nodes, where each node performs a specific task. Data flows between these nodes, allowing for complex automated processes. Workflows are defined using YAML files, specifying the nodes involved and their connections.

## Workflow Definition (YAML)

Workflows are defined in YAML files, typically located in the `config/` directory. A workflow definition consists of two main parts:

1.  **`nodes`**: A list of node names used in the workflow.
2.  **`connections`**: A mapping that defines the flow of data between nodes. Each key is a source node, and its value is a list of destination nodes.

### Example Workflow YAML (`config/example_workflow.yaml` - conceptual)

```yaml
nodes:
  - Trigger
  - Search
  - Slice
  - Details
  - AI
  - Store
  - Invite
  - Wait

connections:
  Trigger:
    - Search
  Search:
    - Slice
  Slice:
    - Details
  Details:
    - AI
  AI:
    - Store
  Store:
    - Invite
  Invite:
    - Wait
  Wait:
    - Search # Example of a loop or branching back
```

## Available Nodes

The system provides several built-in node types, each with a specific function. These nodes are instantiated and configured in the `cmd/workflow/main.go` file.

*   **`ManualTrigger`**: Initiates the workflow with a predefined payload.
    *   Example: `&nodes.ManualTrigger{Payload: []map[string]interface{}{{"keywords": "business development manager fintech"}}}`
*   **`HTTPRequest`**: Performs HTTP requests (GET, POST, etc.) to external APIs. It supports templating for URL and body using data from previous nodes and environment variables.
    *   Example: `nodes.NewHTTPRequest("POST", "https://api15.unipile.com:14501/api/v1/linkedin/search?account_id={{.account_id}}", `{"api":"classic","category":"people","keywords":"{{.keywords}}"}}`
*   **`CodeNode`**: Executes a custom Go function. This is useful for data transformation or custom logic.
    *   Example: `&nodes.CodeNode{Fn: func(items []map[string]interface{}) []map[string]interface{} { ... }}`
*   **`OpenAINode`**: Interacts with the OpenAI API, typically for AI-driven tasks.
    *   Example: `&nodes.OpenAINode{SystemPrompt: os.Getenv("OPENAI_SYSTEM_PROMPT")}`
*   **`DynamoDBUpsert`**: Upserts data into a DynamoDB table.
    *   Example: `&nodes.DynamoDBUpsert{TableName: os.Getenv("DYNAMODB_CONTACTS_TABLE")}`
*   **`WaitNode`**: Pauses the workflow for a specified duration.
    *   Example: `&nodes.WaitNode{MaxSeconds: 30}`

## Creating a Workflow

To create a new workflow:

1.  **Define the Workflow Structure (YAML):** Create a new YAML file (e.g., `config/my_new_workflow.yaml`) and define your nodes and their connections as shown in the "Workflow Definition" section above.
2.  **Instantiate Nodes (Go):** Open `cmd/workflow/main.go`. In the `main` function, locate the `nodesMap` variable. Add new entries to this map for any custom nodes you need, or configure the existing node types with your specific parameters. Ensure the keys in `nodesMap` match the node names you used in your YAML file.
3.  **Configure Environment Variables:** Many nodes (e.g., `HTTPRequest`, `OpenAINode`, `DynamoDBUpsert`) rely on environment variables for API keys, table names, etc. Ensure all necessary environment variables are set before running your workflow. Refer to `cmd/workflow/main.go` for the specific environment variables used by each node.

## Running a Workflow

To run a workflow:

1.  **Build the Workflow Executable:**
    ```bash
    go build -o workflow cmd/workflow/main.go
    ```
2.  **Set Environment Variables:** Export any required environment variables (e.g., `UNIPILE_API_KEY`, `OPENAI_API_KEY`, `DYNAMODB_CONTACTS_TABLE`).
    ```bash
    export UNIPILE_API_KEY="your_unipile_api_key"
    export UNILE_ACCOUNT_ID="your_unipile_account_id"
    export OPENAI_API_KEY="your_openai_api_key"
    export DYNAMODB_CONTACTS_TABLE="your_dynamodb_table_name"
    export AWS_REGION="your_aws_region"
    export AWS_ACCESS_KEY_ID="your_aws_access_key_id"
    export AWS_SECRET_ACCESS_KEY="your_aws_secret_access_key"
    ```
3.  **Execute the Workflow:** Run the compiled executable, specifying your workflow YAML file using the `-config` flag.
    ```bash
    ./workflow -config config/example_workflow.yaml
    ```
    Replace `config/example_workflow.yaml` with the path to your workflow definition file.

## Converting n8n JSON to Workflow YAML

The `cmd/workflow/convert.go` utility can convert an n8n JSON flow definition into a basic workflow YAML structure. This can be a starting point for migrating n8n flows to this system.

To use it:

1.  **Build the Converter Executable:**
    ```bash
    go build -o convert cmd/workflow/convert.go
    ```
2.  **Convert the n8n JSON:**
    ```bash
    ./convert -n8n /path/to/your/n8n_flow.json > config/my_converted_workflow.yaml
    ```
    This will output the generated YAML to standard output, which you can redirect to a new YAML file. You will likely need to manually adjust the generated YAML and configure the nodes in `cmd/workflow/main.go` to match your specific requirements.
