package main

import (
	"encoding/json"
	"log"
	"regexp/syntax"
)

// Examples for unit tests.
var examples = []string{
	"[a,b,c]",
	"[a[aa[aaa],ab,ac],b,c[ca,cb,cc[cca]]]",
	"[a,b, c ]]]]]]",
	" ",
	"◌́",			// fruit unicode, acute "\u0301"
	"123]]]",
	"\n\n",			// Newlines
	"asd",
	"]]][a,b, c ]]]]]]",
	//"",			// Crash the old solution.
}

// Scan - the scanner structure.
var scan struct {
	name string		// Child name
	depth int		// Depth
	err string		// Error
}

type node struct {
	Name     string  `json:"name"`
	Parent   *node   `json:"-"`		// Parent pointer - depth iteration
	Children []*node `json:"children,omitempty"`
}

// Parse next rune as string
func parseNext(v rune, node *node) {
	v2 := string(v)

	switch v2 {
	case "[":
		if checkName() {
			addNode(node, 1)
			scan.depth++
		}
	case "]":
		if checkName() {
			addNode(node, 2)
			scan.depth--
		}
	case ",":
		if checkName() {
			addNode(node, 0)
		}
	default:
		if syntax.IsWordChar(v) {
			scan.name += string(v)		// Accept only alphanumeric (No escapes)
		}
	}
}

func checkName() bool {
	if len(scan.name) > 0 {
		return true
	} else {
		return false
	}
}

// Add a node.
// downUpOrNext input is,
// 	- 0 for next sibling (do nothing),
// 	- 1 for lower level
//  - 2 for upper level
func addNode(curNode *node, downUpOrNext int) {
	var newNode *node

	newNode = &node{}
	newNode.Name = scan.name
	newNode.Parent = curNode	// nest it, depth+=1

	// add created newNode to Children of Parent
	curNode.Children = append(curNode.Children, newNode)

	if downUpOrNext == 1 {
		curNode = newNode
	} else if downUpOrNext == 2 {
		curNode = curNode.Parent
	}

	scan.name = ""
}

func parse(v string) (*node, error) {
	root := node{}
	root.Name = ""
	root.Children = make([]*node, 0)


	scan.depth = 0

	var currentNode = &root	// current node, start with root

	for _, v := range v {
		//curChar := string(v)

		parseNext(v, currentNode)
	}

	return &root, nil
}

func main() {
	log.Print("\n\n========== SOLUTION 1 =========\n\n\n")
	for i, example := range examples {
		result, err := parse(example)
		if err != nil {
			panic(err)
		}

		// Start tests for the YouFoodz implementation.
		x, err := json.MarshalIndent(result, " ", " ")
		if err != nil {
			panic(err)
		}
		log.Printf("Example %d: %s - %s", i, example, string(x))

		// Start tests for the identitii implementation.
		//time.Sleep(3)
		//test()
	}

	// Start tests for the identitii implementation.
	log.Print("\n\n========== SOLUTION 2 =========\n\n\n")
	test()
}
