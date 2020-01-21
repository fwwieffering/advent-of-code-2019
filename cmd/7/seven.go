package seven

import (
	"advent-of-code-2019/logger"
	"advent-of-code-2019/utils"
	"fmt"

	"github.com/spf13/cobra"
)

// Seven day 7
var Seven = &cobra.Command{
	Use:   "seven",
	Short: "does day 7",
	Run:   seven,
}

var input = `3,8,1001,8,10,8,105,1,0,0,21,34,43,64,85,98,179,260,341,422,99999,3,9,1001,9,3,9,102,3,9,9,4,9,99,3,9,102,5,9,9,4,9,99,3,9,1001,9,2,9,1002,9,4,9,1001,9,3,9,1002,9,4,9,4,9,99,3,9,1001,9,3,9,102,3,9,9,101,4,9,9,102,3,9,9,4,9,99,3,9,101,2,9,9,1002,9,3,9,4,9,99,3,9,101,1,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,102,2,9,9,4,9,3,9,102,2,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1001,9,1,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,2,9,4,9,99,3,9,101,1,9,9,4,9,3,9,102,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,101,1,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1002,9,2,9,4,9,99,3,9,101,1,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,1001,9,1,9,4,9,3,9,101,1,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,1,9,9,4,9,3,9,101,1,9,9,4,9,99,3,9,101,1,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,1002,9,2,9,4,9,99,3,9,101,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,1,9,4,9,3,9,102,2,9,9,4,9,99`

func seven(cmd *cobra.Command, args []string) {
	prog, err := utils.NewProgram(input)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	part1Max, part1Seq, err := testAmplifiers(prog, 5)
	if err != nil {
		logger.Fatalf(err.Error())
	}
	logger.Infof("(part 1): Max: %d. Sequence: %+v", part1Max, part1Seq)
	part2Max, part2Seq, err := testAmplifiersFeedback(prog, 5)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	logger.Infof("(part 2): Max %d. Sequence: %+v", part2Max, part2Seq)
}

func runSequence(p *utils.Program, sequence []int) (int, error) {
	// first inputsignal is always 0
	var inputsignal = 0
	for _, phaseSetting := range sequence {
		_, err := p.Run(phaseSetting, inputsignal)
		if err != nil {
			return 0, err
		}
		// chain inputsignal for next amp
		inputsignal = p.Result
	}
	// output of last program is the answer
	return inputsignal, nil
}

func runSequenceFeedback(p *utils.Program, sequence []int) (int, error) {
	type progExec struct {
		program    *utils.Program
		inputChan  chan int
		outputChan chan int
	}
	resultChan := make(chan utils.ExecutionResult, len(sequence))
	runList := make([]progExec, len(sequence))

	// make channels
	chans := make([]chan int, len(sequence))
	for i, phaseSetting := range sequence {
		chans[i] = make(chan int, 2)
		chans[i] <- phaseSetting
		if i == 0 {
			chans[i] <- 0
		}
	}

	for i := range sequence {
		var ipt chan int
		var opt chan int

		ipt = chans[i]

		if i == (len(sequence) - 1) {
			opt = chans[0]
		} else {
			opt = chans[i+1]
		}

		runList[i] = progExec{
			program:    p,
			inputChan:  ipt,
			outputChan: opt,
		}
		// fmt.Printf("program %d: input %+v output %+v\n", i, ipt, opt)
		identifier := fmt.Sprintf("amplifier %d", i)
		go p.RunAsync(identifier, ipt, opt, resultChan)
	}

	errstr := ""
	var answer int
	for i := 0; i < len(runList); i++ {
		res := <-resultChan
		// fmt.Printf("result: %+v\n", res)
		if res.Error != nil {
			errstr += res.Error.Error() + "\n"
		}
		answer = res.Value
	}
	var err error
	if len(errstr) > 0 {
		err = fmt.Errorf(errstr)
	}
	return answer, err
}

func testAmplifiers(program *utils.Program, numAmplifiers int) (int, []int, error) {
	// phase setting must be between 0, 4
	phaseCombinations := generateCombinations([]int{0, 1, 2, 3, 4}, 5)

	maxResult := 0
	var bestSequence []int
	for _, phaseCombo := range phaseCombinations {
		res, err := runSequence(program, phaseCombo)
		if err != nil {
			return 0, nil, err
		}
		if res > maxResult {
			maxResult = res
			bestSequence = phaseCombo
		}
	}
	return maxResult, bestSequence, nil
}

func testAmplifiersFeedback(program *utils.Program, numAmplifiers int) (int, []int, error) {
	phaseCombinations := generateCombinations([]int{5, 6, 7, 8, 9}, 5)

	maxResult := 0
	var bestSequence []int
	for _, phaseCombo := range phaseCombinations {
		res, err := runSequenceFeedback(program, phaseCombo)
		if err != nil {
			return 0, nil, err
		}
		if res > maxResult {
			maxResult = res
			bestSequence = phaseCombo
		}
	}
	return maxResult, bestSequence, nil
}

// r = array size
func generateCombinations(alphabet []int, r int) [][]int {
	data := make([]int, r)
	returnValue := make([][]int, 0)
	combinationRecurse(alphabet, data, 0, r, &returnValue)
	return returnValue
}

func deleteItem(arr []int, index int) []int {
	copyArr := make([]int, len(arr))
	copy(copyArr, arr)

	if len(arr) <= 1 {
		return []int{}
	}
	if index == 0 {
		return copyArr[1:]
	}
	if index == len(arr)-1 {
		return copyArr[:index]
	}
	copyArr = append(copyArr[:index], copyArr[index+1:]...)

	return copyArr
}

func combinationRecurse(alphabet []int, data []int, idx int, r int, returnValue *[][]int) {
	if idx == r {
		newdata := make([]int, len(data))
		copy(newdata, data)
		*returnValue = append(*returnValue, newdata)
		return
	}
	for i := range alphabet {
		data[idx] = alphabet[i]
		// exclude the item that was just used from the new alphabet
		unusedItems := deleteItem(alphabet, i)
		combinationRecurse(unusedItems, data, idx+1, r, returnValue)
	}
}
