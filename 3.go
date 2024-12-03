package main

import (
	"log/slog"
	"strconv"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day3 = runner.DayImplementation{
	DayNumber:          3,
	ExecutePart1:       Day3Part1,
	ExecutePart2:       Day3Part2,
	ExampleInput:       "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))",
	ExamplePart1Answer: "161",
	ExamplePart2Answer: "48",
}

type Action int

const (
	ACTION_RESET Action = iota
	ACTION_ACCEPTED
	ACTION_OPERAND_1
	ACTION_OPERAND_2
	ACTION_DISABLE
	ACTION_ENABLE
	ACTION_COMPLETED
)

type State int

const (
	STATE_INITIAL State = iota
	STATE_M
	STATE_U
	STATE_L
	STATE_D
	STATE_O
	STATE_N
	STATE_APOSTROPHE
	STATE_T
	STATE_DO
	STATE_DONT
	STATE_OPERAND_1
	STATE_OPERAND_2
)

func handleCharacter(char rune, state State, disabled bool) (Action, State) {
	if disabled && state == STATE_INITIAL && char != 'd' {
		return ACTION_RESET, STATE_INITIAL
	}

	switch {
	case char == 'm':
		return ACTION_RESET, STATE_M
	case char == 'u':
		if state == STATE_M {
			return ACTION_ACCEPTED, STATE_U
		}
	case char == 'l':
		if state == STATE_U {
			return ACTION_ACCEPTED, STATE_L
		}
	case char == 'd':
		return ACTION_RESET, STATE_D
	case char == 'o':
		if state == STATE_D {
			return ACTION_ACCEPTED, STATE_O
		}
	case char == 'n':
		if state == STATE_O {
			return ACTION_ACCEPTED, STATE_N
		}
	case char == '\'':
		if state == STATE_N {
			return ACTION_ACCEPTED, STATE_APOSTROPHE
		}
	case char == 't':
		if state == STATE_APOSTROPHE {
			return ACTION_ACCEPTED, STATE_T
		}
	case char == '(':
		switch state {
		case STATE_L:
			return ACTION_ACCEPTED, STATE_OPERAND_1
		case STATE_O:
			return ACTION_ACCEPTED, STATE_DO
		case STATE_T:
			return ACTION_ACCEPTED, STATE_DONT
		}
	case char == ',':
		if state == STATE_OPERAND_1 {
			return ACTION_ACCEPTED, STATE_OPERAND_2
		}
	case char == ')':
		switch state {
		case STATE_OPERAND_2:
			return ACTION_COMPLETED, STATE_INITIAL
		case STATE_DO:
			return ACTION_ENABLE, STATE_INITIAL
		case STATE_DONT:
			return ACTION_DISABLE, STATE_INITIAL
		}
	case char >= '0' && char <= '9':
		if state == STATE_OPERAND_1 {
			return ACTION_OPERAND_1, STATE_OPERAND_1
		} else if state == STATE_OPERAND_2 {
			return ACTION_OPERAND_2, STATE_OPERAND_2
		}
	}
	return ACTION_RESET, STATE_INITIAL
}

func Day3Part1(logger *slog.Logger, input string) (string, any) {
	var operand1, operand2, sum int
	state := STATE_INITIAL
	var action Action
	for _, char := range input {
		action, state = handleCharacter(char, state, false)
		switch action {
		case ACTION_OPERAND_1:
			operand1 *= 10
			operand1 += int(char - '0')
		case ACTION_OPERAND_2:
			operand2 *= 10
			operand2 += int(char - '0')
		case ACTION_COMPLETED:
			sum += (operand1 * operand2)
			fallthrough
		case ACTION_RESET:
			operand1 = 0
			operand2 = 0
		}
	}

	return strconv.Itoa(sum), nil
}

func Day3Part2(logger *slog.Logger, input string, part1Context any) string {
	var operand1, operand2, sum int
	state := STATE_INITIAL
	disabled := false
	var action Action
	for _, char := range input {
		action, state = handleCharacter(char, state, disabled)
		switch action {
		case ACTION_OPERAND_1:
			operand1 *= 10
			operand1 += int(char - '0')
		case ACTION_OPERAND_2:
			operand2 *= 10
			operand2 += int(char - '0')
		case ACTION_COMPLETED:
			sum += (operand1 * operand2)
			fallthrough
		case ACTION_RESET:
			operand1 = 0
			operand2 = 0
		case ACTION_DISABLE:
			disabled = true
		case ACTION_ENABLE:
			disabled = false
		}
	}

	return strconv.Itoa(sum)
}
