package main

import (
	"log/slog"
	"strconv"
	"strings"

	cmap "github.com/orcaman/concurrent-map/v2"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day19 = runner.DayImplementation{
	DayNumber:    19,
	ExecutePart1: Day19Part1,
	ExecutePart2: Day19Part2,
	ExampleInput: `r, wr, b, g, bwu, rb, gb, br

brwrr
bggr
gbbr
rrbgbr
ubwu
bwurrg
brgr
bbrgwb`,
	ExamplePart1Answer: "6",
	ExamplePart2Answer: "16",
}

var towels map[string]nothing
var minTowelLen, maxTowelLen int
var solutions cmap.ConcurrentMap[string, int]

func Day19Part1(logger *slog.Logger, input string) (string, any) {
	solutions = cmap.New[int]()
	minTowelLen = 99
	towels = make(map[string]nothing)
	lines := strings.Split(input, "\n")
	towelsStr := strings.Split(lines[0], ", ")
	for _, towel := range towelsStr {
		towels[towel] = nothing{}
		if len(towel) < minTowelLen {
			minTowelLen = len(towel)
		}
		if len(towel) > maxTowelLen {
			maxTowelLen = len(towel)
		}
	}
	patterns := lines[2:]

	// Solve each pattern on a separate thread.
	c := make(chan int)
	for _, pattern := range patterns {
		go solvePattern(pattern, c)
	}

	// Calculate both parts simultaneously.
	sum := 0
	count := 0
	for range patterns {
		result := <-c
		sum += result
		if result > 0 {
			count++
		}
	}

	return strconv.Itoa(count), sum
}

func solvePattern(pattern string, c chan int) {
	// The input may have blank lines.
	if len(pattern) == 0 {
		c <- 0
		return
	}

	c <- solvePatternRecursive(pattern)
}

func solvePatternRecursive(pattern string) int {
	// If we get down to an empty string, we've found a match.
	if len(pattern) == 0 {
		return 1
	}

	// See whether we've deconstructed exactly this sub-pattern before.
	result, found := solutions.Get(pattern)
	if found {
		return result
	}

	// Check each head length of the current sub-pattern that might
	// match a towel.
	maxLen := maxTowelLen
	if len(pattern) < maxLen {
		maxLen = len(pattern)
	}
	for i := minTowelLen; i <= maxLen; i++ {
		_, found := towels[pattern[:i]]
		if found {
			result += solvePatternRecursive(pattern[i:])
		}
	}

	// Remember the result for this sub-pattern.
	solutions.Set(pattern, result)
	return result
}

func Day19Part2(logger *slog.Logger, input string, part1Context any) string {
	sum := part1Context.(int)
	return strconv.Itoa(sum)
}
