package kaanparser

import (
	"encoding/json"
	"log"
)

func parseOld(v string) (*node, error) {
	root := node{}
	root.Name = ""
	root.Children = make([]*node, 0)

	var name string 		// node name
	var currentNode = &root	// current node, start with root

	for _, v := range v[1:] {	// skip first "["
		currentChar := string(v)
		var newNode *node

		// Iterate char.
		// Assume default case is ASCII tokens, adding to the name. Otherwise   "["  "]"  or  ","
		switch currentChar {
		case "[":				// Go down, depth + 1, new children
			// No children without names ( as with arrays of [[ or [[[ )
			if len(name) > 0 {
				newNode = &node{}
				newNode.Name = name
				newNode.Parent = currentNode	// nest it, depth+=1

				// add created newNode to Children of Parent
				currentNode.Children = append(currentNode.Children, newNode)

				currentNode = newNode	//
				name = "" // -- name reset
			}

		case "]":				// Go up, depth - 1,
			if len(name) > 0 {
				newNode = &node{}
				newNode.Name = name
				newNode.Parent = currentNode

				// add created newNode to Children of Parent
				currentNode.Children = append(currentNode.Children, newNode)

				// end nesting, go up depth
				currentNode = currentNode.Parent
				name = "" // -- name reset
			}

		case ",":				// Sibling children node creation
			if len(name) > 0 {
				newNode = &node{}
				newNode.Name = name
				newNode.Parent = currentNode
				currentNode.Children = append(currentNode.Children, newNode)

				name = "" // -- name reset
			}
		default:
			name += currentChar // add current character to the current name

		}
	}
	return &root, nil
}

func testOld() {
	for i, example := range examples {
		result, err := parseOld(example)
		if err != nil {
			panic(err)
		}

		j, err := json.MarshalIndent(result, " ", " ")
		if err != nil {
			panic(err)
		}
		log.Printf("Example %d: %s - %s", i, example, string(j))
	}
}
