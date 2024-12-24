package main

import (
	"log/slog"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day2 = runner.DayImplementation{
	DayNumber:          2,
	ExecutePart1:       Day2Part1,
	ExecutePart2:       Day2Part2,
	ExampleInput:       "7 6 4 2 1\n1 2 7 8 9\n9 7 6 2 1\n1 3 2 4 5\n8 6 4 4 1\n1 3 6 7 9",
	ExamplePart1Answer: "2",
	ExamplePart2Answer: "4",
}

func Day2Part1(logger *slog.Logger, input string) (string, any) {
	// Parse input into a slice of slices of numbers.
	reportStrings := strings.Split(input, "\n")
	// Cope with a blank line at the end.
	if len(reportStrings[len(reportStrings)-1]) == 0 {
		reportStrings = reportStrings[:len(reportStrings)-1]
	}
	reports := make([][]int, len(reportStrings))
	for ix, report := range reportStrings {
		levels := strings.Split(report, " ")
		reports[ix] = make([]int, len(levels))
		for iy, level := range levels {
			levelInt, err := strconv.Atoi(level)
			if err != nil {
				panic(err)
			}
			reports[ix][iy] = levelInt
		}
	}

	// Check each report in turn to count safe reports.
	safeCount := 0
	for _, report := range reports {
		if safe, _ := checkReport(report); safe {
			safeCount++
		}
	}

	return strconv.Itoa(safeCount), reports
}

func Day2Part2(logger *slog.Logger, input string, part1Context any) string {
	reports := part1Context.([][]int)
	safeCount := 0
	for _, report := range reports {
		safe, problemIx := checkReport(report)
		if !safe {
			// Try again without the level that caused a problem.
			newReport := make([]int, problemIx, len(report))
			copy(newReport, report[:problemIx])
			newReport = append(newReport, report[problemIx+1:]...)
			safe, _ = checkReport(newReport)
		}
		if !safe {
			// It's also possible that it was the previous level
			// that was the real problem. E.g. in the sequence
			// 1 5 6 7..., the 5 is what we identify as unsafe,
			// but it's really the 1 that's the issue.
			newReport := make([]int, problemIx-1, len(report))
			copy(newReport, report[:problemIx-1])
			newReport = append(newReport, report[problemIx:]...)
			safe, _ = checkReport(newReport)
		}
		if !safe && problemIx == 2 {
			// Finally, there's a special case where if a problem
			// is detected with the third element, it might actually
			// be removing the FIRST that helps. For example, in
			// the sequence 3 1 2 4..., the drop from 3 to 1 causes
			// us to think it's a decreasing report, at which point
			// the 2 looks problematic, when it's actually an
			// increasing report where the 3 is problematic.
			safe, _ = checkReport(report[1:])
		}
		if safe {
			safeCount++
		}
	}

	return strconv.Itoa(safeCount)
}

func checkReport(report []int) (bool, int) {
	increasing := true
	safe := true
	problemIx := -1
	if report[1] < report[0] {
		// Assume that a decrease from the first to
		// the second level means this is a
		// decreasing report.
		increasing = false
	}
out:
	for ix := 1; ix < len(report); ix++ {
		// Report a problem any time the level goes up
		// in a decreasing report or down in an
		// increasing report, or any time it changes by
		// more than 3 in either direction.
		switch {
		case increasing && report[ix] <= report[ix-1]:
			fallthrough
		case !increasing && report[ix] >= report[ix-1]:
			fallthrough
		case report[ix]-report[ix-1] > 3:
			fallthrough
		case report[ix]-report[ix-1] < -3:
			safe = false
			problemIx = ix
			break out
		}
	}
	return safe, problemIx
}
