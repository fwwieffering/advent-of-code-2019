package two

import "testing"

func TestProgram(t *testing.T) {
	var testCases = []struct {
		in  string
		out string
	}{
		{
			in:  "1,0,0,0,99",
			out: "2,0,0,0,99",
		},
		{
			in:  "2,3,0,3,99",
			out: "2,3,0,6,99",
		},
		{
			in:  "2,4,4,5,99,0",
			out: "2,4,4,5,99,9801",
		},
		{
			in:  "1,1,1,4,99,5,6,0,99",
			out: "30,1,1,4,2,5,6,0,99",
		},
	}

	for _, c := range testCases {
		out, _ := RunProgramString(c.in)
		if out != c.out {
			t.Fatalf("Input:\n%s\nOuput:\n%s\nExpected:\n%s\n", c.in, out, c.out)
		}
	}
}
