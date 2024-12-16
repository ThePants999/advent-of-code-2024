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
	wall  bool
	nodes [4]d16Node
}

type d16NodeId struct {
	pos gridPos
	dir direction6
}

type d16Node struct {
	id   d16NodeId
	cost int
	heap *d16heap
	next *d16Node
	prev *d16Node
}

type d16heap struct {
	first *d16Node
	last  *d16Node
}

func (heap *d16heap) Push(new *d16Node) {
	new.heap = heap
	if heap.first == nil {
		heap.first = new
		heap.last = new
	} else {
		var pos *d16Node
		for pos = heap.first; pos != nil && pos.cost < new.cost; pos = pos.next {
		}

		if pos == nil {
			// At end
			heap.last.next = new
			new.prev = heap.last
			heap.last = new
			new.next = nil
		} else if pos == heap.first {
			// At front
			heap.first.prev = new
			new.next = heap.first
			heap.first = new
			new.prev = nil
		} else {
			new.prev = pos.prev
			new.next = pos
			pos.prev.next = new
			pos.prev = new
		}
	}
}

func (heap *d16heap) Pop() *d16Node {
	ret := heap.first
	if ret.next != nil {
		ret.next.prev = ret.prev
	}
	heap.first = ret.next
	if heap.first == nil {
		heap.last = nil
	}
	return ret
}

func (heap *d16heap) Remove(node *d16Node) {
	if node == heap.first {
		heap.first = node.next
	}
	if node == heap.last {
		heap.last = node.prev
	}
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	node.heap = nil
}

func (node *d16Node) UpdateCost(newCost int) {
	node.cost = newCost
	if node.heap != nil {
		node.heap.Update(node)
	}
}

func (heap *d16heap) Update(node *d16Node) {
	heap.Remove(node)
	heap.Push(node)
}

func (heap *d16heap) IsEmpty() bool {
	return heap.first == nil
}

type d16context struct {
	grid     [][]d16GridSquare
	startRow int
	startCol int
	endRow   int
	endCol   int
}

func updateFrom(grid [][]d16GridSquare, node *d16Node) {
	frontPos := node.id.pos.move(node.id.dir)
	leftDir := node.id.dir.turn(false)
	rightDir := node.id.dir.turn(true)
	if node.cost+1000 < grid[node.id.pos.row][node.id.pos.col].nodes[leftDir].cost {
		grid[node.id.pos.row][node.id.pos.col].nodes[leftDir].UpdateCost(node.cost + 1000)
	}
	if node.cost+1000 < grid[node.id.pos.row][node.id.pos.col].nodes[rightDir].cost {
		grid[node.id.pos.row][node.id.pos.col].nodes[rightDir].UpdateCost(node.cost + 1000)
	}
	if node.cost+1 < grid[frontPos.row][frontPos.col].nodes[node.id.dir].cost {
		grid[frontPos.row][frontPos.col].nodes[node.id.dir].UpdateCost(node.cost + 1)
	}
}

func Day16Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	grid := make([][]d16GridSquare, len(lines))
	startRow, startCol := 0, 0
	endRow, endCol := 0, 0
	heap := d16heap{}
	for rowIx, row := range lines {
		grid[rowIx] = make([]d16GridSquare, len(row))
		for colIx, square := range row {
			switch square {
			case '#':
				grid[rowIx][colIx].wall = true
			case 'S':
				startRow, startCol = rowIx, colIx
			case 'E':
				endRow, endCol = rowIx, colIx
			}

			if !grid[rowIx][colIx].wall {
				for dir := D6_UP; dir <= D6_LEFT; dir++ {
					grid[rowIx][colIx].nodes[dir].id.pos = gridPos{rowIx, colIx}
					grid[rowIx][colIx].nodes[dir].id.dir = dir
					grid[rowIx][colIx].nodes[dir].cost = math.MaxInt
					heap.Push(&grid[rowIx][colIx].nodes[dir])
				}
			}
		}
	}

	// Run Dijkstra's algorithm over the maze, treating each combination
	// of grid position and facing as a different node in the graph.
	// We run it fully, rather than stopping as soon as we reach the
	// end, to ensure we've found all the best routes, as needed for part 2.
	grid[startRow][startCol].nodes[D6_RIGHT].UpdateCost(0)
	for !heap.IsEmpty() {
		node := heap.Pop()
		updateFrom(grid, node)
	}

	// Bit of cheekiness - we know the end is in the top-right
	// corner, let's assume that all best paths end in the up
	// and/or right directions, so we just need to see which of
	// those has the lowest cost.
	bestCost := grid[endRow][endCol].nodes[D6_UP].cost
	if grid[endRow][endCol].nodes[D6_RIGHT].cost < bestCost {
		bestCost = grid[endRow][endCol].nodes[D6_RIGHT].cost
	}

	return strconv.Itoa(bestCost), d16context{grid, startRow, startCol, endRow, endCol}
}

func Day16Part2(logger *slog.Logger, input string, part1Context any) string {
	context := part1Context.(d16context)
	s := stack.New()
	visited := make(map[d16NodeId]nothing)
	// Repeat of aforementioned cheekiness.
	s.Push(&context.grid[context.endRow][context.endCol].nodes[D6_RIGHT])
	s.Push(&context.grid[context.endRow][context.endCol].nodes[D6_UP])

	// What we're going to do now is run a simple DFS, but we're going to
	// do it in reverse from the end, and we're only going to allow
	// traversal to nodes whose cost as recorded in the earlier Dijkstra
	// run match the current node's cost minus the cost to reach them from
	// the current node. That way, we know we're always sticking to best
	// paths, so we simply need to count the number of grid squares we end
	// up visiting doing this.
	for s.Len() > 0 {
		node := s.Pop().(*d16Node)

		_, alreadyVisited := visited[node.id]
		if alreadyVisited {
			continue
		}
		visited[node.id] = nothing{}

		thisGrid := &context.grid[node.id.pos.row][node.id.pos.col]
		backPos := node.id.pos.move(node.id.dir.turn(false).turn(false))
		backNode := &context.grid[backPos.row][backPos.col].nodes[node.id.dir]
		leftDir := node.id.dir.turn(false)
		leftNode := &thisGrid.nodes[leftDir]
		rightDir := node.id.dir.turn(true)
		rightNode := &thisGrid.nodes[rightDir]
		if leftNode.cost == node.cost-1000 {
			s.Push(leftNode)
		}
		if rightNode.cost == node.cost-1000 {
			s.Push(rightNode)
		}
		if backNode.cost == node.cost-1 {
			s.Push(backNode)
		}
	}

	bestPaths := make(map[gridPos]nothing)
	for node := range visited {
		bestPaths[node.pos] = nothing{}
	}

	return strconv.Itoa(len(bestPaths))
}
