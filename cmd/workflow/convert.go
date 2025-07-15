package main

import (
    "flag"
    "fmt"
    "go-workflow/pkg/framework"
    "gopkg.in/yaml.v3"
)

func main() {
    n8nPath := flag.String("n8n", "", "Path to n8n JSON flow")
    out := flag.String("out", "workflow.yaml", "Output YAML file")
    flag.Parse()

    def, err := framework.ConvertN8nJSON(*n8nPath)
    if err != nil {
        panic(err)
    }
    data, err := yaml.Marshal(def)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s", data)
}
