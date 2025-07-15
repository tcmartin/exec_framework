package test

import (
    "context"
    "encoding/json"
    "testing"

    "go-workflow/pkg/framework"
    "go-workflow/pkg/nodes"
    "github.com/tmc/langchaingo/llms"
)

// fakeLangChainClient implements the minimal interface
// required by OpenAINode
type fakeLangChainClient struct{}

func (f *fakeLangChainClient) GenerateFromSinglePrompt(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
    // echo user payload
    var payload map[string]interface{}
    _ = json.Unmarshal([]byte(prompt[len("test-prompt\n\nPayload: "):]), &payload)
    contentBytes, _ := json.Marshal(payload)
    return string(contentBytes), nil
}

func TestOpenAINode(t *testing.T) {
    fake := &nodes.OpenAINode{SystemPrompt: "test-prompt"}
    ctx := &framework.Context{
        Ctx:       context.Background(),
        LangChain: &fakeLangChainClient{},
    }
    input := []map[string]interface{}{{"foo":"bar"}}
    out, err := fake.Execute(ctx, input)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if val, ok := out[0]["foo"]; !ok || val != "bar" {
        t.Errorf("expected foo=bar, got %v", out[0])
    }
}