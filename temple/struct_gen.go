package temple

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// varRegex is a regular expression that matches valid variable and field names
	// Go requires that variables and fields must be a letter followed by any number
	// of alphanumeric characters.
	varRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`)
)

type Node struct {
	Name            string
	Key             string
	Children        map[string]*Node
	OrderedChildren []*Node
}

func NewNode(name string, key string) *Node {
	return &Node{
		Name:     name,
		Key:      key,
		Children: map[string]*Node{},
	}
}

func createTreeFromPaths(paths []string) (*Node, error) {
	root := NewNode("", "")
	for _, path := range paths {
		currNode := root
		tokens := strings.Split(path, string(os.PathSeparator))
		for i, token := range tokens {
			fieldName := strings.Title(convertUnderscoresToCamelCase(token))
			if !varRegex.MatchString(fieldName) {
				return nil, fmt.Errorf("Cannot generate struct because %s cannot be converted to a valid field name.\n(Must only contain alphanumeric characters and underscores, and the first character must be a letter).", token)
			}
			var newNode *Node
			if i == len(tokens)-1 {
				// Last token, set the key property
				newNode = NewNode(fieldName, path)
			} else {
				newNode = NewNode(fieldName, "")
			}
			if _, found := currNode.Children[fieldName]; !found {
				currNode.Children[fieldName] = newNode
				currNode.OrderedChildren = append(currNode.OrderedChildren, newNode)
			} else if currNode.Children[fieldName].Key != "" || newNode.Key != "" {
				return nil, fmt.Errorf("Cannot generate struct because %s is both a file and directory", filepath.Join(tokens[:i+1]...))
			}
			currNode = currNode.Children[fieldName]
		}
	}
	return root, nil
}

func generateStructDeclr(structName string, root *Node) string {
	if root == nil {
		return ""
	}
	result := fmt.Sprintf("var %s struct {\n", structName)
	for _, child := range root.OrderedChildren {
		result += child.Name + " " + generateInnerDeclr(child)
	}
	result += "}\n"
	return result
}

func generateInnerDeclr(root *Node) string {
	if len(root.Children) == 0 {
		return "string\n"
	}
	result := "struct {\n"
	for _, child := range root.OrderedChildren {
		result += child.Name + " " + generateInnerDeclr(child)
	}
	result += "}\n"
	return result
}

func generateStructInit(baseName string, mapName string, nodes []*Node) string {
	result := ""
	for _, node := range nodes {
		newBase := baseName + "." + node.Name
		if node.Key != "" {
			result += newBase + fmt.Sprintf(" = %s[\"%s\"]\n", mapName, node.Key)
		}
		result += generateStructInit(newBase, mapName, node.OrderedChildren)
	}
	return result
}

// convertUnderscoresToCamelCase converts a string of the form
// foo_bar_baz to fooBarBaz.
func convertUnderscoresToCamelCase(s string) string {
	if len(s) == 0 {
		return ""
	}
	result := ""
	shouldUpper := false
	for _, char := range s {
		if char == '_' {
			shouldUpper = true
			continue
		}
		if shouldUpper {
			result += strings.ToUpper(string(char))
		} else {
			result += string(char)
		}
		shouldUpper = false
	}
	return result
}
