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
// --- Part Two ---
// During the second Go / No Go poll, the Elf in charge of the Rocket Equation Double-Checker stops the launch sequence. Apparently, you forgot to include additional fuel for the fuel you just added.
// Fuel itself requires fuel just like a module - take its mass, divide by three, round down, and subtract 2. However, that fuel also requires fuel, and that fuel requires fuel, and so on. Any mass that would require negative fuel should instead be treated as if it requires zero fuel; the remaining mass, if any, is instead handled by wishing really hard, which has no mass and is outside the scope of this calculation.
// So, for each module mass, calculate its fuel and add it to the total. Then, treat the fuel amount you just calculated as the input mass and repeat the process, continuing until a fuel requirement is zero or negative. For example:
//     A module of mass 14 requires 2 fuel. This fuel requires no further fuel (2 divided by 3 and rounded down is 0, which would call for a negative fuel), so the total fuel required is still just 2.
//     At first, a module of mass 1969 requires 654 fuel. Then, this fuel requires 216 more fuel (654 / 3 - 2). 216 then requires 70 more fuel, which requires 21 fuel, which requires 5 fuel, which requires no further fuel. So, the total fuel required for a module of mass 1969 is 654 + 216 + 70 + 21 + 5 = 966.
//     The fuel required by a module of mass 100756 and its fuel is: 33583 + 11192 + 3728 + 1240 + 411 + 135 + 43 + 12 + 2 = 50346.
// What is the sum of the fuel requirements for all of the modules on your spacecraft when also taking into account the mass of the added fuel? (Calculate the fuel requirements for each module separately, then add them all up at the end.)
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
	totalIncludingFuel := 0
	for scanner.Scan() {
		txt := scanner.Text()
		i, err := strconv.Atoi(txt)
		if err != nil {
			logger.Fatalf("Bad input: unable to convert %s to integer: %s", txt, err.Error())
		}
		total += getFuelReq(i)
		totalIncludingFuel += getFuelReqIncludingFuel(i)
	}
	logger.Infof("(part one) Total Fuel requirement for modules, not including the fuel requirement for the fuel itself: %d", total)
	logger.Infof("(part two) Total Fuel requirement for modules: %d", totalIncludingFuel)
}

// Fuel required to launch a given module is based on its mass. Specifically, to find the fuel required for a module, take its mass, divide by three, round down, and subtract 2.
// For example:
//     For a mass of 12, divide by 3 and round down to get 4, then subtract 2 to get 2.
//     For a mass of 14, dividing by 3 and rounding down still yields 4, so the fuel required is also 2.
//     For a mass of 1969, the fuel required is 654.
//     For a mass of 100756, the fuel required is 33583.
func getFuelReq(i int) int {
	// golang integer division automatically rounds down
	divideThree := i / 3
	return divideThree - 2
}

// So, for each module mass, calculate its fuel and add it to the total. Then, treat the fuel amount you just calculated as the input mass and repeat the process, continuing until a fuel requirement is zero or negative. For example:
//     A module of mass 14 requires 2 fuel. This fuel requires no further fuel (2 divided by 3 and rounded down is 0, which would call for a negative fuel), so the total fuel required is still just 2.
//     At first, a module of mass 1969 requires 654 fuel. Then, this fuel requires 216 more fuel (654 / 3 - 2). 216 then requires 70 more fuel, which requires 21 fuel, which requires 5 fuel, which requires no further fuel. So, the total fuel required for a module of mass 1969 is 654 + 216 + 70 + 21 + 5 = 966.
//     The fuel required by a module of mass 100756 and its fuel is: 33583 + 11192 + 3728 + 1240 + 411 + 135 + 43 + 12 + 2 = 50346
func getFuelReqIncludingFuel(moduleWeight int) int {
	var total = 0

	var current = getFuelReq(moduleWeight)

	for current > 0 {
		total += current
		current = getFuelReq(current)
	}

	return total
}
