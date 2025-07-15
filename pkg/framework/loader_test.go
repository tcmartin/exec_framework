package framework

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadFromYAML(t *testing.T) {
	yamlContent := `
nodes:
  - node1
  - node2
connections:
  node1:
    - node2
`
	tmpfile, err := ioutil.TempFile("", "test.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(yamlContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	def, err := LoadFromYAML(tmpfile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(def.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(def.Nodes))
	}
	if len(def.Connections["node1"]) != 1 {
		t.Errorf("expected 1 connection from node1, got %d", len(def.Connections["node1"]))
	}
	if def.Connections["node1"][0] != "node2" {
		t.Errorf("expected connection from node1 to node2, got %s", def.Connections["node1"][0])
	}
}

func TestConvertN8nJSON(t *testing.T) {
	jsonContent := `{
		"connections": {
			"9": {
				"main": [
					[
						{
							"node": "10",
							"type": "main"
						}
					]
				]
			}
		}
	}`
	tmpfile, err := ioutil.TempFile("", "test.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(jsonContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	def, err := ConvertN8nJSON(tmpfile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(def.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(def.Nodes))
	}
	if len(def.Connections["9"]) != 1 {
		t.Errorf("expected 1 connection from node 9, got %d", len(def.Connections["9"]))
	}
	if def.Connections["9"][0] != "10" {
		t.Errorf("expected connection from node 9 to node 10, got %s", def.Connections["9"][0])
	}
}
