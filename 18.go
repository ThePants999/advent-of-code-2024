package main

import (
	"log/slog"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day18 = runner.DayImplementation{
	DayNumber:          18,
	ExecutePart1:       Day18Part1,
	ExecutePart2:       Day18Part2,
	ExampleInput:       ``,
	ExamplePart1Answer: "",
	ExamplePart2Answer: "",
}

func Day18Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	result := int(lines[0][0])
	return strconv.Itoa(result), nil
}

func Day18Part2(logger *slog.Logger, input string, part1Context any) string {

	return ""
}
