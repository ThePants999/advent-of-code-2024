package main

import (
	"log/slog"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day20 = runner.DayImplementation{
	DayNumber:    20,
	ExecutePart1: Day20Part1,
	ExecutePart2: Day20Part2,
	ExampleInput: `###############
#...#...#.....#
#.#.#.#.#.###.#
#S#...#.#.#...#
#######.#.#.###
#######.#.#...#
#######.#.###.#
###..E#...#...#
###.#######.###
#...###...#...#
#.#####.#.###.#
#.#...#.#.#...#
#.#.#.#.#.#.###
#...#...#...###
###############`,
	ExamplePart1Answer: "44",
	ExamplePart2Answer: "285",
}

const (
	D20_WALL            int = -1
	D20_UNREACHED_SPACE int = -2
)

func (p gridPos) adjacenciesUnbounded() []gridPos {
	adj := make([]gridPos, 4)
	adj[0] = gridPos{p.row - 1, p.col}
	adj[1] = gridPos{p.row + 1, p.col}
	adj[2] = gridPos{p.row, p.col - 1}
	adj[3] = gridPos{p.row, p.col + 1}
	return adj
}

func Day20Part1(logger *slog.Logger, input string) (string, any) {
	// Parse the input - build up the complete grid,
	// and also record the start and end co-ordinates.
	lines := strings.Fields(input)
	numRows := len(lines)
	numCols := len(lines[0])
	grid := make([][]int, len(lines))
	var start, end gridPos
	for rowIx, row := range lines {
		grid[rowIx] = make([]int, len(row))
		for colIx, square := range row {
			var val int
			switch square {
			case '#':
				val = D20_WALL
			case '.':
				val = D20_UNREACHED_SPACE
			case 'S':
				start = gridPos{rowIx, colIx}
				val = 0
			case 'E':
				end = gridPos{rowIx, colIx}
				val = D20_UNREACHED_SPACE
			}
			grid[rowIx][colIx] = val
		}
	}

	// Figure out the path through, and record the
	// point in time at which we reach each grid square.
	cur := start
	dist := 0
	for cur != end {
		dist++
		adjs := cur.adjacenciesUnbounded()
		for _, adj := range adjs {
			if grid[adj.row][adj.col] == D20_UNREACHED_SPACE {
				grid[adj.row][adj.col] = dist
				cur = adj
				break
			}
		}
	}

	// The approach we take for part 1 is simplistic - we're going to
	// look at every wall, and figure out whether the spaces above and
	// below and/or left and right constitute a cheat that saves time
	// over the threshold.
	sum := 0
	threshold := 0
	if numRows > 20 {
		// Real input
		threshold = 100
	}
	for rowIx := 1; rowIx < numRows-1; rowIx++ {
		for colIx := 1; colIx < numCols-1; colIx++ {
			if grid[rowIx][colIx] == D20_WALL {
				top, bottom := grid[rowIx-1][colIx], grid[rowIx+1][colIx]
				left, right := grid[rowIx][colIx-1], grid[rowIx][colIx+1]
				if top >= 0 && bottom >= 0 {
					diff := bottom - top
					if top > bottom {
						diff = top - bottom
					}
					diff -= 2
					if diff >= threshold {
						sum++
					}
				}
				if left >= 0 && right >= 0 {
					diff := right - left
					if left > right {
						diff = left - right
					}
					diff -= 2
					if diff >= threshold {
						sum++
					}
				}
			}
		}
	}

	return strconv.Itoa(sum), grid
}

func Day20Part2(logger *slog.Logger, input string, part1Context any) string {
	grid := part1Context.([][]int)
	threshold := 50
	if len(grid) > 20 {
		threshold = 100
	}

	// Our approach to part 2 is a little different, though only
	// a little. This time, we're going to look at every point
	// that's on the course, and then consider the diamond of
	// points around it that you could get to within 20 time
	// units and figure out whether each one constitutes a
	// qualifying cheat. So pretty much brute force, except
	// that we can consider different starting points
	// independently and in parallel - to take straightforward
	// advantage of that, we give each row of the grid to a
	// separate goroutine.
	c := make(chan int)
	for rowIx := 1; rowIx < len(grid)-1; rowIx++ {
		go day20Part2HandleRow(grid, threshold, rowIx, c)
	}

	sum := 0
	for range len(grid) - 2 {
		sum += <-c
	}

	return strconv.Itoa(sum)
}

func day20Part2HandleRow(grid [][]int, threshold int, rowIx int, c chan int) {
	sum := 0
	for colIx := 1; colIx < len(grid[0])-1; colIx++ {
		if grid[rowIx][colIx] >= 0 {
			// Non-wall
			minRowIx := rowIx - 20
			if minRowIx < 1 {
				minRowIx = 1
			}
			maxRowIx := rowIx + 20
			if maxRowIx >= len(grid) {
				maxRowIx = len(grid) - 1
			}
			for targetRowIx := minRowIx; targetRowIx <= maxRowIx; targetRowIx++ {
				var rowDiff int
				if targetRowIx >= rowIx {
					rowDiff = targetRowIx - rowIx
				} else {
					rowDiff = rowIx - targetRowIx
				}
				remainingDist := 20 - rowDiff

				minColIx := colIx - remainingDist
				if minColIx < 1 {
					minColIx = 1
				}
				maxColIx := colIx + remainingDist
				if maxColIx >= len(grid[0]) {
					maxColIx = len(grid[0]) - 1
				}
				for targetColIx := minColIx; targetColIx <= maxColIx; targetColIx++ {
					if grid[targetRowIx][targetColIx] > grid[rowIx][colIx] {
						// This is a cheat
						var colDiff int
						if targetColIx >= colIx {
							colDiff = targetColIx - colIx
						} else {
							colDiff = colIx - targetColIx
						}
						dist := rowDiff + colDiff
						if grid[targetRowIx][targetColIx]-grid[rowIx][colIx]-dist >= threshold {
							// Legal and qualifying cheat
							sum++
						}
					}
				}
			}
		}
	}
	c <- sum
}
