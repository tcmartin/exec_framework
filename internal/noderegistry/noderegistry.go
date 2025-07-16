package noderegistry

import (
	"log"

	"go-workflow/pkg/framework"
	"go-workflow/pkg/nodes"

	"gopkg.in/yaml.v3"
)

func init() {
	framework.RegisterNodeFactory("manualTrigger", func(nodeDef *yaml.Node) (framework.Node, error) {
		var temp struct {
			Type string `yaml:"type"`
			Payload []map[string]interface{} `yaml:"payload"`
		}
		if err := nodeDef.Decode(&temp); err != nil {
			return nil, err
		}
		return nodes.NewManualTrigger(temp.Payload), nil
	})

	framework.RegisterNodeFactory("setNode", func(nodeDef *yaml.Node) (framework.Node, error) {
		var temp struct {
			Type string `yaml:"type"`
			SetValues map[string]interface{} `yaml:"setValues"`
			RemoveKeys []string `yaml:"removeKeys"`
		}
		if err := nodeDef.Decode(&temp); err != nil {
			return nil, err
		}
		return nodes.NewSetNode(temp.SetValues, temp.RemoveKeys), nil
	})

	framework.RegisterNodeFactory("dedupeNode", func(nodeDef *yaml.Node) (framework.Node, error) {
		var temp struct {
			Type string `yaml:"type"`
			Key string `yaml:"key"`
		}
		if err := nodeDef.Decode(&temp); err != nil {
			return nil, err
		}
		return nodes.NewDedupeNode(temp.Key), nil
	})

	framework.RegisterNodeFactory("mergeNode", func(nodeDef *yaml.Node) (framework.Node, error) {
		var temp struct {
			Type string `yaml:"type"`
			Key string `yaml:"key"`
		}
		if err := nodeDef.Decode(&temp); err != nil {
			return nil, err
		}
		return nodes.NewMergeNode(temp.Key), nil
	})

	framework.RegisterNodeFactory("splitInBatchesNode", func(nodeDef *yaml.Node) (framework.Node, error) {
		var temp struct {
			Type string `yaml:"type"`
			BatchSize int `yaml:"batchSize"`
		}
		if err := nodeDef.Decode(&temp); err != nil {
			return nil, err
		}
		return nodes.NewSplitInBatchesNode(temp.BatchSize), nil
	})

	framework.RegisterNodeFactory("switchNode", func(nodeDef *yaml.Node) (framework.Node, error) {
		var temp struct {
			Type string `yaml:"type"`
			Conditions map[string]nodes.Condition `yaml:"conditions"`
		}
		if err := nodeDef.Decode(&temp); err != nil {
			return nil, err
		}
		return nodes.NewSwitchNode(temp.Conditions), nil
	})

	framework.RegisterNodeFactory("waitForNode", func(nodeDef *yaml.Node) (framework.Node, error) {
		var temp struct {
			Type string `yaml:"type"`
			TimestampKey string `yaml:"timestampKey"`
		}
		if err := nodeDef.Decode(&temp); err != nil {
			return nil, err
		}
		return nodes.NewWaitForNode(temp.TimestampKey), nil
	})

	framework.RegisterNodeFactory("errorHandlerNode", func(nodeDef *yaml.Node) (framework.Node, error) {
		var temp struct {
			Type string `yaml:"type"`
		}
		if err := nodeDef.Decode(&temp); err != nil {
			return nil, err
		}
		return nodes.NewErrorHandlerNode(), nil
	})

	framework.RegisterNodeFactory("httpRequest", func(nodeDef *yaml.Node) (framework.Node, error) {
		var temp struct {
			Type string `yaml:"type"`
			URLKey string `yaml:"urlKey"`
			MethodKey string `yaml:"methodKey"`
			HeadersKey string `yaml:"headersKey"`
			BodyKey string `yaml:"bodyKey"`
		}
		if err := nodeDef.Decode(&temp); err != nil {
			return nil, err
		}
		return nodes.NewHTTPRequest(temp.URLKey, temp.MethodKey, temp.HeadersKey, temp.BodyKey), nil
	})

	framework.RegisterNodeFactory("openaiNode", func(nodeDef *yaml.Node) (framework.Node, error) {
		var temp struct {
			Type string `yaml:"type"`
			SystemPrompt string `yaml:"systemPrompt"`
		}
		if err := nodeDef.Decode(&temp); err != nil {
			return nil, err	
		}
		return nodes.NewOpenAINode(temp.SystemPrompt), nil
	})

	framework.RegisterNodeFactory("waitNode", func(nodeDef *yaml.Node) (framework.Node, error) {
		var temp struct {
			Type string `yaml:"type"`
			MaxSeconds int `yaml:"maxSeconds"`
		}
		if err := nodeDef.Decode(&temp); err != nil {
			return nil, err
		}
		return nodes.NewWaitNode(temp.MaxSeconds), nil
	})

	framework.RegisterNodeFactory("dynamodbUpsert", func(nodeDef *yaml.Node) (framework.Node, error) {
		var temp struct {
			Type string `yaml:"type"`
			TableNameKey string `yaml:"tableNameKey"`
			AWSRegionKey string `yaml:"awsRegionKey"`
			AWSAccessKeyIDKey string `yaml:"awsAccessKeyIDKey"`
			AWSSecretAccessKeyKey string `yaml:"awsSecretAccessKeyKey"`
		}
		if err := nodeDef.Decode(&temp); err != nil {
			return nil, err
		}
		return nodes.NewDynamoDBUpsert(temp.TableNameKey, temp.AWSRegionKey, temp.AWSAccessKeyIDKey, temp.AWSSecretAccessKeyKey), nil
	})

	framework.RegisterNodeFactory("codeNode", func(nodeDef *yaml.Node) (framework.Node, error) {
		var temp struct {
			Type string `yaml:"type"`
			// CodeNode's Fn field is a function, which cannot be directly decoded from YAML.
			// This factory will need to handle how the function is provided or referenced.
			// For now, we'll return a basic CodeNode without a functional Fn.
		}
		if err := nodeDef.Decode(&temp); err != nil {
			return nil, err
		}
		// This will create a CodeNode with a nil Fn, which will likely cause a panic if executed.
		// A proper solution would involve a way to define and load Go functions dynamically or via a registry.
		return nodes.NewCodeNode(nil), nil
	})
}