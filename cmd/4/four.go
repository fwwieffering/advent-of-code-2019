package four

import (
	"advent-of-code-2019/logger"

	"github.com/spf13/cobra"
)

// Four does day 4
var Four = &cobra.Command{
	Use:   "four",
	Short: "runs day 4 of advent of code",
	Run:   four,
}

var rangeLow = 171309

var rangeHigh = 643603

func findMatches(alg func(int) bool) []int {
	res := make([]int, 0)
	for i := rangeLow; i <= rangeHigh; i++ {
		possible := alg(i)
		if possible {
			res = append(res, i)
		}
	}
	return res
}

// pw finder
// It is a six-digit number.
// The value is within the range given in your puzzle input.
// Two adjacent digits are the same (like 22 in 122345).
// Going from left to right, the digits never decrease; they only ever increase or stay the same (like 111123 or 135679).

// Other than the range rule, the following are true:

// 111111 meets these criteria (double 11, never decreases).
// 223450 does not meet these criteria (decreasing pair of digits 50).
// 123789 does not meet these criteria (no double).
func checkMatchPart1(i int) bool {
	repeats := false
	increases := true

	prevDigit := -1

	// iterate over digits in number from right to left
	for i > 0 {
		currentDigit := i % 10
		if currentDigit == prevDigit {
			repeats = true
		}
		// check if prevDigit (digit to the right) is smaller, and if so break
		if prevDigit < currentDigit && prevDigit != -1 {
			increases = false
			break
		}

		// move on to next number
		i = i / 10
		prevDigit = currentDigit
	}

	return repeats && increases
}

// returns the number of times the rightmost digit repeats
// always at least 1
func numRepeats(i int) int {
	if i <= 10 {
		return 1
	}
	matcher := i % 10
	i = i / 10

	repeats := 1
	for i > 0 {
		currentDigit := i % 10
		if currentDigit == matcher {
			repeats++
		} else {
			break
		}
		i = i / 10
	}
	return repeats
}

// An Elf just remembered one more important detail: the two adjacent matching digits are not part of a larger group of matching digits.

// Given this additional criterion, but still ignoring the range rule, the following are now true:

//     112233 meets these criteria because the digits never decrease and all repeated digits are exactly two digits long.
//     123444 no longer meets the criteria (the repeated 44 is part of a larger group of 444).
//     111122 meets the criteria (even though 1 is repeated more than twice, it still contains a double 22).
func checkMatchPart2(i int) bool {
	repeats := false
	increases := true

	prevDigit := -1
	// prevMatches := -1
	// digitTwoBack := -2
	// iterate over digits in number from right to left
	for i > 0 {
		currentDigit := i % 10

		repeatsToTheLeft := numRepeats(i)

		if currentDigit != prevDigit && repeatsToTheLeft == 2 {
			repeats = true
		}
		// check if prevDigit (digit to the right) is smaller, and if so break
		if prevDigit < currentDigit && prevDigit != -1 {
			increases = false
			break
		}

		// move on to next number
		i = i / 10
		// digitTwoBack = prevDigit
		prevDigit = currentDigit
	}

	return repeats && increases

}

func four(cmd *cobra.Command, args []string) {
	part1Answer := findMatches(checkMatchPart1)
	part2Answer := findMatches(checkMatchPart2)

	logger.Infof("(part 1): %d passwords match the given criteria", len(part1Answer))
	logger.Infof("(part 2): %d passwords match the given criteria", len(part2Answer))
}
