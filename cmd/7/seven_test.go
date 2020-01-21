package seven

import (
	"advent-of-code-2019/utils"
	"fmt"
	"reflect"
	"testing"
)

func TestDeleteItem(t *testing.T) {
	testCases := []struct {
		in    []int
		index int
		out   []int
	}{
		{
			in:    []int{0, 1, 2, 3, 4},
			index: 0,
			out:   []int{1, 2, 3, 4},
		},
		{
			in:    []int{0, 1, 2, 3, 4},
			index: 4,
			out:   []int{0, 1, 2, 3},
		},
		{
			in:    []int{0, 1, 2, 3, 4},
			index: 2,
			out:   []int{0, 1, 3, 4},
		},
		{
			in:    []int{4},
			index: 2,
			out:   []int{},
		},
		{
			in:    []int{2, 4},
			index: 0,
			out:   []int{4},
		},
		{
			in:    []int{2, 4},
			index: 1,
			out:   []int{2},
		},
	}

	for _, c := range testCases {
		res := deleteItem(c.in, c.index)
		if !reflect.DeepEqual(res, c.out) {
			t.Fatalf("\nres: %+v\ntest: %+v", res, c)
		}
	}
}

func TestCombinations(t *testing.T) {
	res := generateCombinations([]int{0, 1, 2, 3, 4}, 5)
	positiveChecks := [][]int{[]int{4, 3, 2, 1, 0}, []int{0, 1, 2, 3, 4}, []int{1, 0, 4, 3, 2}}
	for _, check := range positiveChecks {
		found := false
		printStr := ""
		for _, item := range res {
			printStr += fmt.Sprintf("%+v\n", item)
			if reflect.DeepEqual(item, check) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("\n%s\n%+v was not found", printStr, check)
		}
	}
	negativeChecks := [][]int{[]int{4, 4, 4, 4, 4}, []int{4, 4, 2, 2, 1}}
	for _, check := range negativeChecks {
		found := false
		printStr := ""
		for _, item := range res {
			printStr += fmt.Sprintf("%+v\n", item)
			if reflect.DeepEqual(item, check) {
				found = true
				break
			}
		}
		t.Log(found)
		if found {
			t.Fatalf("%s\n%+v was found", printStr, check)
		}
	}

	if len(res) == 0 {
		t.Fatalf("should be at least one combo: %+v", res)
	}
}

func TestAmplifiers(t *testing.T) {
	testCases := []struct {
		program          string
		expectedSequence []int
		result           int
	}{
		{
			program:          "3,15,3,16,1002,16,10,16,1,16,15,15,4,15,99,0,0",
			expectedSequence: []int{4, 3, 2, 1, 0},
			result:           43210,
		},
		{
			program:          "3,23,3,24,1002,24,10,24,1002,23,-1,23,101,5,23,23,1,24,23,23,4,23,99,0,0",
			expectedSequence: []int{0, 1, 2, 3, 4},
			result:           54321,
		},
		{
			program:          "3,31,3,32,1002,32,10,32,1001,31,-2,31,1007,31,0,33,1002,33,7,33,1,33,31,31,1,32,31,31,4,31,99,0,0,0",
			expectedSequence: []int{1, 0, 4, 3, 2},
			result:           65210,
		},
	}

	for _, c := range testCases {
		prog, err := utils.NewProgram(c.program)
		if err != nil {
			t.Fatalf(err.Error())
		}
		sequenceRes, err := runSequence(prog, c.expectedSequence)
		if sequenceRes != c.result {
			t.Fatalf("\n%s\nexpected: %d when running %+v but got: %d", c.program, c.result, c.expectedSequence, sequenceRes)
		}

		max, sequence, err := testAmplifiers(prog, 5)
		if err != nil {
			t.Fatalf("error running program %s: %s", c.program, err.Error())
		}
		if max != c.result || !reflect.DeepEqual(sequence, c.expectedSequence) {
			t.Fatalf("\n%s\nExpected max: %d sequence: %+v\nGot max: %d sequence %+v", c.program, c.result, c.expectedSequence, max, sequence)
		}
	}
}

func TestAmplifierFeedback(t *testing.T) {
	testCases := []struct {
		program          string
		expectedSequence []int
		result           int
	}{
		{
			program:          "3,26,1001,26,-4,26,3,27,1002,27,2,27,1,27,26,27,4,27,1001,28,-1,28,1005,28,6,99,0,0,5",
			expectedSequence: []int{9, 8, 7, 6, 5},
			result:           139629729,
		},
		{
			program:          "3,52,1001,52,-5,52,3,53,1,52,56,54,1007,54,5,55,1005,55,26,1001,54,-5,54,1105,1,12,1,53,54,53,1008,54,0,55,1001,55,1,55,2,53,55,53,4,53,1001,56,-1,56,1005,56,6,99,0,0,0,0,10",
			expectedSequence: []int{9, 7, 8, 5, 6},
			result:           18216,
		},
	}

	for _, c := range testCases {
		prog, err := utils.NewProgram(c.program)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// some iterations have had issues with concurrency.
		// ensure its repeatable
		for i := 0; i < 99; i++ {
			fmt.Printf("--------------------------------- %d -------------------------------\n", i)
			sequenceRes, err := runSequenceFeedback(prog, c.expectedSequence)
			if err != nil {
				t.Fatalf("error running sequence: %s", err.Error())
			}
			if sequenceRes != c.result {
				t.Fatalf("\n%s\nexpected: %d when running %+v #%d but got: %d", c.program, c.result, c.expectedSequence, i, sequenceRes)
			}
			max, sequence, err := testAmplifiersFeedback(prog, 5)
			if err != nil {
				t.Fatalf("error running program %s: %s", c.program, err.Error())
			}
			if max != c.result || !reflect.DeepEqual(sequence, c.expectedSequence) {
				t.Fatalf("\n%s\nExpected max: %d sequence: %+v #%d\nGot max: %d sequence %+v", c.program, c.result, c.expectedSequence, i, max, sequence)
			}
		}
	}
}
