package nodes

import (
    "encoding/json"
    "fmt"

    "go-workflow/pkg/framework"
)

func toJSON(m map[string]interface{}) string {
    b, _ := json.Marshal(m)
    return string(b)
}

// OpenAINode wraps a LangChain chat completion
type OpenAINode struct{ SystemPrompt string }

func (n *OpenAINode) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
    var out []map[string]interface{}
    for _, rec := range inputs {
        prompt := fmt.Sprintf("%s\n\nPayload: %s", n.SystemPrompt, toJSON(rec))
        resp, err := ctx.LangChain.GenerateFromSinglePrompt(ctx.Ctx, prompt)
        if err != nil {
            return nil, err
        }
        var parsed map[string]interface{}
        if err := json.Unmarshal([]byte(resp), &parsed); err != nil {
            return nil, err
        }
        out = append(out, parsed)
    }
    return out, nil
}