package frontmatter

import (
	"gopkg.in/yaml.v3"
)

func yamlFindLineAndCol(node *yaml.Node, parts []string) (line int, col int) {
	part := parts[0]

	for _, node := range node.Content {
		if node.Kind == yaml.MappingNode {
			for i := 0; i < len(node.Content); i += 2 {
				if node.Content[i].Value == part {
					if len(parts) == 1 {
						return node.Content[i].Line, node.Content[i].Column
					}
					return yamlFindLineAndCol(node, parts[1:])
				}
			}
		}
	}
	return -1, -1
}
