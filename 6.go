package main

import (
	"log/slog"
	"maps"
	"runtime"
	"slices"
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

type obstacleHitState struct {
	pos gridPos
	dir direction6
}

type d6context struct {
	obstaclesByRow     [][]int
	obstaclesByCol     [][]int
	obstacleCandidates []gridPos
	startRow           int
	startCol           int
}

func Day6Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	visited := make([][]bool, len(lines))
	obstacles := make(map[gridPos]nothing)
	obstaclesByRow := make([][]int, len(lines))
	obstaclesByCol := make([][]int, len(lines[0]))
	for ix := range lines[0] {
		obstaclesByCol[ix] = make([]int, 0, 20)
	}

	// Parse the input. We're looking to build up the following.
	// - startRow and startCol tell us where the guard starts.
	// - grid will be used to track the guard's movements, and
	//   records where the obstacles are in a way that's useful
	//   for part 1.
	// - obstaclesByRow records, for each row, the column indexes
	//   that contain obstacles. Additionally recording obstacles
	//   this way helps for part 2, where we don't need to
	//   simulate the guard moving square by square but rather can
	//   "teleport" him to the next obstacle in a given line.
	// - obstaclesByCol is similar.
	startRow, startCol := -1, -1
	dir := D6_UP
	for row, line := range lines {
		obstaclesByRow[row] = make([]int, 0, 20)
		visited[row] = make([]bool, len(line))
		for col, gridItem := range line {
			if gridItem == '#' {
				obstacles[gridPos{row, col}] = nothing{}
				obstaclesByRow[row] = append(obstaclesByRow[row], col)
				obstaclesByCol[col] = append(obstaclesByCol[col], row)
			} else if gridItem == '^' {
				startRow, startCol = row, col
				visited[row][col] = true
			}
		}
	}

	// Simulate the guard moving around the grid.
	numRows, numCols := len(visited), len(visited[0])
	curRow, curCol := startRow, startCol
	obstacleCandidates := make(map[gridPos]nothing)
	visitedCount := 1
	for {
		var inBounds bool
		curRow, curCol, dir, inBounds = move(obstacles, curRow, curCol, numRows, numCols, dir)
		if !inBounds {
			// The guard has left the grid - we're done.
			break
		}
		if !visited[curRow][curCol] {
			// This is the first time the guard has entered this space.
			// Increase our count, which is our part 1 answer, and also
			// record this space, as every space the guard visits is
			// somewhere we'll need to consider generating a new obstacle
			// in part 2.
			obstacleCandidates[gridPos{curRow, curCol}] = nothing{}
			visitedCount++
			visited[curRow][curCol] = true
		}
	}

	// Make sure we don't try to spawn an obstacle on top of the guard.
	delete(obstacleCandidates, gridPos{startRow, startCol})

	return strconv.Itoa(visitedCount), d6context{obstaclesByRow, obstaclesByCol, slices.Collect(maps.Keys(obstacleCandidates)), startRow, startCol}
}

func Day6Part2(logger *slog.Logger, input string, part1Context any) string {
	context := part1Context.(d6context)

	// The basic idea of how we tackle part 2 is that we're going to
	// try spawning an obstacle at every location the guard visited
	// in part 1.

	// A bit of optimisation. In an earlier version of this code, we
	// spawned a goroutine for every obstacle location we wanted to
	// try. Every instance of trying an obstacle location requires a
	// map (set, really) of obstacles that the guard has hit - but
	// with ~2500 candidate obstacle locations in my input, we were
	// spending half our runtime just on the memory allocations within
	// those goroutines. So what we do now is to spawn only as many
	// goroutines as we have CPU cores, and have each one process a
	// number of obstacle candidates in series, reusing the same map
	// for each candidate, which vastly brings down the allocations.
	threads := runtime.NumCPU()
	obstaclesPerThread := len(context.obstacleCandidates) / threads
	c := make(chan int)

	for ix := range threads {
		go tryFindLoop(context, ix*obstaclesPerThread, obstaclesPerThread, c)
	}
	remainder := len(context.obstacleCandidates) % threads
	if remainder > 0 {
		go tryFindLoop(context, threads*obstaclesPerThread, remainder, c)
	}

	// Every candidate tried will send either a 0 or a 1 on the channel
	// depending on whether that candidate caused a loop.
	loopCount := 0
	for range len(context.obstacleCandidates) {
		loopCount += <-c
	}

	return strconv.Itoa(loopCount)
}

func tryFindLoop(context d6context, firstObstacleIx int, numObstacles int, c chan int) {
	obstaclesHit := make(map[obstacleHitState]nothing, 150)

	// As explained above, this function will process a number of obstacle
	// candidates in series.
	for ix := range numObstacles {
		newObstacle := context.obstacleCandidates[firstObstacleIx+ix]
		curRow, curCol, dir := context.startRow, context.startCol, D6_UP
		for {
			var inBounds, loopDetected bool
			// Unlike in part 1, we don't need to move square by square. The
			// guard will walk forwards until he either hits an obstacle or
			// exits the grid, so we just figure out what's next in his path
			// and teleport him straight there.
			curRow, curCol, dir, inBounds, loopDetected = moveToNextObstacle(context.obstaclesByRow, context.obstaclesByCol, newObstacle.row, newObstacle.col, obstaclesHit, curRow, curCol, dir)
			if !inBounds {
				// The guard has left the grid before we detected a loop,
				// so this obstacle candidate didn't cause a loop.
				c <- 0
				break
			}
			if loopDetected {
				// The guard hit an obstacle that he's hit before, in the
				// same direction, which means he's now looping. Report
				// that, and we're done.
				c <- 1
				break
			}
		}

		// Clear out the record of hit obstacles ready for the next candidate.
		clear(obstaclesHit)
	}
}

func moveToNextObstacle(obstaclesByRow [][]int, obstaclesByCol [][]int, newObstacleRow int, newObstacleCol int, obstaclesHit map[obstacleHitState]nothing, curRow int, curCol int, curDir direction6) (newRow int, newCol int, newDir direction6, inBounds bool, loopDetected bool) {
	// As noted above, we don't need to move the guard square by
	// square. We just want to figure out what's next in his path
	// and teleport him straight there.

	// If we're moving up or down, then the column isn't going to
	// change from the current one, and we're going to hit an
	// obstacle on the current column. Similarly, if we're moving
	// left or right, the row isn't going to change and we'll hit
	// an obstacle on this row. To simplify the code somewhat, we
	// therefore initialise the final answers to the current
	// position, and then create some pointers to the variables
	// that we're going to change, so that later code doesn't care
	// whether it's row or column that changes.
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

	// "rightwards" means row or column is increasing, so the
	// next obstacle we'd hit is the one AFTER our current
	// position.
	if rightwards {
		// Find the first obstacle that's after our current
		// position.
		for _, obstacle := range obstacles {
			if obstacle > *position {
				// To avoid moving lots of memory around, we
				// don't actually add the candidate obstacle
				// to the obstacles list - we just check
				// against it separately.
				if newObstacle > *position && newObstacle < obstacle {
					// We'd hit the new obstacle first.
					*obstaclePosition = newObstacle
					*position = newObstacle - 1
				} else {
					// We'd hit the obstacle we just found.
					*obstaclePosition = obstacle
					*position = obstacle - 1
				}
				inBounds = true
				break
			}
		}
	} else {
		// We're moving "leftwards" (left or up), meaning we'd
		// hit the obstacle immediately BEFORE our current
		// position in the obstacles array. Work backwards
		// from the end, so the first one we find with an
		// earlier position is the one we'd hit.
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

	// The code above won't catch us hitting the new obstacle
	// if it's after the last one while heading rightwards,
	// or before the first one while heading leftwards. We'll
	// separately check that here.
	if !inBounds && newObstacle > -1 {
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
		// We've stayed in bounds, so we must have hit an
		// obstacle. See if we've already hit that obstacle
		// before, in the same direction. If we have, the
		// guard is now looping.
		_, loopDetected = obstaclesHit[obstacleHit]
		if !loopDetected {
			// First time hitting this obstacle in this
			// direction - record it.
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
