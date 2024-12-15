package main

import (
	"log/slog"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day22 = runner.DayImplementation{
	DayNumber:          22,
	ExecutePart1:       Day22Part1,
	ExecutePart2:       Day22Part2,
	ExampleInput:       ``,
	ExamplePart1Answer: "",
	ExamplePart2Answer: "",
}

func Day22Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	result := int(lines[0][0])
	return strconv.Itoa(result), nil
}

func Day22Part2(logger *slog.Logger, input string, part1Context any) string {

	return ""
}
