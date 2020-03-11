package ten

import (
	"advent-of-code-2019/logger"
	"bufio"
	"bytes"
	"math"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var input = `
.#..#..##.#...###.#............#.
.....#..........##..#..#####.#..#
#....#...#..#.......#...........#
.#....#....#....#.#...#.#.#.#....
..#..#.....#.......###.#.#.##....
...#.##.###..#....#........#..#.#
..#.##..#.#.#...##..........#...#
..#..#.......................#..#
...#..#.#...##.#...#.#..#.#......
......#......#.....#.............
.###..#.#..#...#..#.#.......##..#
.#...#.................###......#
#.#.......#..####.#..##.###.....#
.#.#..#.#...##.#.#..#..##.#.#.#..
##...#....#...#....##....#.#....#
......#..#......#.#.....##..#.#..
##.###.....#.#.###.#..#..#..###..
#...........#.#..#..#..#....#....
..........#.#.#..#.###...#.....#.
...#.###........##..#..##........
.###.....#.#.###...##.........#..
#.#...##.....#.#.........#..#.###
..##..##........#........#......#
..####......#...#..........#.#...
......##...##.#........#...##.##.
.#..###...#.......#........#....#
...##...#..#...#..#..#.#.#...#...
....#......#.#............##.....
#......####...#.....#...#......#.
...#............#...#..#.#.#..#.#
.#...#....###.####....#.#........
#.#...##...#.##...#....#.#..##.#.
.#....#.###..#..##.#.##...#.#..##
`

// Ten command for ten
var Ten = &cobra.Command{
	Use:   "ten",
	Short: "day 10",
	Run:   ten,
}

func ten(cmd *cobra.Command, args []string) {
	asteroids := parseInput(input)

	coord, visible := findBestAsteroid(asteroids)
	logger.Infof("(part 1): coordinate: (%d, %d). visible asteroids: %d", coord.x, coord.y, len(visible))
	destroyed := destroyAsteroids(coord, visible)
	logger.Infof("(part 2): 200th asteroid: %+v. arbitrary math on it: %d", destroyed[199], destroyed[199].x*100+destroyed[199].y)
}

type coord struct {
	x int
	y int
}

// parses input and gives coordinates of asteroids
func parseInput(input string) []coord {
	res := make([]coord, 0)

	buf := bufio.NewScanner(bytes.NewBufferString(input))
	buf.Split(bufio.ScanLines)

	y := 0
	for buf.Scan() {
		txt := strings.TrimSpace(buf.Text())
		if len(txt) > 0 {
			for x, char := range txt {
				if string(char) == "#" {
					res = append(res, coord{x: x, y: y})

				}
			}
			y++
		}
	}
	return res
}

func findBestAsteroid(asteroids []coord) (coord, []coord) {
	max := 0
	var station coord
	var visible []coord
	for i, c := range asteroids {
		// have to copy array to modify it
		others := make([]coord, 0)
		others = append(others, asteroids[:i]...)
		others = append(others, asteroids[i+1:]...)
		sort.Slice(others, func(i, j int) bool {
			iDist := getDistance(c, others[i])
			jDist := getDistance(c, others[j])
			return iDist < jDist
		})
		numVisible := getVisibleAsteroids(c, others)
		if len(numVisible) > max {
			max = len(numVisible)
			station = c
			visible = numVisible
		}
	}

	return station, visible
}

func destroyAsteroids(station coord, asteroids []coord) []coord {
	//https://stackoverflow.com/questions/6989100/sort-points-in-clockwise-order
	asteroidsCopy := asteroids
	// sort by slope and then distance
	sort.Slice(asteroidsCopy, func(i, j int) bool {
		if asteroids[i].x-station.x >= 0 && asteroids[j].x-station.x < 0 {
			return true
		}
		if asteroids[i].x-station.x < 0 && asteroids[j].x-station.x >= 0 {
			return false
		}
		if asteroids[i].x-station.x == 0 && asteroids[j].x-station.x == 0 {
			if asteroids[i].y-station.y >= 0 || asteroids[j].y-station.y > 0 {
				return asteroids[i].y > asteroids[j].y
			}
			return asteroids[j].y > asteroids[i].y
		}

		// get the z magnitude of the cross product to determine
		// where the points are in relation to each other
		// right hand rule baby
		crossProductZ := ((asteroids[i].x - station.x) * (asteroids[j].y - station.y)) - ((asteroids[i].y - station.y) * (asteroids[j].x - station.x))
		// if crossProductZ is positive, i is to the right of j
		if crossProductZ > 0 {
			return true
		} else if crossProductZ == 0 { // if 0, they are on the same line
			iDistance := getDistance(station, asteroidsCopy[i])
			jDistance := getDistance(station, asteroidsCopy[j])

			return iDistance < jDistance
		} else {
			return false
		}
	})

	return asteroidsCopy
}

func getDistance(a, b coord) float64 {
	return math.Sqrt(float64(((b.x - a.x) * (b.x - a.x)) + ((b.y - a.y) * (b.y - a.y))))
}

// expect others to be ordered by distance
func getVisibleAsteroids(current coord, others []coord) []coord {
	visible := make([]coord, 0)
	for _, c := range others {

		blocked := false
		for _, v := range visible {
			// get line between visible point + current
			var start coord
			var end coord
			// always start with lower x
			if v.x < current.x {
				start = v
				end = current
			} else {
				start = current
				end = v
			}

			// don't have to be functions, so slope could be NaN. check for x/y equality
			if start.x == end.x {
				// if all xs same, may be blocked
				if c.x == start.x {
					// check if v is between c + current
					if (c.y > v.y && v.y > current.y) || (c.y < v.y && v.y < current.y) {
						blocked = true
						break

					}
				}
			} else { // otherwise find slope and check if point is on line
				slope := float64(end.y-start.y) / float64(end.x-start.x)
				intercept := float64(start.y) - slope*float64(start.x)
				// check if current is on the line, to arbitrary two decimal places b/c floating points
				calcY := float64(c.x)*slope + float64(intercept)
				rounded := math.Round(calcY*1000) / 1000
				if float64(c.y) == rounded {
					// see if its blocked by visible item
					if (c.x > v.x && v.x > current.x) || (c.x < v.x && v.x < current.x) {
						blocked = true
						break
					}
				}
			}

		}
		if !blocked {
			visible = append(visible, c)
		}
	}
	return visible
}
