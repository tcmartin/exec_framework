package nodes

import (
	"go-workflow/pkg/framework"
)

// WebhookTrigger represents an incoming webhook payload.
type WebhookTrigger struct{}

// NewWebhookTrigger creates a new WebhookTrigger node.
func NewWebhookTrigger() *WebhookTrigger {
	return &WebhookTrigger{}
}

// Execute simply returns the input records, as the webhook payload is the input.
func (n *WebhookTrigger) Execute(ctx *framework.Context, inputs []map[string]interface{}) ([]map[string]interface{}, error) {
	return inputs, nil
}
