package utils

import (
	"path"
	"path/filepath"
	"strconv"
)

// GetInputPath returns the absolute filepath to the input file
// as long as input is stored in the cmd/{day} module and named input
// executable also has to be run in directory root
func GetInputPath(day int) string {
	// read in file input
	cwd, _ := filepath.Abs("./")
	input := path.Join(cwd, "cmd", strconv.Itoa(day), "input")
	return input
}
