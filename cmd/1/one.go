package one

import (
	"advent-of-code-2019/logger"
	"bufio"
	"bytes"
	"strconv"

	"github.com/spf13/cobra"
)

var input = `51585
137484
73634
71535
87274
74243
127025
66829
138729
145459
118813
82326
82518
145032
148699
105958
103969
72689
145061
70385
53104
107851
103392
107051
123475
123918
56709
89284
86208
71943
109257
108272
124811
142709
115650
53607
142891
144135
114277
138671
111998
70838
69802
107210
103319
60377
58639
131863
100807
118360
52573
108207
128009
96180
148492
112914
72867
140991
131267
125123
58393
129615
87239
63085
59231
95007
147712
109838
89829
55634
96163
52323
106701
141511
125349
137267
50694
53692
57466
117769
63535
101708
113593
79163
112327
91994
129674
58076
145062
122730
102481
109994
136271
111178
117920
107933
104305
99613
68482
126543
`

// One produces the solution for advent of code problem one
// --- Day 1: The Tyranny of the Rocket Equation ---
// Santa has become stranded at the edge of the Solar System while delivering presents to other planets! To accurately calculate his position in space, safely align his warp drive, and return to Earth in time to save Christmas, he needs you to bring him measurements from fifty stars.
// Collect stars by solving puzzles. Two puzzles will be made available on each day in the Advent calendar; the second puzzle is unlocked when you complete the first. Each puzzle grants one star. Good luck!
// The Elves quickly load you into a spacecraft and prepare to launch.
// At the first Go / No Go poll, every Elf is Go until the Fuel Counter-Upper. They haven't determined the amount of fuel required yet.
// Fuel required to launch a given module is based on its mass. Specifically, to find the fuel required for a module, take its mass, divide by three, round down, and subtract 2.
// For example:
//     For a mass of 12, divide by 3 and round down to get 4, then subtract 2 to get 2.
//     For a mass of 14, dividing by 3 and rounding down still yields 4, so the fuel required is also 2.
//     For a mass of 1969, the fuel required is 654.
//     For a mass of 100756, the fuel required is 33583.
// The Fuel Counter-Upper needs to know the total fuel requirement. To find it, individually calculate the fuel needed for the mass of each module (your puzzle input), then add together all the fuel values.
// What is the sum of the fuel requirements for all of the modules on your spacecraft?
var One = &cobra.Command{
	Use:   "one",
	Short: "one runs advent of code problem one https://adventofcode.com/2019/day/1",
	Run:   one,
}

func one(cmd *cobra.Command, args []string) {
	buf := bytes.NewBufferString(input)
	scanner := bufio.NewScanner(buf)
	scanner.Split(bufio.ScanLines)

	total := 0
	for scanner.Scan() {
		txt := scanner.Text()
		i, err := strconv.Atoi(txt)
		if err != nil {
			logger.Fatalf("Bad input: unable to convert %s to integer: %s", txt, err.Error())
		}
		total += getFuleReq(i)
	}
	logger.Infof("Total Fuel requirement: %d", total)
}

// Fuel required to launch a given module is based on its mass. Specifically, to find the fuel required for a module, take its mass, divide by three, round down, and subtract 2.
// For example:
//     For a mass of 12, divide by 3 and round down to get 4, then subtract 2 to get 2.
//     For a mass of 14, dividing by 3 and rounding down still yields 4, so the fuel required is also 2.
//     For a mass of 1969, the fuel required is 654.
//     For a mass of 100756, the fuel required is 33583.
func getFuleReq(i int) int {
	// golang integer division automatically rounds down
	divideThree := i / 3
	return divideThree - 2
}
