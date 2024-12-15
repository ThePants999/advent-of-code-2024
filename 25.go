package main

import (
	"log/slog"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day25 = runner.DayImplementation{
	DayNumber:          25,
	ExecutePart1:       Day25Part1,
	ExecutePart2:       Day25Part2,
	ExampleInput:       ``,
	ExamplePart1Answer: "",
	ExamplePart2Answer: "",
}

func Day25Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	result := int(lines[0][0])
	return strconv.Itoa(result), nil
}

func Day25Part2(logger *slog.Logger, input string, part1Context any) string {

	return ""
}
