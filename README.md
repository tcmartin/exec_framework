# Go Workflow Framework

This is a modular, composable engine for defining and executing workflows in Go. It allows you to define workflows via YAML, use reusable Node implementations for common tasks, and orchestrate complex branching and loops.

## Getting Started

1.  **Clone & Install**

    ```bash
    git clone <repo>
    cd go-workflow
    go mod tidy
    ```

2.  **Set Environment Variables**

    ```bash
    export UNIPILE_API_KEY=...
    export UNIPILE_ACCOUNT_ID=...
    export OPENAI_API_KEY=...
    export OPENAI_MODEL=gpt-4.1
    export DYNAMODB_CONTACTS_TABLE=Contacts
    ```

    Ensure AWS credentials/region are configured (e.g. via `AWS_PROFILE`).

3.  **Define Workflow**

    *   Edit `config/example_workflow.yaml`, or
    *   Convert from n8n JSON:

        ```bash
        go run cmd/workflow/convert.go -n8n path/to/flow.json > config/myflow.yaml
        ```

4.  **Run Workflow**

    ```bash
    go run cmd/workflow/main.go -config config/example_workflow.yaml
    ```

5.  **Monitor & Test**

    *   **Tests**:

        ```bash
        go test ./pkg/... ./test/... -cover
        ```
    *   **Metrics**: integrate `promhttp.Handler()` on `/metrics` to scrape metrics
    *   **Logs**: structured JSON logs via Zap on stdout
