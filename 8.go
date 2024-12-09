package main

import (
	"log/slog"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/mowshon/iterium"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day8 = runner.DayImplementation{
	DayNumber:    8,
	ExecutePart1: Day8Part1,
	ExecutePart2: Day8Part2,
	ExampleInput: `............
........0...
.....0......
.......0....
....0.......
......A.....
............
............
........A...
.........A..
............
............`,
	ExamplePart1Answer: "14",
	ExamplePart2Answer: "34",
}

type day8context struct {
	combinations [][]gridPos
	numRows      int
	numCols      int
}

func Day8Part1(logger *slog.Logger, input string) (string, any) {
	rows := strings.Fields(input)
	antennae := make(map[rune][]gridPos)
	for rowIx, row := range rows {
		for colIx, frequency := range row {
			if frequency != '.' {
				list, found := antennae[frequency]
				if !found {
					list = make([]gridPos, 0, 10)
				}
				antennae[frequency] = append(list, gridPos{rowIx, colIx})
			}
		}
	}

	allCombinations := make([][]gridPos, 0, 1000)
	for _, coords := range antennae {
		combinations := iterium.Combinations(coords, 2)
		slice, _ := combinations.Slice()
		allCombinations = append(allCombinations, slice...)
	}

	numRows, numCols := len(rows), len(rows[0])
	context := day8context{allCombinations, numRows, numCols}

	set := mapset.NewSet[gridPos]()
	for _, combination := range allCombinations {
		locationA := gridPos{
			combination[0].row + combination[0].row - combination[1].row,
			combination[0].col + combination[0].col - combination[1].col}
		locationB := gridPos{
			combination[1].row + combination[1].row - combination[0].row,
			combination[1].col + combination[1].col - combination[0].col}
		if locationInGrid(locationA, numRows, numCols) {
			set.Add(locationA)
		}
		if locationInGrid(locationB, numRows, numCols) {
			set.Add(locationB)
		}
	}

	return strconv.Itoa(set.Cardinality()), context
}

func locationInGrid(loc gridPos, numRows int, numCols int) bool {
	return loc.row >= 0 && loc.row < numRows && loc.col >= 0 && loc.col < numCols
}

func (this gridPos) Add(other gridPos) gridPos {
	return gridPos{this.row + other.row, this.col + other.col}
}

func (this gridPos) Subtract(other gridPos) gridPos {
	return gridPos{this.row - other.row, this.col - other.col}
}

func Day8Part2(logger *slog.Logger, input string, part1Context any) string {
	context := part1Context.(day8context)

	set := mapset.NewSet[gridPos]()
	for _, combination := range context.combinations {
		start := combination[0]
		delta := combination[1].Subtract(start)
		for next := start; locationInGrid(next, context.numRows, context.numCols); next = next.Add(delta) {
			set.Add(next)
		}
		for next := start; locationInGrid(next, context.numRows, context.numCols); next = next.Subtract(delta) {
			set.Add(next)
		}
	}

	return strconv.Itoa(set.Cardinality())
}
