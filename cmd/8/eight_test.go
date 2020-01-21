package eight

import "testing"

func TestNewImage(t *testing.T) {
	var testCases = []struct {
		input  string
		width  int
		height int
	}{
		{
			input:  `123456789012`,
			width:  3,
			height: 2,
		},
	}

	for _, c := range testCases {
		_, err := NewImage(c.input, c.width, c.height)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
}
