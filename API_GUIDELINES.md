Go Workflow Framework API & Configuration Guidelines

1. Supporting Multiple Workflows

Store each workflow definition as its own YAML file in config/ (e.g. flowA.yaml, flowB.yaml).

On startup, scan config/ and load all .yaml files via framework.LoadFromYAML:

files, _ := ioutil.ReadDir("config")
flows := map[string]*framework.Workflow{}
for _, f := range files {
    if strings.HasSuffix(f.Name(), ".yaml") {
        name := strings.TrimSuffix(f.Name(), ".yaml")
        def, _ := framework.LoadFromYAML("config/" + f.Name())
        flows[name] = &framework.Workflow{Nodes: nodesMap, Connections: def.Connections}
    }
}

Expose an HTTP endpoint pattern /run/{flowName} to trigger a named workflow.

2. HTTP API Design

Route: POST /run/{flowName}

Authentication: require X-Secret header matching WORKFLOW_SECRET env var.

Payload: JSON mapping to your ManualTrigger input (e.g. { "keywords": "..." }).

Async Response: return 202 Accepted immediately, spawning go wf.Run(...).

Healthcheck: GET /healthz returns 200 OK.

Metrics: expose GET /metrics via Prometheus handler.

mux.HandleFunc("/run/", runHandler)
mux.HandleFunc("/healthz", healthHandler)
mux.Handle("/metrics", promhttp.Handler())

3. Concurrency & Rate-Limiting

Async Execution: each /run spawns its own goroutine. Go scheduler handles thousands concurrently.

External API Limits: use a semaphore channel or worker pool in your HTTP and OpenAI nodes:

var sem = make(chan struct{}, 10) // limit to 10 concurrent calls

sem <- struct{}{}
resp, err := client.Do(req)
<-sem

4. Looping & Branching

Cycles: define cycles in connections: (e.g. Wait: [Details]) to re-enter nodes.

Termination: ensure upstream nodes eventually emit empty slices to break loops.

5. Deduplication Strategies

In-Memory Node: track seen keys in a map inside a DedupeNode. Good for single-instance.

External Store: use Redis or DynamoDB ConditionExpression to skip existing items.

DynamoDB Conditional Write:

input := &dynamodb.PutItemInput{
    TableName: &table,
    Item: item,
    ConditionExpression: aws.String("attribute_not_exists(pk)"),
}

6. Scheduling & Triggers

Cron Node: use robfig/cron or time.Ticker to invoke workflows on schedules.

Webhook Trigger: combine your HTTP /run endpoint with JSON payload input.

7. Error Handling & Retries

Global Retry: use retryablehttp.Client defaults (3 retries, exponential backoff).

Per-Node Error Branch: catch errors in engine.go and route items to an error node:

error: [ErrorHandler]

Logging: Zap logs include node name and error.

8. Observability

Logging: structured JSON logs via Zap.

Metrics: Prometheus histograms (node_duration_seconds) and counters (node_errors_total).

Tracing: integrate OpenTelemetry if needed by wrapping node execution with spans.

9. Extensibility

Add New Nodes: implement framework.Node, register in main.go and reference in YAML.

Dynamic Definitions: swap LoadFromYAML with a store-backed loader (S3, DB) for dynamic workflows.

These guidelines ensure your Go service can securely, concurrently, and flexibly serve any number of workflows, giving you both developer agility and production-grade robustness.
Go Workflow Framework API & Configuration Guidelines

1. Supporting Multiple Workflows

Store each workflow definition as its own YAML file in config/ (e.g. flowA.yaml, flowB.yaml).

On startup, scan config/ and load all .yaml files via framework.LoadFromYAML:

files, _ := ioutil.ReadDir("config")
flows := map[string]*framework.Workflow{}
for _, f := range files {
    if strings.HasSuffix(f.Name(), ".yaml") {
        name := strings.TrimSuffix(f.Name(), ".yaml")
        def, _ := framework.LoadFromYAML("config/" + f.Name())
        flows[name] = &framework.Workflow{Nodes: nodesMap, Connections: def.Connections}
    }
}

Expose an HTTP endpoint pattern /run/{flowName} to trigger a named workflow.

10. OAuth Token Propagation

Token Passing: While your service authenticates workflow runs via X-Secret, downstream OAuth‑protected APIs require a valid access token. Require clients to include their OAuth token in the run payload, e.g.:

{ "keywords": "...", "access_token": "<user_oauth_token>" }

Header Templating: In your HTTPRequest nodes, template the Authorization header from the record:

req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rec["access_token"]))

Or in YAML:

headers:
  Authorization: Bearer {{ .access_token }}

Security: Always use HTTPS to protect tokens in transit. Treat access tokens as sensitive secrets—avoid logging them and store refresh tokens securely if needed.

Token Refresh: For expired tokens, implement a TokenRefreshNode that exchanges a refresh token for a new access token before downstream calls.

Least Privilege: Request only the minimal scopes needed for each workflow.

These guidelines ensure your Go service can securely, concurrently, and flexibly serve any number of workflows, giving you both developer agility and production-grade robustness.
