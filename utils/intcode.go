package utils

import (
	"fmt"
	"log"
	"os"
)

// Program holds a program
type Program struct {
	InitialMemory []int
	Executions    []*Execution
}

type Execution struct {
	Input        []int
	InputChan    chan int
	OutputChan   chan int
	Output       []int
	relativeBase int
	inputCounter int
	Memory       *Memory
}

// ExecutionResult is the result of an execution
type ExecutionResult struct {
	Memory []int
	Result []int
	Error  error
}

// Memory holds programs memory
type Memory struct {
	mem    []int
	length int
}

func NewMemory(initial []int) *Memory {
	return &Memory{
		mem:    CopyIntArray(initial),
		length: len(initial),
	}
}

// Len since memory can be accessed outside of the program, store original length
// for program running purposes
func (m *Memory) Len() int {
	return m.length
}

func (m *Memory) Get(position int) int {
	// extend memory as needed, by 2x needed amount
	if position > len(m.mem) {
		m.mem = append(m.mem, make([]int, ((position-len(m.mem))*2))...)
	}
	return m.mem[position]
}

func (m *Memory) Set(position int, val int) {
	if position > (len(m.mem) - 1) {
		m.mem = append(m.mem, make([]int, ((position-len(m.mem)+1)*2))...)
	}
	m.mem[position] = val
}

type parameterMode int

const (
	position  parameterMode = 0
	immediate parameterMode = 1
	relative  parameterMode = 2
	invalid   parameterMode = -1
)

func parseParameterMode(i int) parameterMode {
	switch i {
	case 0:
		return position
	case 1:
		return immediate
	case 2:
		return relative
	default:
		return invalid
	}
}

type parameter struct {
	mode     parameterMode
	position int
}

func (prm *parameter) getValue(mem *Memory, relativeBase int) int {
	var i int
	switch prm.mode {
	case position:
		idx := mem.Get(prm.position)
		i = mem.Get(idx)
	case immediate:
		i = mem.Get(prm.position)
	case relative:
		p := mem.Get(prm.position)
		i = mem.Get(p + relativeBase)
	default:
		return -1
	}
	return i
}

func (prm *parameter) getPosition(mem *Memory, relativeBase int) int {
	var i int
	switch prm.mode {
	case position, immediate:
		i = mem.Get(prm.position)
	case relative:
		p := mem.Get(prm.position)
		i = p + relativeBase
	}

	return i
}

type VMCommand struct {
	Name         string
	CommandCode  int
	Operation    func(e *Execution, ptr int, parameters []*parameter, logger *log.Logger) (int, error)
	NumberParams int
}

var (
	AddCMD = &VMCommand{
		Name: "Add",
		Operation: func(e *Execution, ptr int, params []*parameter, logger *log.Logger) (int, error) {
			res := 0
			for i := 0; i < 2; i++ {
				res += params[i].getValue(e.Memory, e.relativeBase)
			}
			// third param is output position
			outputPosition := params[2].getPosition(e.Memory, e.relativeBase)
			e.Memory.Set(outputPosition, res)
			return ptr + 4, nil
		},
		CommandCode:  1,
		NumberParams: 3,
	}
	MultiplyCMD = &VMCommand{
		Name: "Multiply",
		Operation: func(e *Execution, ptr int, params []*parameter, logger *log.Logger) (int, error) {
			result := 0
			for i := 0; i < 2; i++ {
				if i == 0 {
					result = params[i].getValue(e.Memory, e.relativeBase)
				} else {
					result = result * params[i].getValue(e.Memory, e.relativeBase)
				}
			}
			// third param is output position
			outputPosition := params[2].getPosition(e.Memory, e.relativeBase)
			e.Memory.Set(outputPosition, result)
			return ptr + 4, nil
		},
		CommandCode:  2,
		NumberParams: 3,
	}
	InputCMD = &VMCommand{
		Name: "Replace",
		Operation: func(e *Execution, ptr int, params []*parameter, logger *log.Logger) (int, error) {
			var val int
			// input could be from array or channel
			if e.InputChan != nil {
				// this can block
				val = <-e.InputChan
			} else {
				// catch errors
				if e.inputCounter > (len(e.Input) - 1) {
					return 0, fmt.Errorf("not enough input for input instruction. on %d input command but only %d inputs provided",
						e.inputCounter+1, len(e.Input))
				}
				val = e.Input[e.inputCounter]
			}
			e.inputCounter++
			position := params[0].getPosition(e.Memory, e.relativeBase)
			e.Memory.Set(position, val)
			// if logger != nil {
			// 	logger.Printf("input command at %d. Input value %d. Input #%d\n", ptr, val, e.InputCounter)
			// }

			return ptr + 2, nil
		},
		CommandCode:  3,
		NumberParams: 1,
	}
	OutputCMD = &VMCommand{
		Name: "Output",
		Operation: func(e *Execution, ptr int, params []*parameter, logger *log.Logger) (int, error) {
			val := params[0].getValue(e.Memory, e.relativeBase)
			if e.OutputChan != nil {
				e.OutputChan <- val
			}
			e.Output = append(e.Output, val)
			// if logger != nil {
			// 	logger.Printf("output cmd at %d. Value: %d\n", params[0].position, val)
			// }
			return ptr + 2, nil
		},
		CommandCode:  4,
		NumberParams: 1,
	}
	JumpIfTrueCMD = &VMCommand{
		Name: "JumpIfTrue",
		Operation: func(e *Execution, ptr int, params []*parameter, logger *log.Logger) (int, error) {
			isTrue := params[0].getValue(e.Memory, e.relativeBase) != 0
			if isTrue {
				return params[1].getValue(e.Memory, e.relativeBase), nil
			}
			return ptr + 3, nil
		},
		CommandCode:  5,
		NumberParams: 2,
	}
	JumpIfFalseCMD = &VMCommand{
		Name: "JumpIfFalse",
		Operation: func(e *Execution, ptr int, params []*parameter, logger *log.Logger) (int, error) {
			isFalse := params[0].getValue(e.Memory, e.relativeBase) == 0
			if isFalse {
				return params[1].getValue(e.Memory, e.relativeBase), nil
			}
			return ptr + 3, nil
		},
		CommandCode:  6,
		NumberParams: 2,
	}
	LessThan = &VMCommand{
		Name: "LessThan",
		Operation: func(e *Execution, ptr int, params []*parameter, logger *log.Logger) (int, error) {
			firstVal := params[0].getValue(e.Memory, e.relativeBase)
			secondVal := params[1].getValue(e.Memory, e.relativeBase)
			if firstVal < secondVal {
				e.Memory.Set(params[2].getPosition(e.Memory, e.relativeBase), 1)
			} else {
				e.Memory.Set(params[2].getPosition(e.Memory, e.relativeBase), 0)
			}
			return ptr + 4, nil
		},
		CommandCode:  7,
		NumberParams: 3,
	}
	Equals = &VMCommand{
		Name: "Equals",
		Operation: func(e *Execution, ptr int, params []*parameter, logger *log.Logger) (int, error) {
			firstVal := params[0].getValue(e.Memory, e.relativeBase)
			secondVal := params[1].getValue(e.Memory, e.relativeBase)
			if firstVal == secondVal {
				e.Memory.Set(params[2].getPosition(e.Memory, e.relativeBase), 1)
			} else {
				e.Memory.Set(params[2].getPosition(e.Memory, e.relativeBase), 0)
			}
			return ptr + 4, nil
		},
		CommandCode:  8,
		NumberParams: 3,
	}
	RelativeBaseAdjust = &VMCommand{
		Name: "RelativeBaseAdjust",
		Operation: func(e *Execution, ptr int, params []*parameter, logger *log.Logger) (int, error) {
			e.relativeBase += params[0].getValue(e.Memory, e.relativeBase)
			return ptr + 2, nil
		},
		CommandCode:  9,
		NumberParams: 1,
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
	3:  InputCMD,
	4:  OutputCMD,
	5:  JumpIfTrueCMD,
	6:  JumpIfFalseCMD,
	7:  LessThan,
	8:  Equals,
	9:  RelativeBaseAdjust,
	99: EndCMD,
}

func commandParse(memory *Memory, position int) (*VMCommand, []*parameter, error) {
	fullCmd := memory.Get(position)
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

func (p *Program) Run(input ...int) (*ExecutionResult, error) {
	memCopy := NewMemory(p.InitialMemory)
	e := &Execution{
		Memory:       memCopy,
		Output:       make([]int, 0),
		Input:        input,
		inputCounter: 0,
	}
	i := 0
	for i < memCopy.Len() {
		cmd, params, err := commandParse(memCopy, i)
		if err != nil {
			return nil, err
		}
		// fmt.Printf("i: %d. code: %d relative base: %d command: %s %+v\n", i, memCopy.Get(i), e.relativeBase, cmd.Name, memCopy.mem[i:i+len(params)+1])
		if cmd.CommandCode == 99 {
			break
		}
		newPtr, err := cmd.Operation(e, i, params, nil)
		if err != nil {
			return nil, fmt.Errorf("error at %d: %s", i, err.Error())
		}
		i = newPtr
	}
	return &ExecutionResult{Memory: memCopy.mem, Result: e.Output}, nil
}

func (p Program) RunAsync(identifier string, inputchan chan int, outputchan chan int, resultChan chan ExecutionResult) {

	var lgr *log.Logger
	lgr = log.New(os.Stdout, fmt.Sprintf("%s :", identifier), 0)

	memCopy := NewMemory(p.InitialMemory)

	e := &Execution{
		Memory:     memCopy,
		InputChan:  inputchan,
		OutputChan: outputchan,
		Output:     make([]int, 0),
	}
	i := 0
	for i < memCopy.Len() {
		cmd, params, err := commandParse(memCopy, i)
		if err != nil {
			resultChan <- ExecutionResult{Error: err}
		}

		if cmd.CommandCode == 99 {
			break
		}
		newPtr, err := cmd.Operation(e, i, params, lgr)
		if err != nil {
			resultChan <- ExecutionResult{Error: fmt.Errorf("error at %d: %s", i, err.Error())}
		}
		i = newPtr
	}
	res := ExecutionResult{
		Memory: memCopy.mem,
		Result: e.Output,
	}

	resultChan <- res
}

func NewProgram(initialMemory string) (*Program, error) {
	intMem, err := StringToIntArray(initialMemory)
	p := &Program{InitialMemory: intMem}
	return p, err
}
