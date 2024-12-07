package main

import (
	"log/slog"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
	iterium "github.com/mowshon/iterium"
)

var Day7 = runner.DayImplementation{
	DayNumber:    7,
	ExecutePart1: Day7Part1,
	ExecutePart2: Day7Part2,
	ExampleInput: `190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20`,
	ExamplePart1Answer: "3749",
	ExamplePart2Answer: "11387",
}

type equation struct {
	result   int
	operands []int
}

type operator int

const (
	OP_ADD operator = iota
	OP_MULTIPLY
	OP_CONCAT
)

func Day7Part1(logger *slog.Logger, input string) (string, any) {
	numbers := strings.Fields(input)
	equations := make([]*equation, 0, 1000)
	var currEq *equation
	for _, number := range numbers {
		if number[len(number)-1] == ':' {
			result, err := strconv.Atoi(number[:len(number)-1])
			if err != nil {
				panic("Invalid input")
			}
			currEq = &equation{result, make([]int, 0, 20)}
			equations = append(equations, currEq)
		} else {
			num, err := strconv.Atoi(number)
			if err != nil {
				panic("Invalid input")
			}
			currEq.operands = append(currEq.operands, num)
		}
	}

	sum := runTest(equations, []operator{OP_ADD, OP_MULTIPLY})
	return strconv.Itoa(sum), equations
}

func Day7Part2(logger *slog.Logger, input string, part1Context any) string {
	equations := part1Context.([]*equation)
	sum := runTest(equations, []operator{OP_ADD, OP_MULTIPLY, OP_CONCAT})
	return strconv.Itoa(sum)
}

func runTest(equations []*equation, operators []operator) int {
	c := make(chan int)
	for _, eq := range equations {
		go testEquation(eq, operators, c)
	}

	sum := 0
	for i := 0; i < len(equations); i++ {
		sum += <-c
	}

	return sum
}

func testEquation(eq *equation, operators []operator, c chan int) {
	combinations := iterium.Product(operators, len(eq.operands)-1)
	for {
		combination, err := combinations.Next()
		if err != nil {
			break
		}

		value := eq.operands[0]
		for i := 0; i < len(combination); i++ {
			switch combination[i] {
			case OP_ADD:
				value += eq.operands[i+1]
			case OP_MULTIPLY:
				value *= eq.operands[i+1]
			case OP_CONCAT:
				value, err = strconv.Atoi(strconv.Itoa(value) + strconv.Itoa(eq.operands[i+1]))
				if err != nil {
					panic("Possible overflow?")
				}
			}
			if value > eq.result {
				break
			}
		}
		if value == eq.result {
			c <- eq.result
			return
		}
	}

	c <- 0
}
