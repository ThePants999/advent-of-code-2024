package main

import (
	"log/slog"
	"slices"
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

type obstacleHitState struct {
	pos gridPos
	dir direction6
}

type d6context struct {
	obstaclesByRow     [][]int
	obstaclesByCol     [][]int
	obstacleCandidates posSet
	startRow           int
	startCol           int
}

func Day6Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	grid := make([][]gridLocationState, len(lines))
	obstacles := mapset.NewSet[gridPos]()
	obstaclesByRow := make([][]int, len(lines))
	obstaclesByCol := make([][]int, len(lines[0]))
	for ix := range lines[0] {
		obstaclesByCol[ix] = make([]int, 0, 20)
	}

	startRow, startCol := -1, -1
	dir := D6_UP
	for row, line := range lines {
		obstaclesByRow[row] = make([]int, 0, 20)
		grid[row] = make([]gridLocationState, len(line))
		for col, gridItem := range line {
			if gridItem == '#' {
				obstacles.Add(gridPos{row, col})
				obstaclesByRow[row] = append(obstaclesByRow[row], col)
				obstaclesByCol[col] = append(obstaclesByCol[col], row)
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

	return strconv.Itoa(visitedCount), d6context{obstaclesByRow, obstaclesByCol, obstacleCandidates, startRow, startCol}
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
	obstaclesHit := mapset.NewSet[obstacleHitState]()
	obstaclesByRow := slices.Clone(context.obstaclesByRow)
	obstaclesByRow[newObstacleRow] = slices.Clone(context.obstaclesByRow[newObstacleRow])
	obstaclesByCol := slices.Clone(context.obstaclesByCol)
	obstaclesByCol[newObstacleCol] = slices.Clone(context.obstaclesByCol[newObstacleCol])

	added := false
	for ix, obstacle := range obstaclesByRow[newObstacleRow] {
		if obstacle > newObstacleCol {
			obstaclesByRow[newObstacleRow] = slices.Insert(obstaclesByRow[newObstacleRow], ix, newObstacleCol)
			added = true
			break
		}
	}
	if !added {
		obstaclesByRow[newObstacleRow] = append(obstaclesByRow[newObstacleRow], newObstacleCol)
	}
	added = false
	for ix, obstacle := range obstaclesByCol[newObstacleCol] {
		if obstacle > newObstacleRow {
			obstaclesByCol[newObstacleCol] = slices.Insert(obstaclesByCol[newObstacleCol], ix, newObstacleRow)
			added = true
			break
		}
	}
	if !added {
		obstaclesByCol[newObstacleCol] = append(obstaclesByCol[newObstacleCol], newObstacleRow)
	}

	curRow, curCol, dir := context.startRow, context.startCol, D6_UP
	for {
		var inBounds, loopDetected bool
		curRow, curCol, dir, inBounds, loopDetected = moveToNextObstacle(obstaclesByRow, obstaclesByCol, obstaclesHit, curRow, curCol, dir)
		if !inBounds {
			c <- 0
			return
		}
		if loopDetected {
			c <- 1
			return
		}
	}
}

func moveToNextObstacle(obstaclesByRow [][]int, obstaclesByCol [][]int, obstaclesHit mapset.Set[obstacleHitState], curRow int, curCol int, curDir direction6) (newRow int, newCol int, newDir direction6, inBounds bool, loopDetected bool) {
	newRow, newCol, inBounds, loopDetected = curRow, curCol, false, false
	obstacleHit := obstacleHitState{gridPos{curRow, curCol}, curDir}
	var obstacles []int
	var position, obstaclePosition *int
	var rightwards bool
	if curDir == D6_UP || curDir == D6_DOWN {
		obstacles = obstaclesByCol[curCol]
		position = &newRow
		obstaclePosition = &(obstacleHit.pos.row)
		rightwards = (curDir == D6_DOWN)
	} else {
		obstacles = obstaclesByRow[curRow]
		position = &newCol
		obstaclePosition = &(obstacleHit.pos.col)
		rightwards = (curDir == D6_RIGHT)
	}

	if rightwards {
		for _, obstacle := range obstacles {
			if obstacle > *position {
				*obstaclePosition = obstacle
				*position = obstacle - 1
				inBounds = true
				break
			}
		}
	} else {
		for ix := len(obstacles) - 1; ix >= 0; ix-- {
			if obstacles[ix] < *position {
				*obstaclePosition = obstacles[ix]
				*position = obstacles[ix] + 1
				inBounds = true
				break
			}
		}
	}

	if inBounds {
		if obstaclesHit.Contains(obstacleHit) {
			loopDetected = true
		} else {
			obstaclesHit.Add(obstacleHit)
		}
	}

	newDir = turnRight(curDir)

	return
}

func turnRight(curDir direction6) direction6 {
	newDir := curDir + 1
	if newDir > D6_LEFT {
		newDir = D6_UP
	}
	return newDir
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
