# Go Workflow Framework Documentation

## Overview

The **Go Workflow Framework** is a modular, composable engine for defining and executing workflows in Go. Inspired by n8n, it allows you to:

* Define workflows via YAML (or convert n8n JSON to YAML).
* Use reusable **Node** implementations for common tasks:

  * HTTP requests with templating and retries
  * OpenAI (LangChain) chat completions
  * DynamoDB upserts
  * Custom code transformations
  * Timing controls (wait nodes)
* Orchestrate complex branching and loops with a simple engine.
* Load credentials and configuration from environment variables.
* Instrument with structured logging (Zap) and Prometheus metrics.

## Architecture

````
+-------------------------------------------------------------+
|                        cmd/workflow                         |
|  - main.go: loads YAML, constructs nodes, runs engine        |
|  - convert.go: CLI to translate n8n JSON to YAML definition |
+-------------------------------------------------------------+
        ▲                               ▲
        |                               |
        |        config/example_workflow.yaml        
        |                               |
        ▼                               ▼
+-------------------------------------------------------------+
|                         pkg/framework                       |
|  - Context: shared clients & config                         |
|  - Node: interface for steps                                |
|  - Workflow: engine for execution & branching               |
|  - Loader: YAML & n8n JSON parser                           |
|  - Logging: Zap setup                                       |
|  - Metrics: Prometheus collectors                           |
+-------------------------------------------------------------+
        ▲                               ▲
        |                               |
        |  registers/defaults            |
        ▼                               ▼
+-------------------------------------------------------------+
|                           pkg/nodes                         |
|  - ManualTrigger: seed payload                              |
|  - HTTPRequest: templated HTTP with retryablehttp           |
|  - CodeNode: arbitrary Go transform                          |
|  - WaitNode: randomized sleep                               |
|  - OpenAINode: LangChain chat completion wrapper            |
|  - DynamoDBUpsert: writes to DynamoDB table                 |
+-------------------------------------------------------------+
        ▲
        |
        ▼
+-------------------------------------------------------------+
|                           pkg/store                         |
|  - DynamoDBClient: AWS SDK v2 client loader                 |
+-------------------------------------------------------------+


## Getting Started

1. **Clone & Install**
   ```bash
   git clone <repo>
   cd go-workflow
   go mod tidy
````

2. **Set Environment Variables**

   ```bash
   export UNIPILE_API_KEY=...
   export UNIPILE_ACCOUNT_ID=...
   export OPENAI_API_KEY=...
   export OPENAI_MODEL=gpt-4.1
   export DYNAMODB_CONTACTS_TABLE=Contacts
   ```

   Ensure AWS credentials/region are configured (e.g. via `AWS_PROFILE`).

3. **Define Workflow**

   * Edit `config/example_workflow.yaml`, or
   * Convert from n8n JSON:

     ```bash
     go run cmd/workflow/convert.go -n8n path/to/flow.json > config/myflow.yaml
     ```

4. **Run Workflow**

   ```bash
   go run cmd/workflow/main.go -config config/example_workflow.yaml
   ```

5. **Monitor & Test**

   * **Tests**:

     ```bash
     go test ./pkg/... ./test/... -cover
     ```
   * **Metrics**: integrate `promhttp.Handler()` on `/metrics` to scrape metrics
   * **Logs**: structured JSON logs via Zap on stdout

## Extending the Framework

* **Add a new Node**:

  1. Create `pkg/nodes/your_node.go` implementing `framework.Node`.
  2. Register it in `cmd/workflow/main.go`.
  3. Update your YAML with the new node name and connections.

* **Advanced Retry**: configure `retryablehttp.Client` (e.g. `RetryMax`, `Backoff`).

* **Custom Metrics**: register additional Prometheus collectors in `framework/metrics.go`.

* **Configuration**: swap YAML loader for JSON or database-backed definitions.

## Exposing as HTTP API

To trigger workflows via external services (e.g. Python scripts, webhooks) you can wrap the engine in a simple HTTP server:

1. **Create a secure endpoint** (e.g. `/run`) expecting a shared secret header:
   ```go
   http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
       if r.Header.Get("X-Secret") != os.Getenv("WORKFLOW_SECRET") {
           http.Error(w, "Unauthorized", http.StatusUnauthorized)
           return
       }
       go func() {
           if err := wf.Run(ctx, "Trigger"); err != nil {
               logger.Errorf("workflow failed: %v", err)
           }
       }()
       w.WriteHeader(http.StatusAccepted)
       w.Write([]byte("Workflow started"))
   })
   log.Fatal(http.ListenAndServe(":8080", nil))

## Summary

This framework brings the power of n8n-style workflows to Go, combining:

* **Strong typing** & performance of Go
* **Flexibility** of composition and templating
* **Observability** through logging & metrics
* **Extensibility** via pluggable nodes

Use it to automate integrations, data pipelines, or any repeatable process with reliability and clarity.

