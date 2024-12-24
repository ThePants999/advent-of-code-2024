package main

import (
	"log/slog"
	"math"
	"regexp"
	"strconv"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day14 = runner.DayImplementation{
	DayNumber:          14,
	ExecutePart1:       Day14Part1,
	ExecutePart2:       Day14Part2,
	ExampleInput:       "",
	ExamplePart1Answer: "",
	ExamplePart2Answer: "",
}

type d14robot struct {
	pos    gridPos
	vector gridPos
}

type d14context struct {
	robots     []d14robot
	areaHeight int
	areaWidth  int
}

func Day14Part1(logger *slog.Logger, input string) (string, any) {
	// Parse input. Quick bit of regex practice
	// can't hurt I suppose.
	re := regexp.MustCompile(`[^\d-]+([-\d]+)[^\d-]+([-\d]+)[^\d-]+([-\d]+)[^\d-]+([-\d]+)`)
	matches := re.FindAllStringSubmatch(input, -1)
	robots := make([]d14robot, len(matches))
	for ix, match := range matches {
		r := &robots[ix]
		r.pos.col, _ = strconv.Atoi(match[1])
		r.pos.row, _ = strconv.Atoi(match[2])
		r.vector.col, _ = strconv.Atoi(match[3])
		r.vector.row, _ = strconv.Atoi(match[4])
	}

	areaWidth, areaHeight := 101, 103
	middleCol, middleRow := areaWidth/2, areaHeight/2

	// We need to simulate each robot independently, but we don't need
	// to do each second one by one - the fact they wrap from one side
	// to the other means we can just multiply and take the modulus.
	var robotCounts [4]int
	for _, robot := range robots {
		finalRow := (robot.pos.row + (robot.vector.row * 100)) % areaHeight
		finalCol := (robot.pos.col + (robot.vector.col * 100)) % areaWidth
		if finalRow < 0 {
			finalRow += areaHeight
		}
		if finalCol < 0 {
			finalCol += areaWidth
		}

		// Calculate which quadrant this robot ends up in.
		if finalRow < middleRow {
			if finalCol < middleCol {
				robotCounts[0]++
			} else if finalCol > middleCol {
				robotCounts[1]++
			}
		} else if finalRow > middleRow {
			if finalCol < middleCol {
				robotCounts[2]++
			} else if finalCol > middleCol {
				robotCounts[3]++
			}
		}
	}

	return strconv.Itoa(robotCounts[0] * robotCounts[1] * robotCounts[2] * robotCounts[3]), d14context{robots, areaHeight, areaWidth}
}

func Day14Part2(logger *slog.Logger, input string, part1Context any) string {
	// It's fun time.
	//
	// The approach I've taken here is to assume that the picture will
	// appear when variance between robots is at a minimum for both row
	// and column values - i.e. when they're the most clustered.
	//
	// Each robot moves independently along X and Y axes - there's no
	// relationship between their movement on each axis, so we can
	// consider them separately.
	//
	// Different robots may return to their starting X and/or Y
	// co-ordinate on different intervals, but with the grid being 101
	// wide by 103 tall, we know that the motion of ALL robots along
	// columns will repeat every 101 seconds, and along rows every 103
	// seconds. We can therefore model the first 103 seconds to find
	// the point of lowest variance among row coordinates and among
	// column coordinates.
	context := part1Context.(d14context)

	rowSum, colSum := 0, 0
	for _, robot := range context.robots {
		rowSum += robot.pos.row
		colSum += robot.pos.col
	}
	meanRow, meanCol := float64(rowSum)/float64(len(context.robots)), float64(colSum)/float64(len(context.robots))

	numRobots := float64(len(context.robots))
	maxDimension := context.areaHeight

	minRowVariance, minRowVarianceIteration, minColVariance, minColVarianceIteration := math.MaxFloat64, 0, math.MaxFloat64, 0
	for i := 1; i <= maxDimension; i++ {
		for ix := range context.robots {
			newRow := context.robots[ix].pos.row + context.robots[ix].vector.row
			if newRow >= context.areaHeight {
				newRow -= context.areaHeight
			} else if newRow < 0 {
				newRow += context.areaHeight
			}
			// Rather than recalculate the mean in full every second, we can
			// tweak it according to each robot's delta as we process that
			// robot.
			meanRow += float64(newRow-context.robots[ix].pos.row) / numRobots
			context.robots[ix].pos.row = newRow

			newCol := context.robots[ix].pos.col + context.robots[ix].vector.col
			if newCol >= context.areaWidth {
				newCol -= context.areaWidth
			} else if newCol < 0 {
				newCol += context.areaWidth
			}
			meanCol += float64(newCol-context.robots[ix].pos.col) / numRobots
			context.robots[ix].pos.col = newCol
		}

		var rowVariance, colVariance float64
		for _, robot := range context.robots {
			rowDiff := float64(robot.pos.row) - meanRow
			rowVariance += rowDiff * rowDiff
			colDiff := float64(robot.pos.col) - meanCol
			colVariance += colDiff * colDiff
		}

		if rowVariance < minRowVariance {
			minRowVariance = rowVariance
			minRowVarianceIteration = i
		}
		if colVariance < minColVariance {
			minColVariance = colVariance
			minColVarianceIteration = i
		}
	}

	// The bit above was me. Now comes some maths bullshit that I
	// shamelessly stole because I do AoC to practice programming,
	// not maths. We need to find the point where the lowest-row-
	// variance (which repeats every 103 seconds) coincides with
	// the lowest-column-variance (which repeats every 101 seconds).
	// The answer involves Chinese remainder theorem...
	magicAnswer := minColVarianceIteration + (((51 * (minRowVarianceIteration - minColVarianceIteration)) % context.areaHeight) * context.areaWidth)
	if magicAnswer < 0 {
		magicAnswer += (context.areaHeight * context.areaWidth)
	}
	return strconv.Itoa(magicAnswer)
}
