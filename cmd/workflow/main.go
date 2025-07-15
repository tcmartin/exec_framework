package main

import (
    "context"
    "flag"
    "os"

    retryablehttp "github.com/hashicorp/go-retryablehttp"
    "github.com/tmc/langchaingo"
    "go.uber.org/zap"

    "go-workflow/pkg/framework"
    "go-workflow/pkg/nodes"
    "go-workflow/pkg/store"
)

func main() {
    cfgPath := flag.String("config", "config/example_workflow.yaml", "Path to workflow YAML definition")
    flag.Parse()

    logger, err := framework.NewLogger()
    if err != nil {
        panic(err)
    }
    metrics := framework.NewMetrics()

    env := map[string]string{
        "X-API-KEY":  os.Getenv("UNIPILE_API_KEY"),
        "account_id": os.Getenv("UNIPILE_ACCOUNT_ID"),
    }

    httpCli := retryablehttp.NewClient()
    langCli := langchaingo.NewClient(os.Getenv("OPENAI_API_KEY"))
    dbCli, err := store.NewClient(context.Background())
    if err != nil {
        logger.Fatalf("failed to create DynamoDB client: %v", err)
    }

    ctx := &framework.Context{
        Ctx:            context.Background(),
        HTTPClient:     httpCli,
        LangChain:      langCli,
        DynamoDBClient: dbCli,
        Logger:         logger,
        Metrics:        metrics,
        Env:            env,
    }

    def, err := framework.LoadFromYAML(*cfgPath)
    if err != nil {
        logger.Fatalf("failed to load workflow definition: %v", err)
    }

    nodesMap := map[string]framework.Node{
        "Trigger": &nodes.ManualTrigger{Payload: []map[string]interface{}{{"keywords": "business development manager fintech"}}},
        "Search": nodes.NewHTTPRequest(
            "POST",
            "https://api15.unipile.com:14501/api/v1/linkedin/search?account_id={{.account_id}}",
            `{"api":"classic","category":"people","keywords":"{{.keywords}}"}`,
        ),
        "Slice": &nodes.CodeNode{Fn: func(items []map[string]interface{}) []map[string]interface{} {
            if len(items) > 20 {
                items = items[:20]
            }
            for i := range items {
                items[i]["outreach_index"] = i + 1
                items[i]["total_to_process"] = len(items)
            }
            return items
        }},
        "Details": nodes.NewHTTPRequest(
            "GET",
            "https://api15.unipile.com:14501/api/v1/users/{{.public_identifier}}?linkedin_sections=*&account_id={{.account_id}}",
            "",
        ),
        "AI":    &nodes.OpenAINode{SystemPrompt: os.Getenv("OPENAI_SYSTEM_PROMPT")},
        "Store": &nodes.DynamoDBUpsert{TableName: os.Getenv("DYNAMODB_CONTACTS_TABLE")},
        "Invite": nodes.NewHTTPRequest(
            "POST",
            "https://api15.unipile.com:14501/api/v1/users/invite",
            `{"provider_id":"{{.provider_id}}","account_id":"{{.account_id}}"}`,
        ),
        "Wait": &nodes.WaitNode{MaxSeconds: 30},
    }

    wf := &framework.Workflow{Nodes: nodesMap, Connections: def.Connections}
    if err := wf.Run(ctx, "Trigger"); err != nil {
        logger.Fatalf("workflow execution failed: %v", err)
    }
}