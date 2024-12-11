package main

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/set"
	stack "github.com/golang-collections/collections/stack"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day10 = runner.DayImplementation{
	DayNumber:    10,
	ExecutePart1: Day10Part1,
	ExecutePart2: Day10Part2,
	ExampleInput: `89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732`,
	ExamplePart1Answer: "36",
	ExamplePart2Answer: "81",
}

type day10context struct {
	grid       [][]int
	trailheads []gridPos
}

func Day10Part1(logger *slog.Logger, input string) (string, any) {
	rows := strings.Fields(input)
	grid := make([][]int, len(rows))
	trailheads := make([]gridPos, 0, len(input))
	for rowIx, row := range rows {
		grid[rowIx] = make([]int, len(row))
		for colIx, square := range row {
			val := int(square - '0')
			grid[rowIx][colIx] = val
			if val == 0 {
				trailheads = append(trailheads, gridPos{rowIx, colIx})
			}
		}
	}

	return strconv.Itoa(day10dfs(grid, trailheads, true)), day10context{grid, trailheads}
}

func (p gridPos) adjacencies(numRows int, numCols int) []gridPos {
	adj := make([]gridPos, 0, 4)
	if p.row > 0 {
		adj = append(adj, gridPos{p.row - 1, p.col})
	}
	if p.row < numRows-1 {
		adj = append(adj, gridPos{p.row + 1, p.col})
	}
	if p.col > 0 {
		adj = append(adj, gridPos{p.row, p.col - 1})
	}
	if p.col < numCols-1 {
		adj = append(adj, gridPos{p.row, p.col + 1})
	}
	return adj
}

func day10dfs(grid [][]int, trailheads []gridPos, visitedCheck bool) int {
	sum := 0
	for _, trailhead := range trailheads {
		visited := set.New()
		s := stack.New()
		s.Push(trailhead)
		for s.Len() > 0 {
			pos := s.Pop().(gridPos)
			if !visitedCheck || !visited.Has(pos) {
				visited.Insert(pos)
				if grid[pos.row][pos.col] == 9 {
					sum++
				} else {
					for _, adj := range pos.adjacencies(len(grid), len(grid[pos.row])) {
						if grid[adj.row][adj.col]-grid[pos.row][pos.col] == 1 {
							s.Push(adj)
						}
					}
				}
			}
		}
	}

	return sum
}

func Day10Part2(logger *slog.Logger, input string, part1Context any) string {
	context := part1Context.(day10context)
	return strconv.Itoa(day10dfs(context.grid, context.trailheads, false))
}
