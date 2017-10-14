package kaanparser

import (
	"regexp/syntax"
)

// Examples for unit tests.
var examples = []string{
	"[a,b,c]",
	"[a[aa[aaa],ab,ac],b,c[ca,cb,cc[cca]]]",
	"[a,b, c ]]]]]]",
	" ",
	"◌́",
	"123]]]",
	"\n\n",
	"]]][a,b, c ]]]]]]",
	//"asd",			// Err: Seg, nil ptr
	//"",				// Err: Slice range
}

// Scan - the scanner structure.
var scan struct {
	name  string // Child name
	err   string // Error
}

type node struct {
	Name     string  `json:"name"`
	Parent   *node   `json:"-"` // Parent pointer - depth iteration
	Children []*node `json:"children,omitempty"`
}


func checkName() bool {
	if len(scan.name) > 0 {
		return true
	} else {
		return false
	}
}

// Parse next rune as string
func parseNext(v rune, node **node) {
	v2 := string(v)

	switch v2 {
	case "[":
		if checkName() {
			addNode(node, 1)
		}
	case "]":
		if checkName() {
			addNode(node, 2)
		}
	case ",":
		if checkName() {
			addNode(node, 0)
		}
	default:
		if syntax.IsWordChar(v) {
			scan.name += string(v) // Accept only alphanumeric (No escapes)
		}
	}
}

// Add a node.
// downUpOrNext input:
// 	- 0 for next sibling
// 	- 1 for lower level		[
//  - 2 for upper level		]
func addNode(curNode **node, downUpOrNext int) {
		var newNode *node

		newNode = &node{}
		newNode.Name = scan.name
		newNode.Parent = *curNode // nest it, depth+=1

		// add created newNode to Children of Parent
	(*curNode).Children = append((*curNode).Children, newNode)

		if downUpOrNext == 1 {
			*curNode = newNode
		} else if downUpOrNext == 2 {
			*curNode = (*curNode).Parent
		}

		scan.name = ""
}

func parse(v string) (*node, error) {
	root := node{}
	root.Name = ""
	root.Children = make([]*node, 0)

	var currentNode = &root // node pointer

	for _, v := range v {
		//curChar := string(v)

		parseNext(v, &currentNode)
	}

	return &root, nil
}


