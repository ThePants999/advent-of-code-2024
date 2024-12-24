package main

import (
	"log/slog"
	"slices"
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

func (prog *d17Program) Execute() []int {
	insPtr := 0
	output := make([]int, 0, len(prog.data))

	for insPtr < len(prog.data) {
		instruction := prog.data[insPtr]
		literal := prog.data[insPtr+1]
		combo := comboOperand(literal, prog.regA, prog.regB, prog.regC)

		switch instruction {
		case INS_ADV:
			prog.regA >>= combo
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
			output = append(output, combo%8)
		case INS_BDV:
			prog.regB = prog.regA >> combo
		case INS_CDV:
			prog.regC = prog.regA >> combo
		}
		insPtr += 2
	}

	return output
}

func Day17Part1(logger *slog.Logger, input string) (string, any) {
	// Parse the input.
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

	// Part 1 is simple enough - genuinely run the
	// program.
	progCopy := prog
	output := progCopy.Execute()
	var result strings.Builder
	for _, val := range output {
		result.WriteString(strconv.Itoa(val))
		result.WriteString(",")
	}
	resStr := result.String()
	return resStr[:len(resStr)-1], prog
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

	candidateA := 0
	// It is approximately the case that each 3 bits of register A will
	// determine one output value, with the least significant bits
	// corresponding to the first output value. So what we're going to do
	// looks like this:
	// -  Try a set of small A values until we find one that outputs the
	//    LAST output value.
	// -  When we find it, "lock it in" by shifting 3 bits left. We don't
	//    technically lock it in - we can change those bits further - but
	//    they're our starting point for what follows.
	// -  We then start trying candidates incrementally to find one that
	//    outputs the last TWO output values. It's possible for bits
	//    more significant than the last three to influence the last
	//    output value, so the bits we've "locked in" might be wrong, and
	//    that's why we need to check the full set of output values so far
	//    on every iteration. But we shouldn't need to change too much of
	//    what we've "locked in" before we find the right combo.
	// -  Repeat until we've got the whole set of output values.
	for pos := len(prog.data) - 1; pos >= 0; pos-- {
		candidateA <<= 3
		for {
			prog.regA = candidateA
			output := prog.Execute()
			if slices.Equal(output, prog.data[pos:]) {
				break
			}
			candidateA++
		}
	}

	return strconv.Itoa(candidateA)
}
