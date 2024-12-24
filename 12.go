package main

import (
	"log/slog"
	"strconv"
	"strings"

	stack "github.com/golang-collections/collections/stack"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day12 = runner.DayImplementation{
	DayNumber:    12,
	ExecutePart1: Day12Part1,
	ExecutePart2: Day12Part2,
	ExampleInput: `RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE`,
	ExamplePart1Answer: "1930",
	ExamplePart2Answer: "1206",
}

type nothing struct{}

type region struct {
	char     rune
	plots    [][]bool
	plotsArr []gridPos
}

func newRegion(char rune, gridSize int) region {
	r := region{char, make([][]bool, gridSize), make([]gridPos, 0, 20)}
	for i := range gridSize {
		r.plots[i] = make([]bool, gridSize)
	}
	return r
}

func Day12Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	grid := make([][]rune, len(lines))
	usedplots := make([][]bool, len(lines))
	for rowIx, line := range lines {
		grid[rowIx] = make([]rune, len(line))
		usedplots[rowIx] = make([]bool, len(line))
		for colIx, char := range line {
			grid[rowIx][colIx] = char
		}
	}

	regions := make([]region, 0, 100)
	for rowIx, row := range grid {
		for colIx, char := range row {
			loc := gridPos{rowIx, colIx}
			if !usedplots[loc.row][loc.col] {
				// New region
				region := newRegion(char, len(grid))
				s := stack.New()
				s.Push(loc)
				for s.Len() > 0 {
					loc = s.Pop().(gridPos)
					if grid[loc.row][loc.col] == region.char && !usedplots[loc.row][loc.col] {
						usedplots[loc.row][loc.col] = true
						region.plots[loc.row][loc.col] = true
						region.plotsArr = append(region.plotsArr, loc)
						adjs := loc.adjacencies(len(grid), len(grid[0]))
						for _, adj := range adjs {
							s.Push(adj)
						}
					}
				}
				regions = append(regions, region)
			}
		}
	}

	price := calcTotalPrice(&regions, func(r *region, pos gridPos) int {
		fences := 4
		for _, adj := range pos.adjacencies(len(grid), len(grid[0])) {
			if r.plots[adj.row][adj.col] {
				fences--
			}
		}
		return fences
	})

	return strconv.Itoa(price), &regions
}

func calcTotalPrice(regions *[]region, plotWeight func(*region, gridPos) int) int {
	price := 0
	c := make(chan int)
	for _, reg := range *regions {
		go reg.calcPrice(plotWeight, c)
	}
	for range *regions {
		price += <-c
	}
	return price
}

func (r *region) calcPrice(plotWeight func(*region, gridPos) int, c chan int) {
	weight := 0

	for _, p := range r.plotsArr {
		weight += plotWeight(r, p)
	}

	c <- weight * len(r.plotsArr)
}

func Day12Part2(logger *slog.Logger, input string, part1Context any) string {
	regions := part1Context.(*[]region)
	price := calcTotalPrice(regions, func(r *region, pos gridPos) int {
		vertices := 0

		allAdj := [9]gridPos{
			{pos.row - 1, pos.col},
			{pos.row - 1, pos.col + 1},
			{pos.row, pos.col + 1},
			{pos.row + 1, pos.col + 1},
			{pos.row + 1, pos.col},
			{pos.row + 1, pos.col - 1},
			{pos.row, pos.col - 1},
			{pos.row - 1, pos.col - 1},
			{pos.row - 1, pos.col},
		}
		var inRegion [9]bool
		for ix := 0; ix < 9; ix++ {
			adj := allAdj[ix]
			inRegion[ix] = adj.row >= 0 && adj.row < len(r.plots) && adj.col >= 0 && adj.col < len(r.plots[adj.row]) && r.plots[allAdj[ix].row][allAdj[ix].col]
		}

		for dir := UP_RIGHT; dir <= UP_LEFT; dir += 2 {
			if !inRegion[dir-1] && !inRegion[dir+1] {
				// Convex vertex
				vertices++
			} else if inRegion[dir-1] && inRegion[dir+1] && !inRegion[dir] {
				// Concave vertex
				vertices++
			}
		}

		return vertices
	})

	return strconv.Itoa(price)
}
