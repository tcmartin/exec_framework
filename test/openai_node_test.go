package test

import (
    "context"
    "encoding/json"
    "testing"
    "go-workflow/pkg/nodes"
    "go-workflow/pkg/framework"
)

// fakeLangChainClient implements the minimal interface needed by OpenAINode
// specifically, the ChatCompletion method
var _ framework.LangChainClient = (*fakeLangChainClient)(nil)

type fakeLangChainClient struct {}

func (f *fakeLangChainClient) ChatCompletion(ctx context.Context, req langchaingo.ChatCompletionRequest) (langchaingo.ChatCompletionResponse, error) {
    // Echo back the user's payload JSON in a deterministic field
    // Assume the last message is the user content
    var payload map[string]interface{}
    _ = json.Unmarshal([]byte(req.Messages[len(req.Messages)-1].Content[len("Payload: "):]), &payload)
    // Construct a fake response wrapping the payload
    contentBytes, _ := json.Marshal(payload)
    return langchaingo.ChatCompletionResponse{
        Choices: []langchaingo.ChatChoice{{Message: langchaingo.ChatMessage{Content: string(contentBytes)}}},
    }, nil
}

func TestOpenAINode(t *testing.T) {
    fake := &nodes.OpenAINode{SystemPrompt: "test-prompt"}
    ctx := &framework.Context{
        Ctx:        context.Background(),
        LangChain:  &fakeLangChainClient{},
        Logger:     nil,
        Metrics:    nil,
        HTTPClient: nil,
        Env:        nil,
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
