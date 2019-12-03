package three

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type wire struct {
	input string
	set   *coordinateSet
}

type coordinateSet struct {
	coords []*coordinate
}

func (cs *coordinateSet) combine(other *coordinateSet) {
	cs.coords = append(cs.coords, other.coords...)
}

func (cs *coordinateSet) add(c *coordinate) {
	cs.coords = append(cs.coords, c)
}

func (cs *coordinateSet) contains(c *coordinate) bool {
	for _, i := range cs.coords {
		if i.x == c.x && i.y == c.y {
			return true
		}
	}
	return false
}

func getIntersections(set1, set2 *coordinateSet) *coordinateSet {
	res := &coordinateSet{
		coords: make([]*coordinate, 0),
	}

	hash := make(map[coordinate]int)

	for _, item := range set2.coords {
		hash[coordinate{
			x: item.x,
			y: item.y}] = item.stepnum
	}

	for _, c := range set1.coords {
		if othersteps, ok := hash[coordinate{x: c.x, y: c.y}]; ok {
			res.add(&coordinate{
				x:       c.x,
				y:       c.y,
				stepnum: c.stepnum + othersteps,
			})
		}
	}
	return res
}

func getNearestIntersection(w1, w2 *wire) int {
	intersections := getIntersections(w1.set, w2.set)
	lowest := -1
	for _, i := range intersections.coords {
		distance := i.manhattanDistance()
		if lowest == -1 || distance < lowest {
			lowest = distance
		}
	}
	return lowest
}

func getLowestLatencyIntersection(w1, w2 *wire) int {
	intersections := getIntersections(w1.set, w2.set)
	lowest := -1
	for _, i := range intersections.coords {
		latency := i.stepnum
		if lowest == -1 || latency < lowest {
			lowest = latency
		}
	}
	return lowest
}

type coordinate struct {
	x       int
	y       int
	stepnum int
}

func (c *coordinate) manhattanDistance() int {
	d := 0
	for _, i := range []int{c.x, c.y} {
		if i > 0 {
			d += i
		} else {
			d += -1 * i
		}
	}
	return d
}

type direction string

const (
	down  direction = "D"
	up    direction = "U"
	left  direction = "L"
	right direction = "R"
)

type instruction struct {
	direction direction
	distance  int
}

func createWire(input string) (*wire, error) {
	res := &wire{
		input: input,
	}

	// split on comma
	instructions := strings.Split(input, ",")
	currentCoordinate := &coordinate{
		x:       0,
		y:       0,
		stepnum: 0,
	}
	set := &coordinateSet{
		coords: make([]*coordinate, 0),
	}
	for idx, i := range instructions {
		cmd, err := parseInstruction(i)
		if err != nil {
			return nil, fmt.Errorf("index %d: %s", idx, err.Error())
		}
		resultantCoordinate, coordSet := processInstruction(cmd, currentCoordinate)
		currentCoordinate = resultantCoordinate
		set.combine(coordSet)
	}
	res.set = set
	return res, nil
}

// instructions of form {direction}{count}
func parseInstruction(input string) (*instruction, error) {
	exp, _ := regexp.Compile(`([RDLU])(\d+)`)
	matches := exp.FindStringSubmatch(input)
	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid instruction %s. Only %+v matches", input, matches)
	}
	distance, _ := strconv.Atoi(matches[2])
	var d direction
	switch matches[1] {
	case "D":
		d = down
	case "L":
		d = left
	case "U":
		d = up
	case "R":
		d = right
	}
	return &instruction{
		direction: d,
		distance:  distance,
	}, nil
}

func goUp(c *coordinate) *coordinate {
	return &coordinate{
		x:       c.x,
		y:       c.y + 1,
		stepnum: c.stepnum + 1,
	}
}

func goDown(c *coordinate) *coordinate {
	return &coordinate{
		x:       c.x,
		y:       c.y - 1,
		stepnum: c.stepnum + 1,
	}
}

func goLeft(c *coordinate) *coordinate {
	return &coordinate{
		x:       c.x - 1,
		y:       c.y,
		stepnum: c.stepnum + 1,
	}
}

func goRight(c *coordinate) *coordinate {
	return &coordinate{
		x:       c.x + 1,
		y:       c.y,
		stepnum: c.stepnum + 1,
	}
}

func processInstruction(i *instruction, c *coordinate) (*coordinate, *coordinateSet) {
	var op func(*coordinate) *coordinate
	switch i.direction {
	case down:
		op = goDown
	case up:
		op = goUp
	case left:
		op = goLeft
	case right:
		op = goRight
	}
	set := make([]*coordinate, i.distance)
	cur := c
	for j := 0; j < i.distance; j++ {
		n := op(cur)
		set[j] = n
		cur = n
	}
	return cur, &coordinateSet{coords: set}
}
