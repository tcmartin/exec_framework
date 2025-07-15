package framework

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v3"
)

// WorkflowDef captures nodes and connections
type WorkflowDef struct {
    Nodes       []string            `yaml:"nodes"`
    Connections map[string][]string `yaml:"connections"`
}

// LoadFromYAML reads a YAML file into WorkflowDef
func LoadFromYAML(path string) (*WorkflowDef, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var def WorkflowDef
    if err := yaml.Unmarshal(data, &def); err != nil {
        return nil, err
    }
    return &def, nil
}

// ConvertN8nJSON parses an n8n JSON flow and emits WorkflowDef stub
func ConvertN8nJSON(path string) (*WorkflowDef, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var n8n struct{ Connections map[string][][]map[string]interface{} }
    if err := json.Unmarshal(data, &n8n); err != nil {
        return nil, err
    }
    def := &WorkflowDef{Connections: map[string][]string{}}
    for from, lists := range n8n.Connections {
        for _, arr := range lists {
            for _, step := range arr {
                to := fmt.Sprintf("%v", step["node"])
                def.Connections[from] = append(def.Connections[from], to)
            }
        }
    }
    // Nodes list inferred from keys
    for k := range def.Connections {
        def.Nodes = append(def.Nodes, k)
    }
    return def, nil
}
