package main

import (
	"log/slog"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day4 = runner.DayImplementation{
	DayNumber:          4,
	ExecutePart1:       Day4Part1,
	ExecutePart2:       Day4Part2,
	ExampleInput:       "MMMSXXMASM\nMSAMXMSMSA\nAMXSXMAAMM\nMSAMASMSMX\nXMASAMXAMM\nXXAMMXXAMA\nSMSMSASXSS\nSAXAMASAAA\nMAMMMXMMMM\nMXMXAXMASX",
	ExamplePart1Answer: "18",
	ExamplePart2Answer: "9",
}

type coords struct {
	row int
	col int
}

type Direction int

const (
	UP Direction = iota
	UP_RIGHT
	RIGHT
	DOWN_RIGHT
	DOWN
	DOWN_LEFT
	LEFT
	UP_LEFT
)

var nextLetter = map[rune]rune{
	'X': 'M',
	'M': 'A',
	'A': 'S',
}

func Day4Part1(logger *slog.Logger, input string) (string, any) {
	// Parse the input, building up both the
	// full grid and also a list of Xs.
	inputRows := strings.Fields(input)
	grid := make([][]rune, len(inputRows))
	xs := make([]coords, 0, len(input))
	for rowIx, inputRow := range inputRows {
		grid[rowIx] = make([]rune, len(inputRow))
		for colIx, inputRune := range inputRow {
			grid[rowIx][colIx] = inputRune
			if inputRune == 'X' {
				xs = append(xs, coords{rowIx, colIx})
			}
		}
	}

	// Go through each X, and for each one, go through
	// each of the 8 directions, looking to see if there's
	// an M, then an A, then an S in that direction.
	sum := 0
	for _, x := range xs {
		for dir := UP; dir <= UP_LEFT; dir++ {
			if testDirection(grid, x.row, x.col, dir) {
				sum++
			}
		}
	}
	return strconv.Itoa(sum), grid
}

func testDirection(grid [][]rune, row int, col int, dir Direction) bool {
	// Determine what letter we're expecting next, based
	// on what letter we've got here.
	expectedNextLetter := nextLetter[grid[row][col]]
	if expectedNextLetter == 0 {
		return true
	}

	// Now determine where "next" is - the coordinates of the
	// next square along in the specified direction.
	newRow := row
	if dir == UP || dir == UP_LEFT || dir == UP_RIGHT {
		newRow -= 1
		if newRow < 0 {
			return false
		}
	} else if dir == DOWN || dir == DOWN_LEFT || dir == DOWN_RIGHT {
		newRow += 1
		if newRow == len(grid) {
			return false
		}
	}
	newCol := col
	if dir == LEFT || dir == UP_LEFT || dir == DOWN_LEFT {
		newCol -= 1
		if newCol < 0 {
			return false
		}
	} else if dir == RIGHT || dir == UP_RIGHT || dir == DOWN_RIGHT {
		newCol += 1
		if newCol == len(grid[newRow]) {
			return false
		}
	}

	if expectedNextLetter == grid[newRow][newCol] {
		// This letter was correct - recurse to check
		// the next one.
		return testDirection(grid, newRow, newCol, dir)
	} else {
		return false
	}
}

// Determine whether a co-ordinate already determined to contain
// an A is the center of two diagonal MASes.
func testMAS(grid [][]rune, row int, col int) bool {
	if ((grid[row-1][col-1] == 'M' && grid[row+1][col+1] == 'S') || (grid[row-1][col-1] == 'S' && grid[row+1][col+1] == 'M')) &&
		((grid[row+1][col-1] == 'M' && grid[row-1][col+1] == 'S') || (grid[row+1][col-1] == 'S' && grid[row-1][col+1] == 'M')) {
		return true
	}
	return false
}

func Day4Part2(logger *slog.Logger, input string, part1Context any) string {
	// You'd think that we'd have built a list of As during
	// initial grid construction. Nope, didn't bother, we'll
	// go through the whole grid looking for them. It runs in
	// less than 100 microseconds anyway.
	grid := part1Context.([][]rune)
	sum := 0
	for rowIx := 1; rowIx < len(grid)-1; rowIx++ {
		for colIx := 1; colIx < len(grid[rowIx])-1; colIx++ {
			if grid[rowIx][colIx] == 'A' && testMAS(grid, rowIx, colIx) {
				sum++
			}
		}
	}
	return strconv.Itoa(sum)
}
