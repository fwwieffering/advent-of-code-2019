package two

import (
	"fmt"
	"strconv"
	"strings"
)

type Program struct {
	InitialMemory []int
}

func (p *Program) Run(input1 int, input2 int) ([]int, error) {
	memoryCopy := make([]int, len(p.InitialMemory))
	copy(memoryCopy, p.InitialMemory)

	memoryCopy[1] = input1
	memoryCopy[2] = input2

	return RunProgram(memoryCopy)
}

var commandCodeMap = map[int]func(intcode []int, inputpositions []int, outputposition int){
	1:  Add,
	2:  Multiply,
	99: End,
}

// Command holds the information to run an intcode command
type Command struct {
	Operation      func(intcode []int, inputpositions []int, outputposition int)
	InputPositions []int
	OutputPosition int
}

// Add adds the inputs at inputpositions and updates the outputPosition
func Add(intcode []int, inputpositions []int, outputposition int) {
	result := 0
	for _, inputIndex := range inputpositions {
		result += intcode[inputIndex]
	}
	intcode[outputposition] = result
}

// End does nothing and terminates the program
func End(intcode []int, inputpositions []int, outputposition int) {
}

// Multiply multiplies the inputs at inputpositions and updates the outputPosition
func Multiply(intcode []int, inputpositions []int, outputposition int) {
	result := 0
	for idx, inputIndex := range inputpositions {
		if idx == 0 {
			result = intcode[inputIndex]
		} else {
			result = result * intcode[inputIndex]
		}
	}
	intcode[outputposition] = result
}

// Run executes the command on the passed intcode
func (c *Command) Run(intcode []int) {
	c.Operation(intcode, c.InputPositions, c.OutputPosition)
}

func convertStringIntCode(intcodeString string) ([]int, error) {
	splitIntCode := strings.Split(intcodeString, ",")
	parsedIntCode := make([]int, len(splitIntCode))
	for i := range splitIntCode {
		currentInt, err := strconv.Atoi(splitIntCode[i])
		if err != nil {
			return nil, fmt.Errorf("couldn't convert %s at index %d into integer. %s", splitIntCode[i], i, err.Error())
		}
		parsedIntCode[i] = currentInt
	}
	return parsedIntCode, nil
}

func convertIntCodeString(intCode []int) string {
	s := ""
	for _, i := range intCode {
		s += fmt.Sprintf("%d,", i)
	}
	return strings.TrimRight(s, ",")
}

func RunProgram(intCode []int) ([]int, error) {
	i := 0
	for i < len(intCode) {
		currentInt := intCode[i]
		op, ok := commandCodeMap[currentInt]
		if !ok {
			return nil, fmt.Errorf("Expected command code at index %d but got: %d which does not map to a command code", i, currentInt)
		}
		// 99 ends the program. Break
		if currentInt == 99 {
			break
		}
		command := Command{
			Operation:      op,
			InputPositions: []int{intCode[i+1], intCode[i+2]},
			OutputPosition: intCode[i+3],
		}
		command.Run(intCode)
		i += 4
	}
	return intCode, nil
}

// RunProgramString takes a comma separated intcode string, parses, and runs it
func RunProgramString(intcodestring string) (string, error) {
	intCode, err := convertStringIntCode(intcodestring)
	if err != nil {
		return "", err
	}
	res, err := RunProgram(intCode)
	if err != nil {
		return "", err
	}
	return convertIntCodeString(res), nil
}
