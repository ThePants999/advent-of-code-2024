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
	obstacles posSet
	basePath  [][]gridLocationState
	startRow  int
	startCol  int
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

	curRow, curCol := startRow, startCol
	visitedCount := 1
	for {
		var inBounds bool
		curRow, curCol, dir, inBounds = move(obstacles, curRow, curCol, len(grid), len(grid[curRow]), dir)
		if !inBounds {
			break
		}
		if !grid[curRow][curCol].visited {
			visitedCount++
			grid[curRow][curCol].visited = true
		}
		grid[curRow][curCol].visitedDirs[dir] = true
	}

	context := d6context{
		obstacles: obstacles,
		basePath:  grid,
		startRow:  startRow,
		startCol:  startCol,
	}

	return strconv.Itoa(visitedCount), context
}

func Day6Part2(logger *slog.Logger, input string, part1Context any) string {
	context := part1Context.(d6context)
	candidates := mapset.NewSet[gridPos]()
	numRows, numCols := len(context.basePath), len(context.basePath[0])

	for rowIx, row := range context.basePath {
		for colIx, locState := range row {
			if locState.visited {
				addCandidates(candidates, context.obstacles, rowIx, colIx, numRows, numCols, locState.visitedDirs)
			}
		}
	}
	candidates.Remove(gridPos{context.startRow, context.startCol})

	loopCount := 0
	for _, candidate := range candidates.ToSlice() {
		grid := make([][]gridLocationState, numRows)
		for i := 0; i < len(context.basePath[0]); i++ {
			grid[i] = make([]gridLocationState, numCols)
		}
		curRow, curCol, dir := context.startRow, context.startCol, D6_UP
		grid[curRow][curCol].visitedDirs[D6_UP] = true

		context.obstacles.Add(candidate)

		for {
			var inBounds bool
			curRow, curCol, dir, inBounds = move(context.obstacles, curRow, curCol, len(grid), len(grid[curRow]), dir)
			if !inBounds {
				break
			}
			if grid[curRow][curCol].visitedDirs[dir] {
				loopCount++
				break
			}
			grid[curRow][curCol].visitedDirs[dir] = true
		}

		context.obstacles.Remove(candidate)
	}

	return strconv.Itoa(loopCount)
}

func addCandidates(candidates posSet, obstacles posSet, row int, col int, numRows int, numCols int, dirs dirsArray) {
	for dir := D6_UP; dir <= D6_LEFT; dir++ {
		if dirs[dir] {
			newRow, newCol, inBounds := moveSimple(row, col, dir, numRows, numCols)
			newPos := gridPos{newRow, newCol}
			if inBounds && !obstacles.Contains(newPos) {
				candidates.Add(newPos)
			}
		}
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
