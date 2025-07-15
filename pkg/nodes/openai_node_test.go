package nodes

import (
    "context"
    "encoding/json"
    "testing"
    "errors"

    "go-workflow/pkg/framework"
    "github.com/tmc/langchaingo/llms"
)

// fakeLangChainClient implements the minimal interface
// required by OpenAINode
type fakeLangChainClient struct {
    generateFromSinglePrompt func(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
}

func (f *fakeLangChainClient) GenerateFromSinglePrompt(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
    if f.generateFromSinglePrompt != nil {
        return f.generateFromSinglePrompt(ctx, prompt, options...)
    }
    // echo user payload
    var payload map[string]interface{}
    _ = json.Unmarshal([]byte(prompt[len("test-prompt\n\nPayload: "):]), &payload)
    contentBytes, _ := json.Marshal(payload)
    return string(contentBytes), nil
}

func TestOpenAINode(t *testing.T) {
    fake := &OpenAINode{SystemPrompt: "test-prompt"}
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

func TestOpenAINode_LangChainError(t *testing.T) {
    fake := &OpenAINode{SystemPrompt: "test-prompt"}
    ctx := &framework.Context{
        Ctx: context.Background(),
        LangChain: &fakeLangChainClient{
            generateFromSinglePrompt: func(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
                return "", errors.New("LangChain error")
            },
        },
    }
    input := []map[string]interface{}{{"foo":"bar"}}
    _, err := fake.Execute(ctx, input)
    if err == nil {
        t.Fatal("expected an error, but got nil")
    }
    if err.Error() != "LangChain error" {
        t.Errorf("expected LangChain error, got %v", err)
    }
}

func TestOpenAINode_InvalidJSONResponse(t *testing.T) {
    fake := &OpenAINode{SystemPrompt: "test-prompt"}
    ctx := &framework.Context{
        Ctx: context.Background(),
        LangChain: &fakeLangChainClient{
            generateFromSinglePrompt: func(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
                return "invalid json", nil
            },
        },
    }
    input := []map[string]interface{}{{"foo":"bar"}}
    _, err := fake.Execute(ctx, input)
    if err == nil {
        t.Fatal("expected an error, but got nil")
    }
}

func TestToJSON(t *testing.T) {
    input := map[string]interface{}{"foo": "bar"}
    expected := `{"foo":"bar"}`
    result := toJSON(input)
    if result != expected {
        t.Errorf("expected %s, got %s", expected, result)
    }
}
