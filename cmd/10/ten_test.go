package ten

import (
	"fmt"
	"testing"
)

type testCase struct {
	input      string
	station    coord
	numvisible int
}

var testCases = []testCase{
	{
		input: `
		.#..#
		.....
		#####
		....#
		...##`,
		station:    coord{x: 3, y: 4},
		numvisible: 8,
	},
	{
		input: `
		......#.#.
		#..#.#....
		..#######.
		.#.#.###..
		.#..#.....
		..#....#.#
		#..#....#.
		.##.#..###
		##...#..#.
		.#....####`,
		station:    coord{x: 5, y: 8},
		numvisible: 33,
	},
	{
		input: `
		#.#...#.#.
		.###....#.
		.#....#...
		##.#.#.#.#
		....#.#.#.
		.##..###.#
		..#...##..
		..##....##
		......#...
		.####.###.`,
		station:    coord{x: 1, y: 2},
		numvisible: 35,
	},
	{
		input: `
		.#..#..###
		####.###.#
		....###.#.
		..###.##.#
		##.##.#.#.
		....###..#
		..#.#..#.#
		#..#.#.###
		.##...##.#
		.....#.#..`,
		station:    coord{x: 6, y: 3},
		numvisible: 41,
	},
	{
		input: `
		.#..##.###...#######
		##.############..##.
		.#.######.########.#
		.###.#######.####.#.
		#####.##.#.##.###.##
		..#####..#.#########
		####################
		#.####....###.#.#.##
		##.#################
		#####.##.###..####..
		..######..##.#######
		####.##.####...##..#
		.#####..#.######.###
		##...#.##########...
		#.##########.#######
		.####.#.###.###.#.##
		....##.##.###..#####
		.#.#.###########.###
		#.#.#.#####.####.###
		###.##.####.##.#..##`,
		station:    coord{x: 11, y: 13},
		numvisible: 210,
	},
}

func TestCountVisibleAsteroids(t *testing.T) {
	for _, c := range testCases {
		asteroids := parseInput(c.input)
		station, visible := findBestAsteroid(asteroids)
		if station.x != c.station.x || station.y != c.station.y {
			t.Fatalf("%s\nexpected station to be: %+v with %d visible. Received station: %+v with %d visible\n%+v", c.input, c.station, c.numvisible, station, len(visible), visible)
		}
		if len(visible) != c.numvisible {
			t.Fatalf("Expected to receive %d visible stations. Got: %d", c.numvisible, len(visible))
		}
	}
}

func TestDestroyedAsteroids(t *testing.T) {
	input := `
	.#..##.###...#######
	##.############..##.
	.#.######.########.#
	.###.#######.####.#.
	#####.##.#.##.###.##
	..#####..#.#########
	####################
	#.####....###.#.#.##
	##.#################
	#####.##.###..####..
	..######..##.#######
	####.##.####...##..#
	.#####..#.######.###
	##...#.##########...
	#.##########.#######
	.####.#.###.###.#.##
	....##.##.###..#####
	.#.#.###########.###
	#.#.#.#####.####.###
	###.##.####.##.#..##`

	asteroids := parseInput(input)
	station, visible := findBestAsteroid(asteroids)
	destroyed := destroyAsteroids(station, visible)

	checkAsteroid := func(a, b coord) {
		if !(a.x == b.x && a.y == b.y) {
			fmt.Printf("%v\n", destroyed)
			t.Fatalf("%+v != %+v", a, b)
		}
	}

	checkAsteroid(destroyed[0], coord{x: 11, y: 12})
	checkAsteroid(destroyed[1], coord{x: 12, y: 1})
	checkAsteroid(destroyed[2], coord{x: 12, y: 2})
	checkAsteroid(destroyed[9], coord{x: 12, y: 8})
	checkAsteroid(destroyed[19], coord{x: 16, y: 0})
	checkAsteroid(destroyed[49], coord{x: 16, y: 9})
	checkAsteroid(destroyed[99], coord{x: 10, y: 16})
	checkAsteroid(destroyed[198], coord{x: 9, y: 6})
}
