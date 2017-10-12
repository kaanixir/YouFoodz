package main

import (
	"encoding/json"
	"log"
)

// Examples for unit tests.
var examples = []string{
	"[a,b,c]",
	"[a[aa[aaa],ab,ac],b,c[ca,cb,cc[cca]]]",
}

// The states.
const (
	// Continue.
	scanContinue     = iota // default content (alphanumeric etc..)
	scanBeginLiteral        // end implied by next result != scanContinue
	scanBeginObject         // begin object
	scanObjectKey           // just finished object key (string)
	scanObjectValue         // just finished non-last object value
	scanEndObject           // end object (implies scanObjectValue if possible)
	scanBeginArray          // begin array
	scanArrayValue          // just finished array value
	scanEndArray            // end array (implies scanArrayValue if possible)
	scanSkipSpace           // space byte; can skip; known to be last "continue" result

	// Stop.
	scanEnd   // top-level value ended *before* this byte; known to be first "stop" result
	scanError // hit an error, scanner.err.
)

const (
	parseObjectKey   = iota // parsing object key (before colon)
	parseObjectValue        // parsing object value (after colon)
	parseArrayValue         // parsing array value
)

// Kaan, the scanner.
type kaan struct {
	// The step is a func to be called to execute the next transition.
	step func(*kaan, byte) int

	// End of this level - []
	endTop bool

	// Stack of states.
	parseState []int

	// Errors
	err error

	// 1-byte redo (see undo method)
	redo      bool
	redoCode  int
	redoState func(*kaan, byte) int

	// total bytes consumed, updated by decoder.Decode
	bytes int64
}


// ========= PUSH&POP PARSE STATES ============ //
// pushParseState pushes a new parse state p onto the parse stack.
func (s *kaan) pushParseState(p int) {
	s.parseState = append(s.parseState, p)
}
// popParseState pops a parse state (already obtained) off the stack
// and updates s.step accordingly.
func (s *kaan) popParseState() {
	n := len(s.parseState) - 1
	s.parseState = s.parseState[0:n]
	s.redo = false
	if n == 0 {
		s.step = stateEndTop
		s.endTop = true
	} else {
		s.step = stateEndValue		// TODO: stateEndValue  == stateEndAll
	}
}// ========= PUSH&POP PARSE STATES ============ //


// eof tells the scanner that the end of input has been reached.
// It returns a scan status just as s.step does.
func (s *kaan) eof() int {
	if s.err != nil {
		return scanError
	}
	if s.endTop {
		return scanEnd
	}
	s.step(s, ' ')
	if s.endTop {
		return scanEnd
	}
	if s.err == nil {
		s.err = &SyntaxError{"unexpected end of JSON input", s.bytes}
	}
	return scanError
}

// ======= ERRORS ========== //
// A SyntaxError is a description of a JSON syntax error.
type SyntaxError struct {
	msg    string // description of error
	Offset int64  // error occurred after reading Offset bytes
}
func (e *SyntaxError) Error() string { return e.msg }
// error records an error and switches to the error state.
func (s *kaan) error(c byte, context string) int {
	s.step = stateError
	s.err = &SyntaxError{"invalid character " + quoteChar(c) + " " + context, s.bytes}
	return scanError
}// ======= ERRORS ========== //















































// ===================================== //

type node struct {
	Name     string  `json:"name"`
	Parent   *node   `json:"-"`		// Extra parent pointer placeholder - depth iteration
	Children []*node `json:"children,omitempty"`
}


func parse(v string) (*node, error) {
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

func main() {
	for i, example := range examples {
		result, err := parse(example)
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
