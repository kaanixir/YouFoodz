// +build gofuzz

package kaanparser

import "log"

// Fuzz tests.
func Fuzz(data []byte) int {
	log.Print("\n\n========== FUZZING =========\n\n\n")
	result, err := parse(string(data))
	if result == nil && err == nil {
		panic("Out nil and err nil!")
	}

	return 0
}