package kaanparser

import (
	"encoding/json"
	"log"
	"testing"
)

// +build gofuzz

// 1-) Start tests for the YouFoodz implementation.
func Test1(t *testing.T) {
	log.Print("\n\n========== SOLUTION 1 =========\n\n\n")
	for i, example := range examples {
		result, err := parse(example)
		if err != nil {
			panic(err)
		}


		x, err := json.MarshalIndent(result, " ", " ")
		if err != nil {
			panic(err)
		}
		log.Printf("Example %d: %s - %s", i, example, string(x))
	}
}

// 2-) Start tests for the identitii implementation.
func Test2(t *testing.T) {
	log.Print("\n\n========== SOLUTION 2 =========\n\n\n")
	testOld()
}

