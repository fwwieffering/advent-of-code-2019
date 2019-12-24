package utils

import (
	"fmt"
)

// Program holds a program
type Program struct {
	Input         []int
	Result        int
	InitialMemory []int
}

type parameterMode int

const (
	position  parameterMode = 0
	immediate parameterMode = 1
	invalid   parameterMode = -1
)

func parseParameterMode(i int) parameterMode {
	switch i {
	case 0:
		return position
	case 1:
		return immediate
	default:
		return invalid
	}
}

type parameter struct {
	mode     parameterMode
	position int
}

func (prm *parameter) getValue(intcode []int) int {
	var i int
	switch prm.mode {
	case position:
		idx := intcode[prm.position]
		i = intcode[idx]
	case immediate:
		i = intcode[prm.position]
	default:
		return -1
	}
	return i
}

type VMCommand struct {
	Name         string
	CommandCode  int
	Operation    func(program *Program, ptr int, memory []int, parameters []*parameter) int
	NumberParams int
}

var (
	AddCMD = &VMCommand{
		Name: "Add",
		Operation: func(p *Program, ptr int, memory []int, params []*parameter) int {
			res := 0
			for _, p := range params {
				res += p.getValue(memory)
			}
			// output position is after param
			outputPosition := memory[params[len(params)-1].position+1]
			memory[outputPosition] = res
			return ptr + 4
		},
		CommandCode:  1,
		NumberParams: 2,
	}
	MultiplyCMD = &VMCommand{
		Name: "Multiply",
		Operation: func(p *Program, ptr int, memory []int, params []*parameter) int {
			result := 0
			for idx, p := range params {
				if idx == 0 {
					result = p.getValue(memory)
				} else {
					result = result * p.getValue(memory)
				}
			}
			// output position is after param
			outputPosition := memory[params[len(params)-1].position+1]
			memory[outputPosition] = result
			return ptr + 4
		},
		CommandCode:  2,
		NumberParams: 2,
	}
	ReplaceCMD = &VMCommand{
		Name: "Replace",
		Operation: func(p *Program, ptr int, memory []int, params []*parameter) int {
			fmt.Printf("input command at %d. Input value %d\n", ptr, p.Input[0])
			input := p.Input[0]
			position := memory[params[0].position]
			memory[position] = input
			return ptr + 2
		},
		CommandCode:  3,
		NumberParams: 1,
	}
	OutputCMD = &VMCommand{
		Name: "Output",
		Operation: func(p *Program, ptr int, memory []int, params []*parameter) int {
			val := params[0].getValue(memory)
			fmt.Printf("output cmd at %d. Value: %d\n", params[0].position, val)
			p.Result = val
			return ptr + 2
		},
		CommandCode:  4,
		NumberParams: 1,
	}
	JumpIfTrueCMD = &VMCommand{
		Name: "JumpIfTrue",
		Operation: func(p *Program, ptr int, memory []int, params []*parameter) int {
			isTrue := params[0].getValue(memory) != 0
			if isTrue {
				return params[1].getValue(memory)
			}
			return ptr + 3
		},
		CommandCode:  5,
		NumberParams: 2,
	}
	JumpIfFalseCMD = &VMCommand{
		Name: "JumpIfFalse",
		Operation: func(p *Program, ptr int, memory []int, params []*parameter) int {
			isFalse := params[0].getValue(memory) == 0
			if isFalse {
				return params[1].getValue(memory)
			}
			return ptr + 3
		},
		CommandCode:  6,
		NumberParams: 2,
	}
	LessThan = &VMCommand{
		Name: "LessThan",
		Operation: func(p *Program, ptr int, memory []int, params []*parameter) int {
			firstVal := params[0].getValue(memory)
			secondVal := params[1].getValue(memory)
			if firstVal < secondVal {
				memory[memory[params[2].position]] = 1
			} else {
				memory[memory[params[2].position]] = 0
			}
			return ptr + 4
		},
		CommandCode:  7,
		NumberParams: 3,
	}
	Equals = &VMCommand{
		Name: "Equals",
		Operation: func(p *Program, ptr int, memory []int, params []*parameter) int {
			firstVal := params[0].getValue(memory)
			secondVal := params[1].getValue(memory)
			if firstVal == secondVal {
				memory[memory[params[2].position]] = 1
			} else {
				memory[memory[params[2].position]] = 0
			}
			return ptr + 4
		},
		CommandCode:  8,
		NumberParams: 3,
	}
	EndCMD = &VMCommand{
		Name:         "End",
		CommandCode:  99,
		NumberParams: 0,
	}
)

var vmcommandCodeMap = map[int]*VMCommand{
	1:  AddCMD,
	2:  MultiplyCMD,
	3:  ReplaceCMD,
	4:  OutputCMD,
	5:  JumpIfTrueCMD,
	6:  JumpIfFalseCMD,
	7:  LessThan,
	8:  Equals,
	99: EndCMD,
}

func commandParse(memory []int, position int) (*VMCommand, []*parameter, error) {
	fullCmd := memory[position]
	// command codes are the rightmost 2 digits
	cmdCode := fullCmd % 100
	vmCmd, ok := vmcommandCodeMap[cmdCode]
	if !ok {
		return nil, nil, fmt.Errorf("No command found for %d at position %d", fullCmd, position)
	}
	params := make([]*parameter, vmCmd.NumberParams)
	// paramModes should be of magnitude 10 * len(vmCmd.NumParams)
	paramModeFull := fullCmd / 100
	// get parameter mode by digit by %10
	for i := 0; i < vmCmd.NumberParams; i++ {
		modeInt := paramModeFull % 10
		mode := parseParameterMode(modeInt)
		params[i] = &parameter{mode: mode, position: position + i + 1}
		paramModeFull = paramModeFull / 10
	}
	return vmCmd, params, nil
}

func (p *Program) Run(input ...int) ([]int, error) {
	memCopy := CopyIntArray(p.InitialMemory)
	i := 0
	p.Input = input
	for i < len(memCopy) {
		cmd, params, err := commandParse(memCopy, i)
		if err != nil {
			return nil, err
		}

		if cmd.CommandCode == 99 {
			break
		}
		newPtr := cmd.Operation(p, i, memCopy, params)
		i = newPtr
	}
	return memCopy, nil
}

func NewProgram(initialMemory string) (*Program, error) {
	intMem, err := StringToIntArray(initialMemory)
	p := &Program{InitialMemory: intMem}
	return p, err
}
