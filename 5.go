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

// Record that both _pageNum_ and anything for which
// it's a prerequisite is no longer valid after this
// point.
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

// Record that _page_ must come before _prereqOf_, by
// adding _prereqOf_ to the slice recorded for _page_.
func addPrereq(pages [][]int, page int, prereqOf int) {
	if pages[page] == nil {
		pages[page] = make([]int, 1, 10)
		pages[page][0] = prereqOf
	} else {
		pages[page] = append(pages[page], prereqOf)
	}
}

func Day5Part1(logger *slog.Logger, input string) (string, any) {
	// Parse the first half of the input. What we're going
	// to construct here is a "map" of prerequisites,
	// so if we see "47|53" , we add 53 to the set of pages
	// for which 47 is a prerequisite. (Since all page numbers
	// are two digits, we use an array rather than a map for
	// speed.)
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

	// We now go through the second part of the input, both
	// parsing and solving simultaneously.
	//
	// The approach we take here is to reverse the pages in
	// the update. We can then make use of our prereqs list
	// - any time we see page X, because we're effectively
	// working backwards, we know that anything for which X
	// is a prerequisite is not valid "after" this point.
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
			// This update is valid - and fortunately, reversing
			// didn't change the middle element!
			sum += pagesInThisUpdate[numPages/2]
		} else {
			// This update isn't legal. Record it ready for part 2.
			incorrectUpdates = append(incorrectUpdates, pagesInThisUpdate)
		}
	}
	return strconv.Itoa(sum), p1Context{prereqs, incorrectUpdates}
}

// Returns -1 for a valid update, else index of the last illegally-placed page.
func checkUpdate(prereqs [][]int, pagesInThisUpdate []int) int {
	pageIllegal := make([]bool, NUM_PAGES)
	invalidIndex := -1
	for ix, pageNumber := range pagesInThisUpdate {
		if pageIllegal[pageNumber] {
			invalidIndex = ix
			break
		}
		// Having seen page X, it's not valid to see either
		// X again, or anything for which it's a prerequisite,
		// later in the list.  (Remember, we're looking
		// through the list backwards.)
		makePageIllegal(prereqs, pageIllegal, pageNumber)
	}
	return invalidIndex
}

func Day5Part2(logger *slog.Logger, input string, part1Context any) string {
	context := part1Context.(p1Context)
	prereqs := context.prereqs

	// Go through each update that was identified as illegal in part 1.
	sum := 0
	for ix := 0; ix < len(context.incorrectUpdates); ix++ {
		pagesInThisUpdate := context.incorrectUpdates[ix]
		// The same function that told us this was illegal can also tell
		// us where the problematic element is. The problem statement
		// allows us to assume that every update can be made legal, so
		// we take a simple approach whereby every time we find a
		// problem, we swap the problematic page "forwards" (which is
		// really backwards since the pages are reversed) and try
		// again.
		if invalidPageIx := checkUpdate(prereqs, pagesInThisUpdate); invalidPageIx == -1 {
			sum += pagesInThisUpdate[len(pagesInThisUpdate)/2]
		} else {
			pagesInThisUpdate[invalidPageIx], pagesInThisUpdate[invalidPageIx-1] = pagesInThisUpdate[invalidPageIx-1], pagesInThisUpdate[invalidPageIx]
			ix--
		}
	}
	return strconv.Itoa(sum)
}
