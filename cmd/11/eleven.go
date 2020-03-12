package eleven

import (
	"advent-of-code-2019/logger"
	"advent-of-code-2019/utils"
	"fmt"

	"github.com/spf13/cobra"
)

type color int

const (
	black color = 0
	white color = 1
)

type coord struct {
	x int
	y int
}

var (
	left  = coord{x: -1, y: 0}
	right = coord{x: 1, y: 0}
	up    = coord{x: 0, y: 1}
	down  = coord{x: 0, y: -1}
)

type PaintingRobot struct {
	Direction     coord
	Position      coord
	Path          []coord
	VisitedCoords map[coord]color
}

func (r *PaintingRobot) Run(p *utils.Program, startingColor color) {
	r.Direction = up
	r.Position = coord{x: 0, y: 0}
	r.Path = make([]coord, 1)
	r.VisitedCoords = make(map[coord]color)
	r.Path[0] = r.Position
	r.VisitedCoords[r.Position] = startingColor
	inputChan := make(chan int)
	outputChan := make(chan int, 2)
	resultChan := make(chan utils.ExecutionResult)
	go p.RunAsync("robotprog", inputChan, outputChan, resultChan)

	for {
		var c color
		// check if we have been at this position before
		c, haveVisited := r.VisitedCoords[r.Position]
		// default is black
		if !haveVisited {
			c = black
		}
		select {
		case <-resultChan:
			return
		// tell the program our current color
		case inputChan <- int(c):
			// receive the new color
			rawNewColor := <-outputChan
			if rawNewColor == 0 {
				r.VisitedCoords[r.Position] = black
			} else {
				r.VisitedCoords[r.Position] = white
			}
			// receive the direction to turn
			turnDirection := <-outputChan
			// fmt.Printf("got new direction\n")
			switch r.Direction {
			case up:
				if turnDirection == 0 {
					r.Direction = left
				} else {
					r.Direction = right
				}
			case left:
				if turnDirection == 0 {
					r.Direction = down
				} else {
					r.Direction = up
				}
			case down:
				if turnDirection == 0 {
					r.Direction = right
				} else {
					r.Direction = left
				}
			case right:
				if turnDirection == 0 {
					r.Direction = up
				} else {
					r.Direction = down
				}
			}
			// move
			newPosition := coord{
				x: r.Position.x + r.Direction.x,
				y: r.Position.y + r.Direction.y,
			}
			r.Position = newPosition
			r.Path = append(r.Path, r.Position)
		}
	}
}

func renderPaint(visitedCoords map[coord]color) {
	// need to get max and min x/y for later use
	var maxY int
	var maxX int
	var minY int
	var minX int

	for crd, _ := range visitedCoords {
		// track max/minx
		if crd.x > maxX {
			maxX = crd.x
		} else if crd.x < minX {
			minX = crd.x
		}
		if crd.y > maxY {
			maxY = crd.y
		} else if crd.y < minY {
			minY = crd.y
		}
	}

	// combine quadrants
	xSize := (-1 * minX) + maxX + 1
	ySize := (-1 * minY) + maxY + 1
	zeroXIndex := -1 * minX
	zeroYIndex := -1 * minY
	fullQuad := make([][]color, ySize)

	for crd, colr := range visitedCoords {
		var xIndex int
		var yIndex int

		if crd.y >= 0 && crd.x >= 0 {
			yIndex = crd.y + zeroYIndex
			xIndex = crd.x + zeroXIndex
		} else if crd.y < 0 && crd.x >= 0 {
			yIndex = zeroYIndex - (crd.y * -1)
			xIndex = crd.x + zeroXIndex
		} else if crd.y < 0 && crd.x < 0 {
			yIndex = zeroYIndex - (crd.y * -1)
			xIndex = zeroXIndex - (crd.x * -1)
		} else {
			yIndex = crd.y + zeroYIndex
			xIndex = zeroXIndex - (crd.x * -1)
		}
		if len(fullQuad[yIndex]) == 0 {
			fullQuad[yIndex] = make([]color, xSize)
		}
		currentRow := fullQuad[yIndex]
		currentRow[xIndex] = colr
		fullQuad[yIndex] = currentRow
	}

	// render coordinate grid
	for i := len(fullQuad) - 1; i >= 0; i-- {
		for j := 0; j < len(fullQuad[i]); j++ {
			colr := fullQuad[i][j]
			if colr == black {
				fmt.Printf(" ")
			} else {
				fmt.Printf("\u001b[37m#\u001b[0m")
			}
		}
		fmt.Printf("\n")
	}
}

// Eleven does day 11
var Eleven = &cobra.Command{
	Use:   "eleven",
	Short: "day eleven",
	Run:   eleven,
}

var input = `3,8,1005,8,306,1106,0,11,0,0,0,104,1,104,0,3,8,1002,8,-1,10,1001,10,1,10,4,10,108,1,8,10,4,10,1002,8,1,28,2,107,3,10,1,101,19,10,3,8,1002,8,-1,10,1001,10,1,10,4,10,1008,8,0,10,4,10,102,1,8,59,2,5,13,10,3,8,102,-1,8,10,1001,10,1,10,4,10,1008,8,0,10,4,10,1001,8,0,85,3,8,1002,8,-1,10,101,1,10,10,4,10,1008,8,1,10,4,10,1001,8,0,107,1006,0,43,3,8,1002,8,-1,10,1001,10,1,10,4,10,1008,8,1,10,4,10,101,0,8,132,3,8,102,-1,8,10,1001,10,1,10,4,10,1008,8,0,10,4,10,1001,8,0,154,2,4,1,10,2,4,9,10,3,8,1002,8,-1,10,101,1,10,10,4,10,108,0,8,10,4,10,1001,8,0,183,1,1102,5,10,1,1102,1,10,1006,0,90,2,9,12,10,3,8,102,-1,8,10,1001,10,1,10,4,10,1008,8,0,10,4,10,1001,8,0,221,1006,0,76,1006,0,27,1,102,9,10,3,8,1002,8,-1,10,1001,10,1,10,4,10,108,1,8,10,4,10,102,1,8,252,2,4,9,10,1006,0,66,3,8,1002,8,-1,10,101,1,10,10,4,10,1008,8,1,10,4,10,101,0,8,282,1,102,19,10,101,1,9,9,1007,9,952,10,1005,10,15,99,109,628,104,0,104,1,21102,1,387240010644,1,21101,0,323,0,1105,1,427,21102,846541370112,1,1,21101,334,0,0,1106,0,427,3,10,104,0,104,1,3,10,104,0,104,0,3,10,104,0,104,1,3,10,104,0,104,1,3,10,104,0,104,0,3,10,104,0,104,1,21102,3425718295,1,1,21102,381,1,0,1105,1,427,21102,179410541715,1,1,21101,0,392,0,1106,0,427,3,10,104,0,104,0,3,10,104,0,104,0,21101,0,718078255872,1,21101,0,415,0,1105,1,427,21102,1,868494234468,1,21102,1,426,0,1105,1,427,99,109,2,21202,-1,1,1,21101,0,40,2,21101,458,0,3,21101,0,448,0,1106,0,491,109,-2,2106,0,0,0,1,0,0,1,109,2,3,10,204,-1,1001,453,454,469,4,0,1001,453,1,453,108,4,453,10,1006,10,485,1102,0,1,453,109,-2,2105,1,0,0,109,4,2102,1,-1,490,1207,-3,0,10,1006,10,508,21102,1,0,-3,22102,1,-3,1,22101,0,-2,2,21102,1,1,3,21102,1,527,0,1106,0,532,109,-4,2105,1,0,109,5,1207,-3,1,10,1006,10,555,2207,-4,-2,10,1006,10,555,22101,0,-4,-4,1105,1,623,22101,0,-4,1,21201,-3,-1,2,21202,-2,2,3,21101,574,0,0,1105,1,532,21202,1,1,-4,21102,1,1,-1,2207,-4,-2,10,1006,10,593,21102,0,1,-1,22202,-2,-1,-2,2107,0,-3,10,1006,10,615,21201,-1,0,1,21101,615,0,0,106,0,490,21202,-2,-1,-2,22201,-4,-2,-4,109,-5,2105,1,0`

func eleven(cmd *cobra.Command, args []string) {
	prog, err := utils.NewProgram(input)
	if err != nil {
		logger.Fatalf("failed to initialize program: %s", err.Error())
	}

	r := PaintingRobot{}
	r.Run(prog, black)
	logger.Infof("(part 1): painted %d panels", len(r.VisitedCoords))
	r.Run(prog, white)
	renderPaint(r.VisitedCoords)
}
