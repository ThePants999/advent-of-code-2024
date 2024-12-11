package main

import (
	"log/slog"
	"maps"
	"math"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

type intPair struct {
	one int
	two int
}

// I briefly tried caching the result of a blink on specific
// stone values; interestingly it actually slowed it down.
//var day11Cache map[int]intPair

var Day11 = runner.DayImplementation{
	DayNumber:          11,
	ExecutePart1:       Day11Part1,
	ExecutePart2:       Day11Part2,
	ExampleInput:       "125 17",
	ExamplePart1Answer: "55312",
	ExamplePart2Answer: "65601038650482",
}

func Day11Part1(logger *slog.Logger, input string) (string, any) {
	//day11Cache = make(map[int]intPair)
	first := parseDay11Input(input)
	return strconv.Itoa(doDay11Calc(first, 25)), nil
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

func doDay11Calc(stones map[int]int, iterations int) int {
	stones = maps.Clone(stones)
	for i := 0; i < iterations; i++ {
		newStones := make(map[int]int)
		for num, count := range stones {
			var newNums intPair
			/*newNums, found := day11Cache[num]
			if !found {*/
			if num == 0 {
				newNums = intPair{1, -1}
			} else {
				numDigits := int(math.Log10((float64)(num))) + 1
				if numDigits%2 == 0 {
					divisor := int(math.Pow10(numDigits / 2))
					newNums = intPair{num % divisor, num / divisor}
				} else {
					newNums = intPair{num * 2024, -1}
				}
			}
			/*day11Cache[num] = newNums
			}*/
			existingNum1 := newStones[newNums.one]
			newStones[newNums.one] = existingNum1 + count
			if newNums.two >= 0 {
				existingNum2 := newStones[newNums.two]
				newStones[newNums.two] = existingNum2 + count
			}
		}
		stones = newStones
	}

	sum := 0
	for _, count := range stones {
		sum += count
	}

	return sum
}

func Day11Part2(logger *slog.Logger, input string, part1Context any) string {
	first := parseDay11Input(input)
	return strconv.Itoa(doDay11Calc(first, 75))
}
