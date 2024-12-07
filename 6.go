package main

import (
	"log/slog"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day6 = runner.DayImplementation{
	DayNumber:    6,
	ExecutePart1: Day6Part1,
	ExecutePart2: Day6Part2,
	ExampleInput: `....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...`,
	ExamplePart1Answer: "41",
	ExamplePart2Answer: "6",
}

type direction6 int

const (
	D6_UP direction6 = iota
	D6_RIGHT
	D6_DOWN
	D6_LEFT
)

type gridPos struct {
	row int
	col int
}
type dirsArray [D6_LEFT + 1]bool
type posSet mapset.Set[gridPos]

type gridLocationState struct {
	visited     bool
	visitedDirs dirsArray
}

type d6context struct {
	obstacles          posSet
	obstacleCandidates posSet
	startRow           int
	startCol           int
	numRows            int
	numCols            int
}

func Day6Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	grid := make([][]gridLocationState, len(lines))
	obstacles := mapset.NewSet[gridPos]()
	startRow, startCol := -1, -1
	dir := D6_UP
	for row, line := range lines {
		grid[row] = make([]gridLocationState, len(line))
		for col, gridItem := range line {
			if gridItem == '#' {
				obstacles.Add(gridPos{row, col})
			} else if gridItem == '^' {
				startRow, startCol = row, col
				grid[row][col].visited = true
				grid[row][col].visitedDirs[D6_UP] = true
			}
		}
	}

	numRows, numCols := len(grid), len(grid[0])
	curRow, curCol := startRow, startCol
	obstacleCandidates := mapset.NewSet[gridPos]()
	visitedCount := 1
	for {
		var inBounds bool
		curRow, curCol, dir, inBounds = move(obstacles, curRow, curCol, numRows, numCols, dir)
		if !inBounds {
			break
		}
		if !grid[curRow][curCol].visited {
			obstacleCandidates.Add(gridPos{curRow, curCol})
			visitedCount++
			grid[curRow][curCol].visited = true
		}
		grid[curRow][curCol].visitedDirs[dir] = true
	}
	obstacleCandidates.Remove(gridPos{startRow, startCol})

	return strconv.Itoa(visitedCount), d6context{obstacles, obstacleCandidates, startRow, startCol, numRows, numCols}
}

func Day6Part2(logger *slog.Logger, input string, part1Context any) string {
	context := part1Context.(d6context)

	c := make(chan int)
	for _, candidate := range context.obstacleCandidates.ToSlice() {
		go tryFindLoop(context, candidate.row, candidate.col, c)
	}

	loopCount := 0
	for i := 0; i < context.obstacleCandidates.Cardinality(); i++ {
		loopCount += <-c
	}

	return strconv.Itoa(loopCount)
}

func tryFindLoop(context d6context, newObstacleRow int, newObstacleCol int, c chan int) {
	grid := make([][]gridLocationState, context.numRows)
	for i := 0; i < context.numRows; i++ {
		grid[i] = make([]gridLocationState, context.numCols)
	}
	curRow, curCol, dir := context.startRow, context.startCol, D6_UP
	grid[curRow][curCol].visitedDirs[D6_UP] = true

	obstacles := context.obstacles.Clone()
	obstacles.Add(gridPos{newObstacleRow, newObstacleCol})

	for {
		var inBounds bool
		curRow, curCol, dir, inBounds = move(obstacles, curRow, curCol, context.numRows, context.numCols, dir)
		if !inBounds {
			c <- 0
			return
		}
		if grid[curRow][curCol].visitedDirs[dir] {
			c <- 1
			return
		}
		grid[curRow][curCol].visitedDirs[dir] = true
	}
}

func move(obstacles posSet, curRow int, curCol int, numRows int, numCols int, curDir direction6) (newRow int, newCol int, newDir direction6, inBounds bool) {
	newDir = curDir
	newRow, newCol, inBounds = moveSimple(curRow, curCol, curDir, numRows, numCols)

	pos := gridPos{newRow, newCol}
	if obstacles.Contains(pos) {
		newRow, newCol = curRow, curCol
		newDir++
		if newDir > D6_LEFT {
			newDir = D6_UP
		}
	}
	return
}

func moveSimple(curRow int, curCol int, dir direction6, numRows int, numCols int) (newRow int, newCol int, inBounds bool) {
	newRow, newCol, inBounds = curRow, curCol, true
	switch dir {
	case D6_UP:
		newRow--
	case D6_DOWN:
		newRow++
	case D6_LEFT:
		newCol--
	case D6_RIGHT:
		newCol++
	}
	if newRow < 0 || newRow == numRows || newCol < 0 || newCol == numCols {
		inBounds = false
	}
	return
}
