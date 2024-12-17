package main

import (
	"log/slog"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day17 = runner.DayImplementation{
	DayNumber:    17,
	ExecutePart1: Day17Part1,
	ExecutePart2: Day17Part2,
	ExampleInput: `Register A: 2024
Register B: 0
Register C: 0

Program: 0,3,5,4,3,0`,
	ExamplePart1Answer: "5,7,3,0",
	ExamplePart2Answer: "117440",
}

const (
	INS_ADV int = iota
	INS_BXL
	INS_BST
	INS_JNZ
	INS_BXC
	INS_OUT
	INS_BDV
	INS_CDV
)

type d17Program struct {
	data []int
	regA int
	regB int
	regC int
}

type programResult struct {
	output         string
	successfulCopy bool
}

func (prog *d17Program) Execute(successfulCopyRequired bool) programResult {
	insPtr, cpyPtr := 0, 0
	var output strings.Builder

	for insPtr < len(prog.data) {
		instruction := prog.data[insPtr]
		literal := prog.data[insPtr+1]
		combo := comboOperand(literal, prog.regA, prog.regB, prog.regC)

		switch instruction {
		case INS_ADV:
			divisor := 1
			for i := 0; i < combo; i++ {
				divisor *= 2
			}
			prog.regA /= divisor
		case INS_BXL:
			prog.regB ^= literal
		case INS_BST:
			prog.regB = combo % 8
		case INS_JNZ:
			if prog.regA != 0 {
				insPtr = literal - 2
			}
		case INS_BXC:
			prog.regB = prog.regB ^ prog.regC
		case INS_OUT:
			value := combo % 8
			if successfulCopyRequired {
				if cpyPtr >= len(prog.data) || value != prog.data[cpyPtr] {
					return programResult{"", false}
				} else {
					cpyPtr++
				}
			} else {
				output.WriteString(strconv.Itoa(value))
				output.WriteString(",")
			}
		case INS_BDV:
			divisor := 1
			for i := 0; i < combo; i++ {
				divisor *= 2
			}
			prog.regB = prog.regA / divisor
		case INS_CDV:
			divisor := 1
			for i := 0; i < combo; i++ {
				divisor *= 2
			}
			prog.regC = prog.regA / divisor
		}
		insPtr += 2
	}

	if successfulCopyRequired {
		if cpyPtr == len(prog.data) {
			return programResult{"", true}
		} else {
			return programResult{"", false}
		}
	} else {
		result := output.String()
		return programResult{result[:len(result)-1], false}
	}
}

func Day17Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Split(input, "\n")
	prog := d17Program{}
	prog.regA, _ = strconv.Atoi(lines[0][12:])
	prog.regB, _ = strconv.Atoi(lines[1][12:])
	prog.regC, _ = strconv.Atoi(lines[2][12:])
	dataStrs := strings.Split(lines[4][9:], ",")
	prog.data = make([]int, len(dataStrs))
	for ix, str := range dataStrs {
		prog.data[ix], _ = strconv.Atoi(str)
	}

	progCopy := prog
	return progCopy.Execute(false).output, prog
}

func comboOperand(operand int, regA int, regB int, regC int) int {
	switch {
	case operand >= 0 && operand <= 3:
		return operand
	case operand == 4:
		return regA
	case operand == 5:
		return regB
	case operand == 6:
		return regC
	default:
		panic("Unexpected operand")
	}
}

func Day17Part2(logger *slog.Logger, input string, part1Context any) string {
	prog := part1Context.(d17Program)

	regValue := 0
	for {
		progCopy := prog
		progCopy.regA = regValue
		if progCopy.Execute(true).successfulCopy {
			break
		}
		regValue++
	}

	return strconv.Itoa(regValue)
}
