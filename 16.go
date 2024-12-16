package main

import (
	"log/slog"
	"math"
	"strconv"
	"strings"

	stack "github.com/golang-collections/collections/stack"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day16 = runner.DayImplementation{
	DayNumber:    16,
	ExecutePart1: Day16Part1,
	ExecutePart2: Day16Part2,
	ExampleInput: `###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
###############`,
	ExamplePart1Answer: "7036",
	ExamplePart2Answer: "45",
}

func (dir direction6) turn(clockwise bool) direction6 {
	newDir := dir + 1
	if !clockwise {
		newDir = dir - 1
	}

	if newDir > D6_LEFT {
		newDir = D6_UP
	} else if newDir < 0 {
		newDir = D6_LEFT
	}

	return newDir
}

func (pos gridPos) move(dir direction6) gridPos {
	switch dir {
	case D6_UP:
		pos.row--
	case D6_DOWN:
		pos.row++
	case D6_LEFT:
		pos.col--
	case D6_RIGHT:
		pos.col++
	}
	return pos
}

type d16GridSquare struct {
	wall bool
	cost [4]int
}

type d16candidate struct {
	pos  gridPos
	dir  direction6
	cost int
	next *d16candidate
}

type d16candidate2 struct {
	pos  gridPos
	dir  direction6
	path []gridPos
}

type d16heap struct {
	first *d16candidate
}

func (heap *d16heap) Push(new *d16candidate) {
	pos := &(heap.first)
	for *pos != nil && (**pos).cost < new.cost {
		pos = &((**pos).next)
	}
	new.next = *pos
	*pos = new
}

func (heap *d16heap) Pop() *d16candidate {
	ret := heap.first
	heap.first = heap.first.next
	return ret
}

type d16context struct {
	grid     [][]d16GridSquare
	startRow int
	startCol int
	endRow   int
	endCol   int
}

func Day16Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	grid := make([][]d16GridSquare, len(lines))
	startRow, startCol := 0, 0
	endRow, endCol := 0, 0
	for rowIx, row := range lines {
		grid[rowIx] = make([]d16GridSquare, len(row))
		for colIx, square := range row {
			grid[rowIx][colIx].cost[D6_UP] = -1
			grid[rowIx][colIx].cost[D6_DOWN] = -1
			grid[rowIx][colIx].cost[D6_LEFT] = -1
			grid[rowIx][colIx].cost[D6_RIGHT] = -1
			switch square {
			case '#':
				grid[rowIx][colIx].wall = true
			case 'S':
				startRow, startCol = rowIx, colIx
			case 'E':
				endRow, endCol = rowIx, colIx
			}
		}
	}

	heap := d16heap{}
	heap.Push(&d16candidate{gridPos{startRow, startCol}, D6_RIGHT, 0, nil})
	bestCost := math.MaxInt

	for heap.first != nil {
		candidate := heap.Pop()
		if grid[candidate.pos.row][candidate.pos.col].cost[candidate.dir] > -1 {
			// Been here before
			continue
		}

		grid[candidate.pos.row][candidate.pos.col].cost[candidate.dir] = candidate.cost

		if candidate.pos.row == endRow && candidate.pos.col == endCol {
			if candidate.cost < bestCost {
				bestCost = candidate.cost
			} else if candidate.cost > bestCost {
				// We don't record the cost on the destination square if
				// it's higher than reaching the square from another
				// direction, else we'll consider paths in this direction
				// to be valid in part 2.
				grid[candidate.pos.row][candidate.pos.col].cost[candidate.dir] = -1
			}

			continue
		}

		grid[candidate.pos.row][candidate.pos.col].cost[candidate.dir] = candidate.cost

		frontPos := candidate.pos.move(candidate.dir)
		leftDir := candidate.dir.turn(false)
		rightDir := candidate.dir.turn(true)
		heap.Push(&d16candidate{candidate.pos, leftDir, candidate.cost + 1000, nil})
		heap.Push(&d16candidate{candidate.pos, rightDir, candidate.cost + 1000, nil})
		if !grid[frontPos.row][frontPos.col].wall {
			heap.Push(&d16candidate{frontPos, candidate.dir, candidate.cost + 1, nil})
		}
	}

	return strconv.Itoa(bestCost), d16context{grid, startRow, startCol, endRow, endCol}
}

func Day16Part2(logger *slog.Logger, input string, part1Context any) string {
	context := part1Context.(d16context)
	bestPaths := make(map[gridPos]nothing)
	s := stack.New()
	s.Push(d16candidate2{gridPos{context.startRow, context.startCol}, D6_RIGHT, make([]gridPos, 0, 1000)})

	for s.Len() > 0 {
		candidate := s.Pop().(d16candidate2)

		if candidate.pos.row == context.endRow && candidate.pos.col == context.endCol {
			// We're done with this path
			for _, pos := range candidate.path {
				bestPaths[pos] = nothing{}
			}
			continue
		}

		newPath := append(candidate.path, candidate.pos)
		frontPos := candidate.pos.move(candidate.dir)
		leftDir := candidate.dir.turn(false)
		rightDir := candidate.dir.turn(true)
		if context.grid[candidate.pos.row][candidate.pos.col].cost[leftDir] == context.grid[candidate.pos.row][candidate.pos.col].cost[candidate.dir]+1000 {
			s.Push(d16candidate2{candidate.pos, leftDir, newPath})
		}
		if context.grid[candidate.pos.row][candidate.pos.col].cost[rightDir] == context.grid[candidate.pos.row][candidate.pos.col].cost[candidate.dir]+1000 {
			s.Push(d16candidate2{candidate.pos, rightDir, newPath})
		}
		if context.grid[frontPos.row][frontPos.col].cost[candidate.dir] == context.grid[candidate.pos.row][candidate.pos.col].cost[candidate.dir]+1 {
			s.Push(d16candidate2{frontPos, candidate.dir, newPath})
		}
	}

	// fmt.Println()
	// for rowIx, row := range context.grid {
	// 	for colIx, square := range row {
	// 		_, best := bestPaths[gridPos{rowIx, colIx}]
	// 		if best {
	// 			fmt.Print("O")
	// 		} else if square.wall {
	// 			fmt.Print("#")
	// 		} else {
	// 			fmt.Print(".")
	// 		}
	// 	}
	// 	fmt.Println()
	// }

	return strconv.Itoa(len(bestPaths) + 1)
}
