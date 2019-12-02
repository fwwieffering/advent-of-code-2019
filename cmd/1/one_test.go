package one

import "testing"

func TestGetFuelReq(t *testing.T) {
	var testCases = []struct {
		in  int
		out int
	}{
		{
			in:  12,
			out: 2,
		},
		{
			in:  14,
			out: 2,
		},
		{
			in:  1969,
			out: 654,
		},
		{
			in:  100756,
			out: 33583,
		},
	}

	for _, c := range testCases {
		res := getFuleReq(c.in)
		if res != c.out {
			t.Fatalf("getFuleReq(%d) should equal %d. Got: %d", c.in, c.out, res)
		}
	}
}
