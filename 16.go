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
	id     d16NodeId
	cost   int
	heap   *d16Heap
	heapIx int
}

// Hand-rolled heap implementation, giving us a priority queue
// we can use to run Dijkstra's algorithm.

type d16Heap struct {
	arr []*d16Node
}

func NewHeap() *d16Heap {
	heap := d16Heap{}
	heap.arr = make([]*d16Node, 1, 10000)
	return &heap
}

func (heap *d16Heap) Push(new *d16Node) {
	heap.arr = append(heap.arr, new)
	new.heap = heap
	new.heapIx = len(heap.arr) - 1
	heap.PercolateUp(new.heapIx)
}

func (heap *d16Heap) PercolateUp(ix int) {
	for ; ix > 1 && heap.arr[ix].cost < heap.arr[ix/2].cost; ix /= 2 {
		heap.arr[ix], heap.arr[ix/2] = heap.arr[ix/2], heap.arr[ix]
		heap.arr[ix].heapIx = ix
		heap.arr[ix/2].heapIx = ix / 2
	}
}

func (heap *d16Heap) PercolateDown(ix int) {
	for 2*ix < len(heap.arr) {
		leftIx := ix * 2
		rightIx := leftIx + 1
		if rightIx < len(heap.arr) && heap.arr[rightIx].cost < heap.arr[leftIx].cost && heap.arr[rightIx].cost < heap.arr[ix].cost {
			heap.arr[ix], heap.arr[rightIx] = heap.arr[rightIx], heap.arr[ix]
			heap.arr[ix].heapIx = ix
			heap.arr[rightIx].heapIx = rightIx
			ix = rightIx
		} else if heap.arr[leftIx].cost < heap.arr[ix].cost {
			heap.arr[ix], heap.arr[leftIx] = heap.arr[leftIx], heap.arr[ix]
			heap.arr[ix].heapIx = ix
			heap.arr[leftIx].heapIx = leftIx
			ix = leftIx
		} else {
			break
		}
	}
}

func (heap *d16Heap) Pop() *d16Node {
	if len(heap.arr) == 1 {
		return nil
	}
	ret := heap.arr[1]
	ret.heapIx = 0
	ret.heap = nil
	if len(heap.arr) == 2 {
		heap.arr = heap.arr[:1]
		return ret
	}

	// Swap last element to root.
	heap.arr[1] = heap.arr[len(heap.arr)-1]
	heap.arr[1].heapIx = 1
	heap.arr = heap.arr[:len(heap.arr)-1]

	// Percolate down to the correct place.
	heap.PercolateDown(1)

	return ret
}

func (heap *d16Heap) Update(node *d16Node) {
	ix := node.heapIx
	if ix > 1 && heap.arr[ix].cost < heap.arr[ix/2].cost {
		heap.PercolateUp(ix)
	} else {
		heap.PercolateDown(ix)
	}
}

func (heap *d16Heap) IsEmpty() bool {
	return len(heap.arr) == 1
}

func (node *d16Node) UpdateCost(newCost int) {
	node.cost = newCost
	if node.heap != nil {
		node.heap.Update(node)
	}
}

type d16context struct {
	grid     [][]d16GridSquare
	startRow int
	startCol int
	endRow   int
	endCol   int
}

func updateFrom(grid [][]d16GridSquare, node *d16Node) {
	// Per Dijkstra's algorithm, update the costs of nodes reachable from this one
	// (which are: moving forwards, turning left, turning right). However, we don't
	// bother to assign costs to turning left and right if they would face a wall
	// as we trivially know that's never going to be optimal.
	frontPos := node.id.pos.move(node.id.dir)
	leftDir := node.id.dir.turn(false)
	leftPos := node.id.pos.move(leftDir)
	rightDir := node.id.dir.turn(true)
	rightPos := node.id.pos.move(rightDir)
	if !grid[leftPos.row][leftPos.col].wall && node.cost+1000 < grid[node.id.pos.row][node.id.pos.col].nodes[leftDir].cost {
		grid[node.id.pos.row][node.id.pos.col].nodes[leftDir].UpdateCost(node.cost + 1000)
	}
	if !grid[rightPos.row][rightPos.col].wall && node.cost+1000 < grid[node.id.pos.row][node.id.pos.col].nodes[rightDir].cost {
		grid[node.id.pos.row][node.id.pos.col].nodes[rightDir].UpdateCost(node.cost + 1000)
	}
	if node.cost+1 < grid[frontPos.row][frontPos.col].nodes[node.id.dir].cost {
		// Node that there's no need to check to see whether there's a wall in front
		// of us - walls nodes never get initialised with an infinite cost, so they'll
		// always appear to be 0-cost already, but they're also not on the heap.
		grid[frontPos.row][frontPos.col].nodes[node.id.dir].UpdateCost(node.cost + 1)
	}
}

func Day16Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	grid := make([][]d16GridSquare, len(lines))
	startRow, startCol := 0, 0
	endRow, endCol := 0, 0
	heap := NewHeap()
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

			// Every space that's not a wall needs to be a graph node
			// (well, technically, four graph nodes, one per direction)
			// and to be in the priority queue for the upcoming
			// Dijkstra.
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
	grid[startRow][startCol].nodes[D6_RIGHT].UpdateCost(0)
	for !heap.IsEmpty() {
		node := heap.Pop()
		if node.id.pos.row == endRow && node.id.pos.col == endCol {
			// Once we've popped the end, no other nodes can be on a best
			// path.
			break
		}
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
