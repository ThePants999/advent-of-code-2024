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

type WallSetStatus int

const (
	UNCONNECTED WallSetStatus = iota
	TOP_RIGHT
	BOTTOM_LEFT
)

type d18GridSquare struct {
	wall              bool // Remaining elements only valid if this is true
	disjointSetParent *d18GridSquare
	disjointSetRank   int
	status            WallSetStatus
}

type d18MergeResult struct {
	newRoot    *d18GridSquare
	compatible bool
}

type d18Context struct {
	grid     [][]d18GridSquare
	gridSize int
	lines    []string
}

func (p gridPos) allAdjacencies(numRows int, numCols int) []gridPos {
	adj := make([]gridPos, 0, 8)
	if p.row > 0 {
		adj = append(adj, gridPos{p.row - 1, p.col})
		if p.col > 0 {
			adj = append(adj, gridPos{p.row - 1, p.col - 1})
		}
		if p.col < numCols-1 {
			adj = append(adj, gridPos{p.row - 1, p.col + 1})
		}
	}
	if p.row < numRows-1 {
		adj = append(adj, gridPos{p.row + 1, p.col})
		if p.col > 0 {
			adj = append(adj, gridPos{p.row + 1, p.col - 1})
		}
		if p.col < numCols-1 {
			adj = append(adj, gridPos{p.row + 1, p.col + 1})
		}
	}
	if p.col > 0 {
		adj = append(adj, gridPos{p.row, p.col - 1})
	}
	if p.col < numCols-1 {
		adj = append(adj, gridPos{p.row, p.col + 1})
	}
	return adj
}

func findDisjointSetRoot(square *d18GridSquare) *d18GridSquare {
	if square.disjointSetParent != nil {
		// We take this opportunity to flatten the tree by replacing
		// all parent pointers along the path with the root pointer.
		square.disjointSetParent = findDisjointSetRoot(square.disjointSetParent)
		return square.disjointSetParent
	} else {
		return square
	}
}

func mergeDisjointSets(squareA *d18GridSquare, squareB *d18GridSquare) d18MergeResult {
	rootA := findDisjointSetRoot(squareA)
	rootB := findDisjointSetRoot(squareB)

	if rootA == rootB {
		// Already in the same set
		return d18MergeResult{rootA, true}
	}

	// Make "A" whichever is the highest-ranked of the two
	if rootA.disjointSetRank < rootB.disjointSetRank {
		rootA, rootB = rootB, rootA
	} else if rootA.disjointSetRank == rootB.disjointSetRank {
		// A's rank will increase through adding B as a child
		rootA.disjointSetRank++
	}

	// Move B under A
	rootB.disjointSetParent = rootA

	if rootA.status == UNCONNECTED {
		// A being unconnected means they're definitely compatible
		// but B's status should have priority.
		rootA.status = rootB.status
		return d18MergeResult{rootA, true}
	}
	if rootB.status == UNCONNECTED {
		// B being unconnnected while A is connected means they're
		// compatible but the status is already correctly recorded.
		// (Moving B under A means B gains A's status.)
		return d18MergeResult{rootA, true}
	}

	// Both sets are attached to grid edges, so whether they're
	// compatible depends on whether they're attached to the same
	// one.
	return d18MergeResult{rootA, rootA.status == rootB.status}
}

// The approach this code takes is to maintain groups of walls that are
// connected to each other (including diagonally) as "disjoint sets"
// (see https://en.wikipedia.org/wiki/Disjoint-set_data_structure). Each
// set knows whether it's touching the bottom-left of the grid (the left
// edge or the bottom edge), or the top-right of the grid, or whether it's
// floating in the middle. Newly-added walls can join sets together,
// which might join a floating set with one connected to an edge, thereby
// making the new combined set a connected one. The moment we join a
// bottom-left-connected set with a top-right-connected set, that's the
// moment we make it impossible to traverse from top left to bottom right.
// (This function returns true when that happens.)
func createWall(grid [][]d18GridSquare, wallStr string, gridSize int) bool {
	commaIx := strings.IndexRune(wallStr, ',')
	col, _ := strconv.Atoi(wallStr[:commaIx])
	row, _ := strconv.Atoi(wallStr[commaIx+1:])
	square := &grid[row][col]
	square.wall = true

	if row == 0 || col == gridSize-1 {
		square.status = TOP_RIGHT
	} else if row == gridSize-1 || col == 0 {
		square.status = BOTTOM_LEFT
	}

	adjs := gridPos{row, col}.allAdjacencies(gridSize, gridSize)
	for _, adj := range adjs {
		otherSquare := &grid[adj.row][adj.col]
		if otherSquare.wall {
			// There's a wall here, we'll need to merge sets.
			result := mergeDisjointSets(square, otherSquare)
			square = result.newRoot
			if !result.compatible {
				return true
			}
		}
	}

	return false
}

func Day18Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	gridSize := GRID_SIZE_EXAMPLE
	startAfter := START_AFTER_EXAMPLE
	if len(lines) > 100 {
		gridSize = GRID_SIZE_REAL
		startAfter = START_AFTER_REAL
	}

	grid := make([][]d18GridSquare, gridSize)
	for rowIx := range gridSize {
		grid[rowIx] = make([]d18GridSquare, gridSize)
	}

	for lineIx := range startAfter {
		_ = createWall(grid, lines[lineIx], gridSize)
	}

	pathLen := runMaze(grid, gridSize)
	return strconv.Itoa(pathLen), d18Context{grid, gridSize, lines[startAfter:]}
}

// Just a basic BFS that returns the length of the best path.
func runMaze(grid [][]d18GridSquare, gridSize int) int {
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
				if !grid[adj.row][adj.col].wall && !visited[adj.row][adj.col] {
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
	for _, line := range context.lines {
		// See the comment above createWall() for an explanation of the
		// algorithm we use in this part.
		if createWall(context.grid, line, context.gridSize) {
			return line
		}
	}
	return ""
}
