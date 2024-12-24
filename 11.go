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

// The approach we take here is stone counting. Let's
// say that at time T we have A stones of value W and
// B stones of value X. When you blink at a W, it
// changes into a Y and a Z. When you blink at an X,
// it changes into Y. At time T+1, we therefore have
// A+B stones of value Y, and A stones of value Z
// (and of course 0 Ws and 0 Xs). The order doesn't
// matter, so we just store how many we have of each
// stone, in a map.
//
// You'd think that some sort of caching/memoisation
// would also help, but in my testing that actually
// slowed things down vs this implementation. Map
// operations can be expensive.

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
	// Each time we blink at all the stones, we need to
	// forget the old stone counts and entirely replace
	// them with new stone counts. However, to avoid the
	// overhead of allocating a new map each time, we
	// alternate between using two different maps, simply
	// clearing each one out when we're done with it.
	secondStones := make(map[int]int)
	stones, newStones := &inputStones, &secondStones
	for i := 0; i < iterations; i++ {
		// Go through each unique stone value we currently
		// have.
		for num, count := range *stones {
			// Figure out what stones of that value turn
			// into on a blink.
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

			// Record updated quantities of the new
			// stone values.
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
	// We've still got our "what do we have after 25
	// blinks" state, so we just need to do another
	// 50 to get to 75.
	stones = doDay11Calc(stones, 50)
	return strconv.Itoa(countStones(stones))
}
