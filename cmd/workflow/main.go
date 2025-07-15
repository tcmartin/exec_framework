package main

import (
    "flag"
    "os"
    "context"
    retryablehttp "github.com/hashicorp/go-retryablehttp"
    "github.com/tmc/langchaingo"
    "go-workflow/pkg/framework"
    "go-workflow/pkg/nodes"
    "go-workflow/pkg/store"
    "go.uber.org/zap"
)

func main() {
    cfgPath := flag.String("config", "config/example_workflow.yaml", "workflow definition YAML")
    flag.Parse()

    logger, _ := framework.NewLogger()
    metrics := framework.NewMetrics()
    env := map[string]string{
        "X-API-KEY": os.Getenv("UNIPILE_API_KEY"),
        "account_id": os.Getenv("UNIPILE_ACCOUNT_ID"),
    }
    httpCli := retryablehttp.NewClient()
    langCli := langchaingo.NewClient(os.Getenv("OPENAI_API_KEY"))
    dbCli, _ := store.NewClient(context.Background())

    ctx := &framework.Context{Ctx: context.Background(), HTTPClient: httpCli,
        LangChain: langCli, DynamoDBClient: dbCli, Logger: logger, Metrics: metrics, Env: env}

    def, err := framework.LoadFromYAML(*cfgPath)
    if err != nil { panic(err) }

    // initialize nodesMap and build workflow

    wf := &framework.Workflow{Nodes: nodesMap, Connections: def.Connections}
    if err := wf.Run(ctx, "Trigger"); err != nil {
        logger.Fatalf("workflow error: %v", err)
    }
}
