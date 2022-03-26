package updater

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type MachConfig struct {
	Components yaml.Node
}

func MachConfigUpdater(src []byte, updates *UpdateSet) []byte {
	data := MachConfig{}
	err := yaml.Unmarshal(src, &data)
	if err != nil {
		log.Fatalln(err)
	}

	// Create a mapping where the key is the component name and the value the
	// yaml node of the component. We do this by iterating all component children
	// and searching for the `name` tag.
	nodes := map[string]*yaml.Node{}
	for _, cn := range data.Components.Content {
		for i, n := range cn.Content {
			if n.Tag == "!!str" && n.Value == "name" {
				name := cn.Content[i+1].Value
				nodes[name] = cn
				break
			}
		}
	}

	// Walk through the updated components and search the corresponding yaml node
	// via the previously created mapping. Withing the node search for the
	// `version` tag and use the line number to change the value in the source
	// document (lines list)
	lines := SplitLines(string(src))
	for _, c := range updates.components {
		node, ok := nodes[c.component.Name]
		if !ok {
			logrus.Warn("Component with update not found in yaml file")
			continue
		}

		for i, n := range node.Content {
			if n.Tag == "!!str" && n.Value == "version" {

				// The value is in the node after this node. Assume it's always
				// sequential
				vn := node.Content[i+1]
				if vn.Value != c.component.Version {
					log.Fatal("Unexpected version")
				}

				// Make sure the version is always quoted
				replacement := c.version
				if lines[vn.Line-1][vn.Column-1] != '"' {
					replacement = fmt.Sprintf(`"%s"`, replacement)
				}

				key := lines[vn.Line-1][:vn.Column-1]
				value := lines[vn.Line-1][vn.Column-1:]

				lines[vn.Line-1] = key + strings.Replace(value, vn.Value, replacement, 1)
				break
			}
		}
	}

	output := strings.Join(lines, "\n") + "\n"
	return []byte(output)
}

// MachFileWriter updates the contents of a mach file with the updated
// version of the components
func MachFileWriter(updates *UpdateSet) {

	input, err := ioutil.ReadFile(updates.filename)
	if err != nil {
		log.Fatalln(err)
	}

	output := MachConfigUpdater(input, updates)

	err = ioutil.WriteFile(updates.filename, output, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
