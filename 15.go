package main

import (
	"log/slog"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day15 = runner.DayImplementation{
	DayNumber:    15,
	ExecutePart1: Day15Part1,
	ExecutePart2: Day15Part2,
	ExampleInput: `##########
#..O..O.O#
#......O.#
#.OO..O.O#
#..O@..O.#
#O#..O...#
#O..O..O.#
#.OO.O.OO#
#....O...#
##########

<vv>^<v^>v>^vv^v>v<>v^v<v<^vv<<<^><<><>>v<vvv<>^v^>^<<<><<v<<<v^vv^v>^
vvv<<^>^v^^><<>>><>^<<><^vv^^<>vvv<>><^^v>^>vv<>v<<<<v<^v>^<^^>>>^<v<v
><>vv>v^v^<>><>>>><^^>vv>v<^^^>>v^v^<^^>v^^>v^<^v>v<>>v^v^<v>v^^<^^vv<
<<v<^>>^^^^>>>v^<>vvv^><v<<<>^^^vv^<vvv>^>v<^^^^v<>^>vvvv><>>v^<<^^^^^
^><^><>>><>^^<<^^v>>><^<v>^<vv>>v>>>^v><>^v><<<<v>>v<v<v>vvv>^<><<>^><
^>><>^v<><^vvv<^^<><v<<<<<><^v<<<><<<^^<v<^^^><^>>^<v^><<<^>>^v<v^v<v^
>^>>^v>vv>^<<^v<>><<><<v<<v><>v<^vv<<<>^^v^>^^>>><<^v>>v^v><^^>>^<>vv^
<><^^>^^^<><vvvvv^v<v<<>^v<v>v<<^><<><<><<<^^<<<^<<>><<><^^^>^^<>^>v<>
^^>vv<^v^v<vv>^<><v<^v>^^^>>>^^vvv^>vvv<>>>^<^>>>>>^<<^v>^vvv<>^<><<v>
v^^>>><<^^<>>^v^<v^vv<>v^<<>^<^v^v><^<<<><<^<v><v<>vv>>v><v^<vv<>v^<<^`,
	ExamplePart1Answer: "10092",
	ExamplePart2Answer: "9021",
}

// There isn't much to say about today, so I'm not going to
// thoroughly comment throughout. We're just genuinely simulating
// robot movement and pushing boxes around, it's plenty fast enough.
// The challenge of the day is just being sure to get the logic
// right.

type gridContents int

const (
	EMPTY gridContents = iota
	WALL
	CRATE
	CRATE_LEFT
	CRATE_RIGHT
)

func Day15Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	grid := make([][]gridContents, len(lines))
	robotRow, robotCol := 0, 0
	var sb strings.Builder
	gridFinished := false
	for rowIx, line := range lines {
		if !gridFinished {
			if line[0] != '#' {
				gridFinished = true
				grid = grid[:rowIx]
			} else {
				grid[rowIx] = make([]gridContents, len(line))
				for colIx, char := range line {
					switch char {
					case '#':
						grid[rowIx][colIx] = WALL
					case 'O':
						grid[rowIx][colIx] = CRATE
					case '@':
						robotRow, robotCol = rowIx, colIx
					}
				}
			}
		}

		if gridFinished {
			sb.WriteString(line)
		}
	}

	movements := sb.String()
	for _, char := range movements {
		curRow, curCol := robotRow, robotCol
		var delta int
		var changingVar, changingRobotVar *int
		switch char {
		case '^':
			delta = -1
			changingVar = &curRow
			changingRobotVar = &robotRow
		case 'v':
			delta = 1
			changingVar = &curRow
			changingRobotVar = &robotRow
		case '>':
			delta = 1
			changingVar = &curCol
			changingRobotVar = &robotCol
		case '<':
			delta = -1
			changingVar = &curCol
			changingRobotVar = &robotCol
		}

		canMove := true
		for {
			*changingVar += delta
			if grid[curRow][curCol] == WALL {
				canMove = false
				break
			} else if grid[curRow][curCol] == EMPTY {
				break
			}
		}
		if canMove {
			*changingRobotVar += delta
			grid[robotRow][robotCol] = EMPTY
			if robotRow != curRow || robotCol != curCol {
				grid[curRow][curCol] = CRATE
			}
		}
	}

	sum := 0
	for rowIx, row := range grid {
		for colIx, space := range row {
			if space == CRATE {
				sum += (rowIx * 100) + colIx
			}
		}
	}
	return strconv.Itoa(sum), nil
}

func Day15Part2(logger *slog.Logger, input string, part1Context any) string {
	lines := strings.Fields(input)
	grid := make([][]gridContents, len(lines))
	robotRow, robotCol := 0, 0
	var sb strings.Builder
	gridFinished := false
	for rowIx, line := range lines {
		if !gridFinished {
			if line[0] != '#' {
				gridFinished = true
				grid = grid[:rowIx]
			} else {
				grid[rowIx] = make([]gridContents, len(line)*2)
				for colIx, char := range line {
					switch char {
					case '#':
						grid[rowIx][colIx*2] = WALL
						grid[rowIx][colIx*2+1] = WALL
					case 'O':
						grid[rowIx][colIx*2] = CRATE_LEFT
						grid[rowIx][colIx*2+1] = CRATE_RIGHT
					case '@':
						robotRow, robotCol = rowIx, colIx*2
					}
				}
			}
		}

		if gridFinished {
			sb.WriteString(line)
		}
	}

	movements := sb.String()
	for _, char := range movements {
		var dir Direction
		targetRow, targetCol := robotRow, robotCol
		switch char {
		case '^':
			dir = UP
			targetRow--
		case 'v':
			dir = DOWN
			targetRow++
		case '>':
			dir = RIGHT
			targetCol++
		case '<':
			dir = LEFT
			targetCol--
		}

		if canMove(grid, targetRow, targetCol, dir) {
			doMove(grid, targetRow, targetCol, dir, EMPTY)
			robotRow = targetRow
			robotCol = targetCol
		}
	}

	sum := 0
	for rowIx, row := range grid {
		for colIx, space := range row {
			if space == CRATE_LEFT {
				sum += (rowIx * 100) + colIx
			}
		}
	}
	return strconv.Itoa(sum)
}

func doMove(grid [][]gridContents, row int, col int, dir Direction, incoming gridContents) {
	if grid[row][col] == WALL {
		panic("Hit a wall")
	}
	if grid[row][col] == CRATE {
		panic("Found a part 1 crate!")
	}
	if grid[row][col] != EMPTY {
		rowDelta := 1
		switch dir {
		case LEFT:
			doMove(grid, row, col-1, dir, grid[row][col])
		case RIGHT:
			doMove(grid, row, col+1, dir, grid[row][col])
		case UP:
			rowDelta = -1
			fallthrough
		case DOWN:
			if grid[row][col] == CRATE_LEFT {
				grid[row][col] = EMPTY
				grid[row][col+1] = EMPTY
				doMove(grid, row+rowDelta, col, dir, CRATE_LEFT)
				doMove(grid, row+rowDelta, col+1, dir, CRATE_RIGHT)
			} else {
				grid[row][col] = EMPTY
				grid[row][col-1] = EMPTY
				doMove(grid, row+rowDelta, col-1, dir, CRATE_LEFT)
				doMove(grid, row+rowDelta, col, dir, CRATE_RIGHT)
			}
		}
	}

	grid[row][col] = incoming
}

func canMove(grid [][]gridContents, row int, col int, dir Direction) bool {
	if grid[row][col] == WALL {
		return false
	}
	if grid[row][col] == EMPTY {
		return true
	}
	rowDelta := 1
	switch dir {
	case LEFT:
		return canMove(grid, row, col-2, dir)
	case RIGHT:
		return canMove(grid, row, col+2, dir)
	case UP:
		rowDelta = -1
		fallthrough
	case DOWN:
		if canMove(grid, row+rowDelta, col, dir) {
			if grid[row][col] == CRATE_LEFT {
				return canMove(grid, row+rowDelta, col+1, dir)
			} else {
				return canMove(grid, row+rowDelta, col-1, dir)
			}
		} else {
			return false
		}
	}
	panic("Should be unreachable")
}
