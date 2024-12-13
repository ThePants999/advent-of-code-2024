package main

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/set"
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
	char  rune
	plots map[gridPos]nothing
}

func Day12Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	grid := make([][]rune, len(lines))
	for rowIx, line := range lines {
		grid[rowIx] = make([]rune, len(line))
		for colIx, char := range line {
			grid[rowIx][colIx] = char
		}
	}

	regions := make([]region, 0, 100)
	usedplots := set.New()
	for rowIx, row := range grid {
		for colIx, char := range row {
			loc := gridPos{rowIx, colIx}
			if !usedplots.Has(loc) {
				// New region
				region := region{char, make(map[gridPos]nothing)}
				s := stack.New()
				s.Push(loc)
				for s.Len() > 0 {
					loc = s.Pop().(gridPos)
					if grid[loc.row][loc.col] == region.char && !usedplots.Has(loc) {
						usedplots.Insert(loc)
						region.plots[loc] = nothing{}
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
			_, found := r.plots[adj]
			if found {
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

	for p := range r.plots {
		weight += plotWeight(r, p)
	}

	c <- weight * len(r.plots)
}

func Day12Part2(logger *slog.Logger, input string, part1Context any) string {
	regions := part1Context.(*[]region)
	price := calcTotalPrice(regions, func(r *region, pos gridPos) int {
		vertices := 0

		// Conveniently, we don't have to worry about bounds checking here
		// as we only care whether an adjacent plot is within the region,
		// so we don't care whether the answer is "no" because it's a valid
		// plot in a different region or because it's an invalid plot.
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
			_, inRegion[ix] = r.plots[allAdj[ix]]
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
