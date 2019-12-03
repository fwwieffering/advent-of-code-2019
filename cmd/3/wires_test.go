package three

import "testing"

func TestWires(t *testing.T) {
	var testCases = []struct {
		wire1    string
		wire2    string
		distance int
	}{
		{
			wire1:    "R75,D30,R83,U83,L12,D49,R71,U7,L72",
			wire2:    "U62,R66,U55,R34,D71,R55,D58,R83",
			distance: 159,
		},
		{
			wire1:    "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51",
			wire2:    "U98,R91,D20,R16,D67,R40,U7,R15,U6,R7",
			distance: 135,
		},
	}

	for _, c := range testCases {
		wire1, err := createWire(c.wire1)
		if err != nil {
			t.Fatalf(err.Error())
		}
		wire2, err := createWire(c.wire2)
		if err != nil {
			t.Fatalf(err.Error())
		}

		res := getNearestIntersection(wire1, wire2)
		if res != c.distance {
			t.Fatalf("Wire1: %s\nWire2: %s\nExpected: %d\nGot: %d\n", c.wire1, c.wire2, c.distance, res)
		}
	}
}

func TestWireLatency(t *testing.T) {
	var testCases = []struct {
		wire1    string
		wire2    string
		distance int
	}{
		{
			wire1:    "R75,D30,R83,U83,L12,D49,R71,U7,L72",
			wire2:    "U62,R66,U55,R34,D71,R55,D58,R83",
			distance: 610,
		},
		{
			wire1:    "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51",
			wire2:    "U98,R91,D20,R16,D67,R40,U7,R15,U6,R7",
			distance: 410,
		},
	}

	for _, c := range testCases {
		wire1, err := createWire(c.wire1)
		if err != nil {
			t.Fatalf(err.Error())
		}
		wire2, err := createWire(c.wire2)
		if err != nil {
			t.Fatalf(err.Error())
		}

		res := getLowestLatencyIntersection(wire1, wire2)
		if res != c.distance {
			t.Fatalf("Wire1: %s\nWire2: %s\nExpected: %d\nGot: %d\n", c.wire1, c.wire2, c.distance, res)
		}
	}
}
