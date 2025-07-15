package nodes

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/tmc/langchaingo"
    "go-workflow/pkg/framework"
)

// OpenAINode calls LangChain chat completion
type OpenAINode struct{ SystemPrompt string }

// Execute sends system+user messages, then parses the JSON response
func (n *OpenAINode) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
    var out []map[string]interface{}
    for _, rec := range inputs {
        msgs := []langchaingo.ChatMessage{
            {Role: "system", Content: n.SystemPrompt},
            {Role: "user", Content: fmt.Sprintf("Payload: %s", toJSON(rec))},
        }
        resp, err := ctx.LangChain.ChatCompletion(ctx.Ctx, langchaingo.ChatCompletionRequest{
            Model:    os.Getenv("OPENAI_MODEL"),
            Messages: msgs,
        })
        if err != nil {
            return nil, err
        }
        var parsed map[string]interface{}
        if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &parsed); err != nil {
            return nil, err
        }
        out = append(out, parsed)
    }
    return out, nil
}

func toJSON(m map[string]interface{}) string {
    b, _ := json.Marshal(m)
    return string(b)
}

