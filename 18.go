package main

import (
	"log/slog"
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
	walls       [][]bool
	gridSize    int
	lines       []string
	currentPath map[gridPos]nothing
}

type d18State struct {
	pos  gridPos
	prev *d18State
}

func Day18Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	gridSize := GRID_SIZE_EXAMPLE
	startAfter := START_AFTER_EXAMPLE
	if len(lines) > 100 {
		gridSize = GRID_SIZE_REAL
		startAfter = START_AFTER_REAL
	}
	walls := make([][]bool, gridSize)
	visited := make([][]bool, gridSize)
	for rowIx := range gridSize {
		walls[rowIx] = make([]bool, gridSize)
		visited[rowIx] = make([]bool, gridSize)
	}
	for lineIx, line := range lines {
		if lineIx == startAfter {
			break
		}
		commaIx := strings.IndexRune(line, ',')
		colIx, _ := strconv.Atoi(line[:commaIx])
		rowIx, _ := strconv.Atoi(line[commaIx+1:])
		walls[rowIx][colIx] = true
	}
	lines = lines[startAfter:]

	path := runMaze(walls, gridSize)
	return strconv.Itoa(len(path) - 1), d18Context{walls, gridSize, lines, path}
}

func runMaze(walls [][]bool, gridSize int) map[gridPos]nothing {
	visited := make([][]bool, gridSize)
	for i := range gridSize {
		visited[i] = make([]bool, gridSize)
	}

	var finalState *d18State
	q := queue.New()
	q.Enqueue(d18State{gridPos{0, 0}, nil})
	for q.Len() > 0 {
		s := q.Dequeue().(d18State)
		if s.pos.row == gridSize-1 && s.pos.col == gridSize-1 {
			// Reached the end
			finalState = &s
			break
		}
		if visited[s.pos.row][s.pos.col] {
			// Been here before
			continue
		}
		visited[s.pos.row][s.pos.col] = true
		adjs := s.pos.adjacencies(gridSize, gridSize)
		for _, adj := range adjs {
			if !walls[adj.row][adj.col] {
				q.Enqueue(d18State{adj, &s})
			}
		}
	}

	if finalState != nil {
		path := make(map[gridPos]nothing)
		for s := finalState; s != nil; s = s.prev {
			path[s.pos] = nothing{}
		}
		return path
	} else {
		return nil
	}
}

func Day18Part2(logger *slog.Logger, input string, part1Context any) string {
	context := part1Context.(d18Context)

	path := context.currentPath
	var line string
	for _, line = range context.lines {
		commaIx := strings.IndexRune(line, ',')
		colIx, _ := strconv.Atoi(line[:commaIx])
		rowIx, _ := strconv.Atoi(line[commaIx+1:])
		context.walls[rowIx][colIx] = true

		_, inPath := path[gridPos{rowIx, colIx}]
		if inPath {
			// We just obstructed the current path, find a new one.
			path = runMaze(context.walls, context.gridSize)
			if path == nil {
				break
			}
		}
	}

	return line
}
