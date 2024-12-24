package main

import (
	"log/slog"
	"math"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
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

func Day7Part1(logger *slog.Logger, input string) (string, any) {
	// Parse the input into a slice of equation structs,
	// each recording both the desired result and the
	// operands we've been given.
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

	sum := runTest(equations, false)
	return strconv.Itoa(sum), equations
}

func Day7Part2(logger *slog.Logger, input string, part1Context any) string {
	equations := part1Context.([]*equation)
	sum := runTest(equations, true)
	return strconv.Itoa(sum)
}

func runTest(equations []*equation, allowConcatenation bool) int {
	// Each equation is completely independent, so farm them
	// out to a separate goroutine each for parallel processing.
	c := make(chan int)
	for _, eq := range equations {
		go func() {
			if testEquation(eq, eq.operands[0], 1, allowConcatenation) {
				c <- eq.result
			} else {
				c <- 0
			}
		}()
	}

	// Collate the results.
	sum := 0
	for i := 0; i < len(equations); i++ {
		sum += <-c
	}

	return sum
}

// Recursive function attempting to solve _equation_.
// Takes the cumulative value calculated so far, and
// tries to figure out the operator before the _index_th
// operand.
func testEquation(eq *equation, value int, index int, allowConcatenation bool) bool {
	if index == len(eq.operands) {
		// We're done, so whether we're successful
		// depends on whether the cumulative result
		// so far equals the final result.
		return eq.result == value
	}

	if value > eq.result {
		// If we go past the desired result, bug out
		// early.
		return false
	}

	// Fork to try +...
	newValue := value + eq.operands[index]
	if testEquation(eq, newValue, index+1, allowConcatenation) {
		return true
	}

	// ...and *.
	newValue = value * eq.operands[index]
	if testEquation(eq, newValue, index+1, allowConcatenation) {
		return true
	}

	// And then if this is part 2, do another fork
	// to try | as well.
	if allowConcatenation {
		numDigits := math.Log10(float64(eq.operands[index])) + 1
		newValue = value * int(math.Pow10(int(numDigits)))
		newValue += eq.operands[index]
		if testEquation(eq, newValue, index+1, allowConcatenation) {
			return true
		}
	}

	return false
}
