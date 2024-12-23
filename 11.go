package main

import (
	"log/slog"
	"math"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

type intPair struct {
	one int
	two int
}

var Day11 = runner.DayImplementation{
	DayNumber:          11,
	ExecutePart1:       Day11Part1,
	ExecutePart2:       Day11Part2,
	ExampleInput:       "125 17",
	ExamplePart1Answer: "55312",
	ExamplePart2Answer: "65601038650482",
}

func Day11Part1(logger *slog.Logger, input string) (string, any) {
	first := parseDay11Input(input)
	stones := doDay11Calc(first, 25)
	return strconv.Itoa(countStones(stones)), stones
}

func parseDay11Input(input string) map[int]int {
	stones := make(map[int]int)
	nums := strings.Fields(input)
	for _, numString := range nums {
		num, _ := strconv.Atoi(numString)
		prev := stones[num]
		stones[num] = prev + 1
	}
	return stones
}

func countDigits(num int) int {
	// At one point, profiling showed 20% of my runtime going
	// on math.Log10(). What we're trying to use it for is
	// relatively simple, so here's a dumb implementation.
	//
	// Yes, it sickens me too. WHATEVER I DON'T CARE THIS GOT
	// ME BELOW 10MS
	switch {
	case num > 99999999999:
		return 12
	case num > 9999999999:
		return 11
	case num > 999999999:
		return 10
	case num > 99999999:
		return 9
	case num > 9999999:
		return 8
	case num > 999999:
		return 7
	case num > 99999:
		return 6
	case num > 9999:
		return 5
	case num > 999:
		return 4
	case num > 99:
		return 3
	case num > 9:
		return 2
	default:
		return 1
	}
}

func doDay11Calc(inputStones map[int]int, iterations int) map[int]int {
	secondStones := make(map[int]int)
	stones, newStones := &inputStones, &secondStones
	for i := 0; i < iterations; i++ {
		for num, count := range *stones {
			var newNums intPair
			if num == 0 {
				newNums = intPair{1, -1}
			} else {
				numDigits := countDigits(num)
				if numDigits%2 == 0 {
					divisor := int(math.Pow10(numDigits / 2))
					newNums = intPair{num % divisor, num / divisor}
				} else {
					newNums = intPair{num * 2024, -1}
				}
			}

			existingNum1 := (*newStones)[newNums.one]
			(*newStones)[newNums.one] = existingNum1 + count
			if newNums.two >= 0 {
				existingNum2 := (*newStones)[newNums.two]
				(*newStones)[newNums.two] = existingNum2 + count
			}
		}

		clear(*stones)
		stones, newStones = newStones, stones
	}

	return *stones
}

func countStones(stones map[int]int) int {
	sum := 0
	for _, count := range stones {
		sum += count
	}

	return sum
}

func Day11Part2(logger *slog.Logger, input string, part1Context any) string {
	stones := part1Context.(map[int]int)
	stones = doDay11Calc(stones, 50)
	return strconv.Itoa(countStones(stones))
}
