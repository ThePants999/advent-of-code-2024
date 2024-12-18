package main

import (
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/queue"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day18 = runner.DayImplementation{
	DayNumber:    18,
	ExecutePart1: Day18Part1,
	ExecutePart2: Day18Part2,
	ExampleInput: `5,4
4,2
4,5
3,0
2,1
6,3
2,4
1,5
0,6
3,3
2,6
5,1
1,2
5,5
2,5
6,5
1,4
0,4
6,4
1,1
6,1
1,0
0,5
1,6
2,0`,
	ExamplePart1Answer: "22",
	ExamplePart2Answer: "6,1",
}

const GRID_SIZE_EXAMPLE int = 7
const GRID_SIZE_REAL int = 71
const START_AFTER_EXAMPLE int = 12
const START_AFTER_REAL int = 1024

type d18Context struct {
	walls    [][]int
	gridSize int
	minTime  int
	lines    []string
}

func Day18Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	gridSize := GRID_SIZE_EXAMPLE
	startAfter := START_AFTER_EXAMPLE
	if len(lines) > 100 {
		gridSize = GRID_SIZE_REAL
		startAfter = START_AFTER_REAL
	}

	walls := make([][]int, gridSize)
	for rowIx := range gridSize {
		walls[rowIx] = make([]int, gridSize)
		for colIx := range gridSize {
			walls[rowIx][colIx] = math.MaxInt
		}
	}

	// _walls_ tells us at what point in time a wall appears
	// at that grid location.  (We initialize above to MaxInt,
	// effectively meaning "never".)
	for lineIx, line := range lines {
		commaIx := strings.IndexRune(line, ',')
		colIx, _ := strconv.Atoi(line[:commaIx])
		rowIx, _ := strconv.Atoi(line[commaIx+1:])
		walls[rowIx][colIx] = lineIx
	}

	pathLen := runMaze(walls, startAfter-1, gridSize)
	return strconv.Itoa(pathLen), d18Context{walls, gridSize, startAfter, lines}
}

// Simulate running the maze at a given time - which means the walls that are
// considered to be in place are those where walls[rowIx][colIx] <= time.
func runMaze(walls [][]int, time int, gridSize int) int {
	visited := make([][]bool, gridSize)
	for i := range gridSize {
		visited[i] = make([]bool, gridSize)
	}

	pathLen := 0
	solved := false
	q := queue.New()
	q.Enqueue(gridPos{0, 0})

outer:
	for q.Len() > 0 {
		currentQueueLen := q.Len()
		for range currentQueueLen {
			pos := q.Dequeue().(gridPos)
			if pos.row == gridSize-1 && pos.col == gridSize-1 {
				// Reached the end
				solved = true
				break outer
			}
			visited[pos.row][pos.col] = true
			adjs := pos.adjacencies(gridSize, gridSize)
			for _, adj := range adjs {
				if walls[adj.row][adj.col] > time && !visited[adj.row][adj.col] {
					// Open space that we've not visited yet
					visited[adj.row][adj.col] = true
					q.Enqueue(adj)
				}
			}
		}
		pathLen++
	}

	if !solved {
		pathLen = -1
	}

	return pathLen
}

func Day18Part2(logger *slog.Logger, input string, part1Context any) string {
	context := part1Context.(d18Context)
	maxTime := len(context.lines)

	// Binary search the remaining ticks.
	for {
		time := context.minTime + ((maxTime - context.minTime) / 2)

		canExit := runMaze(context.walls, time, context.gridSize) >= 0

		if canExit {
			// We're not blocked yet - try to the right
			context.minTime = time + 1
			if context.minTime == maxTime {
				// Finished - the one to the right is the answer.
				return context.lines[time+1]
			}
		} else {
			// We're blocked - try to the left.
			maxTime = time
			if context.minTime == maxTime {
				// Finished - the one to the left is the answer.
				return context.lines[time-1]
			}
		}
	}
}
