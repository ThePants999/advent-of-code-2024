package main

import (
	"log/slog"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day5 = runner.DayImplementation{
	DayNumber:    5,
	ExecutePart1: Day5Part1,
	ExecutePart2: Day5Part2,
	ExampleInput: `47|53
97|13
97|61
97|47
75|29
61|13
75|53
29|13
97|29
53|29
61|53
97|53
61|29
47|13
75|47
97|75
47|61
75|61
47|29
75|13
53|13

75,47,61,53,29
97,61,53,29,13
75,29,13
75,97,47,61,53
61,13,29
97,13,75,29,47`,
	ExamplePart1Answer: "143",
	ExamplePart2Answer: "123",
}

const NUM_PAGES int = 100

func makePageIllegal(prereqs [][]int, pageIllegal []bool, pageNum int) {
	pageIllegal[pageNum] = true
	for _, prereq := range prereqs[pageNum] {
		pageIllegal[prereq] = true
	}
}

type p1Context struct {
	prereqs          [][]int
	incorrectUpdates [][]int
}

func addPrereq(pages [][]int, page int, prereqOf int) {
	if pages[page] == nil {
		pages[page] = make([]int, 1, 10)
		pages[page][0] = prereqOf
	} else {
		pages[page] = append(pages[page], prereqOf)
	}
}

func Day5Part1(logger *slog.Logger, input string) (string, any) {
	prereqs := make([][]int, NUM_PAGES)
	lines := strings.Fields(input)
	ix := 0
	var line string
	for ix, line = range lines {
		if line[2] != '|' {
			break
		}
		firstPage, _ := strconv.Atoi(line[0:2])
		secondPage, _ := strconv.Atoi(line[3:])
		addPrereq(prereqs, firstPage, secondPage)
	}

	sum := 0
	incorrectUpdates := make([][]int, 0, len(lines))
	for ; ix < len(lines); ix++ {
		pagesInThisUpdateStr := strings.Split(lines[ix], ",")
		numPages := len(pagesInThisUpdateStr)
		pagesInThisUpdate := make([]int, numPages)
		// Simultaneously convert to ints and reverse
		for i := range numPages {
			pagesInThisUpdate[numPages-i-1], _ = strconv.Atoi(pagesInThisUpdateStr[i])
		}

		if checkUpdate(prereqs, pagesInThisUpdate) == -1 {
			sum += pagesInThisUpdate[numPages/2]
		} else {
			incorrectUpdates = append(incorrectUpdates, pagesInThisUpdate)
		}
	}
	return strconv.Itoa(sum), p1Context{prereqs, incorrectUpdates}
}

// Returns -1 for a valid update, else index of the last illegally-placed page.
func checkUpdate(prereqs [][]int, pagesInThisUpdate []int) int {
	//legalPages := make(map[string]nothing, len(allPages))
	pageIllegal := make([]bool, NUM_PAGES)
	invalidIndex := -1
	for ix, pageNumber := range pagesInThisUpdate {
		if pageIllegal[pageNumber] {
			invalidIndex = ix
			break
		}
		makePageIllegal(prereqs, pageIllegal, pageNumber)
	}
	return invalidIndex
}

func Day5Part2(logger *slog.Logger, input string, part1Context any) string {
	context := part1Context.(p1Context)
	prereqs := context.prereqs
	sum := 0
	for ix := 0; ix < len(context.incorrectUpdates); ix++ {
		pagesInThisUpdate := context.incorrectUpdates[ix]
		if invalidPageIx := checkUpdate(prereqs, pagesInThisUpdate); invalidPageIx == -1 {
			sum += pagesInThisUpdate[len(pagesInThisUpdate)/2]
		} else {
			pagesInThisUpdate[invalidPageIx], pagesInThisUpdate[invalidPageIx-1] = pagesInThisUpdate[invalidPageIx-1], pagesInThisUpdate[invalidPageIx]
			ix--
		}
	}
	return strconv.Itoa(sum)
}
