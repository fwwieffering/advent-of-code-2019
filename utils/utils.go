package utils

import (
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"strings"
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

func StringToIntArray(commasepstr string) ([]int, error) {
	trimMem := strings.TrimSpace(commasepstr)
	splitMem := strings.Split(trimMem, ",")
	intMem := make([]int, len(splitMem))
	for i, s := range splitMem {
		intval, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("couldn't convert %s to int at position %d", s, i)
		}
		intMem[i] = intval
	}
	return intMem, nil
}

func CopyIntArray(src []int) []int {
	dst := make([]int, len(src))
	copy(dst, src)
	return dst
}

func IntArrayToString(arr []int) string {
	res := ""
	for idx, digit := range arr {
		res += fmt.Sprintf("%d", digit)
		if idx != len(arr)-1 {
			res += ","
		}
	}
	return res
}
