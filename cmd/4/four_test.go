package four

import (
	"testing"
)

func TestCheckMatch(t *testing.T) {
	var testCases = []struct {
		pw    int
		match bool
	}{
		{
			pw:    111111,
			match: true,
		},
		{
			pw:    223450,
			match: false,
		},
		{
			pw:    123789,
			match: false,
		},
	}

	for _, c := range testCases {
		res := checkMatchPart1(c.pw)
		if res != c.match {
			t.Fatalf("\npw: %d\nExpected: %t\nGot: %t\n", c.pw, c.match, res)
		}
	}
}

func TestNumRepeats(t *testing.T) {
	var testCases = []struct {
		input  int
		output int
	}{
		{
			input:  3222,
			output: 3,
		},
		{
			input:  2224,
			output: 1,
		},
		{
			input:  4,
			output: 1,
		},
		{
			input:  11111,
			output: 5,
		},
	}
	for _, c := range testCases {
		res := numRepeats(c.input)
		if res != c.output {
			t.Fatalf("\ninput: %d\nExpected: %d\nGot: %d\n", c.input, c.output, res)
		}
	}
}

func TestCheckMatchPart2(t *testing.T) {
	var testCases = []struct {
		pw    int
		match bool
	}{
		{
			pw:    112233,
			match: true,
		},
		{
			pw:    123444,
			match: false,
		},
		{
			pw:    111122,
			match: true,
		},
		{
			pw:    112345,
			match: true,
		},
		{
			pw:    589999,
			match: false,
		},
		{
			pw:    588999,
			match: true,
		},
		{
			pw:    222222,
			match: false,
		},
	}

	for _, c := range testCases {
		res := checkMatchPart2(c.pw)
		if res != c.match {
			t.Fatalf("\npw: %d\nExpected: %t\nGot: %t\n", c.pw, c.match, res)
		}
	}

}
