package main

import (
	"log/slog"
	"slices"
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

type page struct {
	number   string
	prereqOf []*page
}

func getOrCreatePage(pages map[string]*page, number string) *page {
	thePage, found := pages[number]
	if !found {
		thePage = &page{number, make([]*page, 0, 10)}
		pages[number] = thePage
	}
	return thePage
}

func makePageIllegal(legalPages map[string]nothing, thePage *page) {
	delete(legalPages, thePage.number)
	for _, prereq := range thePage.prereqOf {
		delete(legalPages, prereq.number)
	}
}

type p1Context struct {
	allPages         map[string]*page
	incorrectUpdates [][]string
}

func Day5Part1(logger *slog.Logger, input string) (string, any) {
	allPages := make(map[string]*page)
	lines := strings.Fields(input)
	ix := 0
	var line string
	for ix, line = range lines {
		if line[2] != '|' {
			break
		}
		firstPage := getOrCreatePage(allPages, line[0:2])
		secondPage := getOrCreatePage(allPages, line[3:])
		firstPage.prereqOf = append(firstPage.prereqOf, secondPage)
	}

	sum := 0
	incorrectUpdates := make([][]string, 0, len(lines))
	for ; ix < len(lines); ix++ {
		pagesInThisUpdate := strings.Split(lines[ix], ",")
		slices.Reverse(pagesInThisUpdate)
		if checkUpdate(allPages, pagesInThisUpdate) == -1 {
			middlePage, _ := strconv.Atoi(pagesInThisUpdate[len(pagesInThisUpdate)/2])
			sum += middlePage
		} else {
			incorrectUpdates = append(incorrectUpdates, pagesInThisUpdate)
		}
	}
	return strconv.Itoa(sum), p1Context{allPages, incorrectUpdates}
}

// Returns -1 for a valid update, else index of the last illegally-placed page.
func checkUpdate(allPages map[string]*page, pagesInThisUpdate []string) int {
	legalPages := make(map[string]nothing, len(allPages))
	for page := range allPages {
		legalPages[page] = nothing{}
	}
	invalidIndex := -1
	for ix, pageNumber := range pagesInThisUpdate {
		_, found := legalPages[pageNumber]
		if !found {
			invalidIndex = ix
			break
		}
		thePage := allPages[pageNumber]
		makePageIllegal(legalPages, thePage)
	}
	return invalidIndex
}

func Day5Part2(logger *slog.Logger, input string, part1Context any) string {
	context := part1Context.(p1Context)
	allPages := context.allPages
	sum := 0
	for ix := 0; ix < len(context.incorrectUpdates); ix++ {
		pagesInThisUpdate := context.incorrectUpdates[ix]
		if invalidPageIx := checkUpdate(allPages, pagesInThisUpdate); invalidPageIx == -1 {
			middlePage, _ := strconv.Atoi(pagesInThisUpdate[len(pagesInThisUpdate)/2])
			sum += middlePage
		} else {
			pagesInThisUpdate[invalidPageIx], pagesInThisUpdate[invalidPageIx-1] = pagesInThisUpdate[invalidPageIx-1], pagesInThisUpdate[invalidPageIx]
			ix--
		}
	}
	return strconv.Itoa(sum)
}
