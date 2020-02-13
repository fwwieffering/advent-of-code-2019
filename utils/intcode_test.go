package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"strings"
)

func TestProgram(t *testing.T) {
	var testCases = []struct {
		input  []int
		memory string
		result string
	}{
		{
			memory: "1,0,0,0,99",
			result: "2,0,0,0,99",
		},
		{
			memory: "2,3,0,3,99",
			result: "2,3,0,6,99",
		},
		{
			memory: "2,4,4,5,99,0",
			result: "2,4,4,5,99,9801",
		},
		{
			memory: "1,1,1,4,99,5,6,0,99",
			result: "30,1,1,4,2,5,6,0,99",
		},
		{
			memory: "0001,1,1,4,99,5,6,0,99",
			result: "30,1,1,4,2,5,6,0,99",
		},
		{
			memory: "3,3,99,-1",
			input:  []int{4444},
			result: "3,3,99,4444",
		},
		{
			memory: "1001,5,-1890,6,99,5,-1",
			result: "1001,5,-1890,6,99,5,-1885",
		},
		{
			memory: "101,20,1,5,99,-1",
			result: "101,20,1,5,99,40",
		},
		{
			memory: "108,2,5,6,99,2,-1",
			result: "108,2,5,6,99,2,1",
		},
	}

	for _, c := range testCases {
		fmt.Printf("\n--------------\nNEW TEST\n--------------\n%+v\n", c.memory)
		p, err := NewProgram(c.memory)
		if err != nil {
			t.Fatalf(err.Error())
		}
		res, err := p.Run(c.input...)
		if err != nil {
			t.Fatalf("\n%+v\n%s\n", res, err.Error())
		}
		out := IntArrayToString(res.Memory)
		if out != c.result {
			t.Fatalf("\nInput:\n%s\nOuput:\n%s\nExpected:\n%s\n", c.memory, out, c.result)
		}
	}
}

func TestOutput(t *testing.T) {
	testCases := []struct {
		mem    string
		input  []int
		output []int
	}{
		{
			mem:    "3,9,8,9,10,9,4,9,99,-1,8",
			output: []int{1},
			input:  []int{8},
		},
		{
			mem:    "3,9,8,9,10,9,4,9,99,-1,8",
			output: []int{0},
			input:  []int{18},
		},
		{
			mem:    "3,9,7,9,10,9,4,9,99,-1,8",
			output: []int{1},
			input:  []int{7},
		},
		{
			mem:    "3,9,7,9,10,9,4,9,99,-1,8",
			output: []int{0},
			input:  []int{17},
		},
		{
			mem:    "3,3,1108,-1,8,3,4,3,99",
			output: []int{1},
			input:  []int{8},
		},
		{
			mem:    "3,3,1108,-1,8,3,4,3,99",
			output: []int{0},
			input:  []int{11},
		},
		{
			mem:    "3,3,1107,-1,8,3,4,3,99",
			output: []int{1},
			input:  []int{6},
		},
		{
			mem:    "3,3,1107,-1,8,3,4,3,99",
			output: []int{0},
			input:  []int{1020},
		},
		{
			mem:    "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9",
			output: []int{1},
			input:  []int{2},
		},
		{
			mem:    "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9",
			output: []int{0},
			input:  []int{0},
		},
		{
			mem:    "3,3,1105,-1,9,1101,0,0,12,4,12,99,1",
			output: []int{1},
			input:  []int{2},
		},
		{
			mem:    "3,3,1105,-1,9,1101,0,0,12,4,12,99,1",
			output: []int{0},
			input:  []int{0},
		},
		{
			mem:    "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99",
			output: []int{1000},
			input:  []int{8},
		},
		{
			mem:    "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99",
			output: []int{999},
			input:  []int{7},
		},
		{
			mem:    "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99",
			output: []int{1001},
			input:  []int{9},
		},
		{
			mem:    "104,1125899906842624,99",
			output: []int{1125899906842624},
		},
		{
			mem:    "1102,34915192,34915192,7,4,7,99,0",
			output: []int{1219070632396864},
		},
		{
			mem:    "109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99",
			output: []int{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99},
		},
	}

	for _, c := range testCases {
		stringMem := strings.Split(c.mem, ",")
		intMem := make([]int, len(stringMem))

		for i, s := range stringMem {
			intval, _ := strconv.Atoi(s)
			intMem[i] = intval
		}
		t.Log(c.mem)
		p := Program{InitialMemory: intMem}
		output, _ := p.Run(c.input...)
		if !reflect.DeepEqual(output.Result, c.output) {
			t.Fatalf("%+v\n", output.Result)
		}
	}
}
