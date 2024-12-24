package main

import (
	"log/slog"
	"slices"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day1 = runner.DayImplementation{
	DayNumber:          1,
	ExecutePart1:       Day1Part1,
	ExecutePart2:       Day1Part2,
	ExampleInput:       "3   4\n4   3\n2   5\n1   3\n3   9\n3   3",
	ExamplePart1Answer: "11",
	ExamplePart2Answer: "31",
}

func Day1Part1(logger *slog.Logger, input string) (string, any) {
	// Split on any amount of whitespace, following which
	// we can direct even and odd elements into separate
	// lists.
	numbers := strings.Fields(input)
	var list1, list2 []int = make([]int, 0, len(numbers)/2), make([]int, 0, len(numbers)/2)

	for ix, num := range numbers {
		numInt, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		if ix%2 == 0 {
			list1 = append(list1, numInt)
		} else {
			list2 = append(list2, numInt)
		}
	}

	// Part 1 just requires that we sort the lists before
	// performing an element-wise comparison. We'll leave
	// that to the standard library.
	slices.Sort(list1)
	slices.Sort(list2)

	distanceSum := 0
	for ix, list1Num := range list1 {
		distance := list2[ix] - list1Num
		if distance < 0 {
			distance *= -1
		}
		distanceSum += distance
	}

	return strconv.Itoa(distanceSum), [][]int{list1, list2}
}

func Day1Part2(logger *slog.Logger, input string, part1Context any) string {
	lists := part1Context.([][]int)
	list1, list2 := lists[0], lists[1]

	// Calculate frequencies by iterating through the
	// right-hand list and incrementing values in a map.
	m := make(map[int]int)
	for _, val := range list2 {
		m[val] = m[val] + 1
	}

	// Similarity scores are then elements in the left-
	// hand list multiplied by their frequency in the
	// right.
	similarityScore := 0
	for _, val := range list1 {
		similarityScore += val * m[val]
	}

	return strconv.Itoa(similarityScore)
}
