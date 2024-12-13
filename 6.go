package main

import (
	"log/slog"
	"strconv"
	"strings"

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
	obstacleCandidates map[gridPos]nothing
	startRow           int
	startCol           int
}

func Day6Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	grid := make([][]gridLocationState, len(lines))
	obstacles := make(map[gridPos]nothing)
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
				obstacles[gridPos{row, col}] = nothing{}
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
	obstacleCandidates := make(map[gridPos]nothing)
	visitedCount := 1
	for {
		var inBounds bool
		curRow, curCol, dir, inBounds = move(obstacles, curRow, curCol, numRows, numCols, dir)
		if !inBounds {
			break
		}
		if !grid[curRow][curCol].visited {
			obstacleCandidates[gridPos{curRow, curCol}] = nothing{}
			visitedCount++
			grid[curRow][curCol].visited = true
		}
		grid[curRow][curCol].visitedDirs[dir] = true
	}
	delete(obstacleCandidates, gridPos{startRow, startCol})

	return strconv.Itoa(visitedCount), d6context{obstaclesByRow, obstaclesByCol, obstacleCandidates, startRow, startCol}
}

func Day6Part2(logger *slog.Logger, input string, part1Context any) string {
	context := part1Context.(d6context)

	c := make(chan int)
	for candidate := range context.obstacleCandidates {
		go tryFindLoop(context, candidate.row, candidate.col, c)
	}

	loopCount := 0
	for i := 0; i < len(context.obstacleCandidates); i++ {
		loopCount += <-c
	}

	return strconv.Itoa(loopCount)
}

func tryFindLoop(context d6context, newObstacleRow int, newObstacleCol int, c chan int) {
	obstaclesHit := make(map[obstacleHitState]nothing, 150)

	curRow, curCol, dir := context.startRow, context.startCol, D6_UP
	for {
		var inBounds, loopDetected bool
		curRow, curCol, dir, inBounds, loopDetected = moveToNextObstacle(context.obstaclesByRow, context.obstaclesByCol, newObstacleRow, newObstacleCol, obstaclesHit, curRow, curCol, dir)
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

func moveToNextObstacle(obstaclesByRow [][]int, obstaclesByCol [][]int, newObstacleRow int, newObstacleCol int, obstaclesHit map[obstacleHitState]nothing, curRow int, curCol int, curDir direction6) (newRow int, newCol int, newDir direction6, inBounds bool, loopDetected bool) {
	newRow, newCol, inBounds, loopDetected = curRow, curCol, false, false
	obstacleHit := obstacleHitState{gridPos{curRow, curCol}, curDir}
	var obstacles []int
	var position, obstaclePosition *int
	var rightwards bool
	newObstacle := -1
	if curDir == D6_UP || curDir == D6_DOWN {
		obstacles = obstaclesByCol[curCol]
		position = &newRow
		obstaclePosition = &(obstacleHit.pos.row)
		rightwards = (curDir == D6_DOWN)
		if newObstacleCol == curCol {
			newObstacle = newObstacleRow
		}
	} else {
		obstacles = obstaclesByRow[curRow]
		position = &newCol
		obstaclePosition = &(obstacleHit.pos.col)
		rightwards = (curDir == D6_RIGHT)
		if newObstacleRow == curRow {
			newObstacle = newObstacleCol
		}
	}

	if rightwards {
		for _, obstacle := range obstacles {
			if obstacle > *position {
				if newObstacle > *position && newObstacle < obstacle {
					// We'd hit the new obstacle first.
					*obstaclePosition = newObstacle
					*position = newObstacle - 1
				} else {
					*obstaclePosition = obstacle
					*position = obstacle - 1
				}
				inBounds = true
				break
			}
		}
	} else {
		for ix := len(obstacles) - 1; ix >= 0; ix-- {
			if obstacles[ix] < *position {
				if newObstacle < *position && newObstacle > obstacles[ix] {
					// We'd hit the new obstacle first.
					*obstaclePosition = newObstacle
					*position = newObstacle + 1
				} else {
					*obstaclePosition = obstacles[ix]
					*position = obstacles[ix] + 1
				}
				inBounds = true
				break
			}
		}
	}

	if !inBounds && newObstacle > -1 {
		// Check to see if the new obstacle would have stopped us going out of bounds.
		if rightwards && newObstacle > *position {
			*obstaclePosition = newObstacle
			*position = newObstacle - 1
			inBounds = true
		} else if !rightwards && newObstacle < *position {
			*obstaclePosition = newObstacle
			*position = newObstacle + 1
			inBounds = true
		}
	}

	if inBounds {
		_, loopDetected = obstaclesHit[obstacleHit]
		if !loopDetected {
			obstaclesHit[obstacleHit] = nothing{}
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

func move(obstacles map[gridPos]nothing, curRow int, curCol int, numRows int, numCols int, curDir direction6) (newRow int, newCol int, newDir direction6, inBounds bool) {
	newDir = curDir
	newRow, newCol, inBounds = moveSimple(curRow, curCol, curDir, numRows, numCols)

	pos := gridPos{newRow, newCol}
	_, found := obstacles[pos]
	if found {
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
