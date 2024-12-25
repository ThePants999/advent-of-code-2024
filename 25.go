package main

import (
	"log/slog"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day25 = runner.DayImplementation{
	DayNumber:    25,
	ExecutePart1: Day25Part1,
	ExecutePart2: Day25Part2,
	ExampleInput: `#####
.####
.####
.####
.#.#.
.#...
.....

#####
##.##
.#.##
...##
...#.
...#.
.....

.....
#....
#....
#...#
#.#.#
#.###
#####

.....
.....
#.#..
###..
###.#
###.#
#####

.....
.....
.....
#....
#.#..
#.#.#
#####`,
	ExamplePart1Answer: "3",
	ExamplePart2Answer: "",
}

// You don't need comments today. Merry Christmas!

func Day25Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	locks := make([][]int, 0, len(lines)/7)
	keys := make([][]int, 0, len(lines)/7)
	for ix := 0; ix < len(lines); ix += 7 {
		seq := make([]int, 5)
		isALock := lines[ix][0] == '#'
		first := ix
		last := ix + 7
		if isALock {
			first++
		} else {
			last--
		}
		for _, row := range lines[first:last] {
			for colIx, char := range row {
				if char == '#' {
					seq[colIx]++
				}
			}
		}
		if isALock {
			locks = append(locks, seq)
		} else {
			keys = append(keys, seq)
		}
	}

	sum := 0
	for _, lock := range locks {
	key:
		for _, key := range keys {
			for i := 0; i < 5; i++ {
				if lock[i]+key[i] > 5 {
					continue key
				}
			}
			sum++
		}
	}

	return strconv.Itoa(sum), nil
}

func Day25Part2(logger *slog.Logger, input string, part1Context any) string {

	return ""
}
