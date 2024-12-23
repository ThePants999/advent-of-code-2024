package main

import (
	"log/slog"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/set"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day23 = runner.DayImplementation{
	DayNumber:    23,
	ExecutePart1: Day23Part1,
	ExecutePart2: Day23Part2,
	ExampleInput: `kh-tc
qp-kh
de-cg
ka-co
yn-aq
qp-ub
cg-tb
vc-aq
tb-ka
wh-tc
yn-cg
kh-ub
ta-co
de-co
tc-td
tb-wq
wh-td
ta-ka
td-qp
aq-cg
wq-ub
ub-vc
de-ta
wq-aq
wq-vc
wh-yn
ka-de
kh-ta
co-tc
wh-qp
tb-vc
td-yn`,
	ExamplePart1Answer: "7",
	ExamplePart2Answer: "co,de,ka,ta",
}

type d23Computer struct {
	name        string
	connections *set.Set
}

func newComputer(name string) d23Computer {
	return d23Computer{name, set.New()}
}

func Day23Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	computers := make(map[string]*d23Computer)
	for _, line := range lines {
		comp1, found := computers[line[0:2]]
		if !found {
			computer := newComputer(line[0:2])
			comp1 = &computer
			computers[computer.name] = comp1
		}
		comp2, found := computers[line[3:5]]
		if !found {
			computer := newComputer(line[3:5])
			comp2 = &computer
			computers[computer.name] = comp2
		}
		comp1.connections.Insert(comp2)
		comp2.connections.Insert(comp1)
	}

	sum := 0
	handled := set.New()
	for comp1Name, comp1 := range computers {
		// Only interested in computers starting with T
		if comp1Name[0] != 't' {
			continue
		}

		connsSlice := make([]*d23Computer, 0, comp1.connections.Len())
		comp1.connections.Do(func(item any) {
			connsSlice = append(connsSlice, item.(*d23Computer))
		})
		for ix, comp2 := range connsSlice {
			if handled.Has(comp2) {
				// This computer has already been comp1 in the
				// past so we've counted all its connections
				// already
				continue
			}

			// Now see if any of the others are also connected to
			// this one
			for ix2 := ix + 1; ix2 < len(connsSlice); ix2++ {
				comp3 := connsSlice[ix2]
				if handled.Has(comp3) {
					// This computer has already been comp1 in the
					// past so we've counted all its connections
					// already
					continue
				}

				if comp2.connections.Has(comp3) {
					// comp2 and comp3 are connected to each other
					// as well as comp1
					sum++
				}
			}
		}

		handled.Insert(comp1)
	}

	return strconv.Itoa(sum), computers
}

func Day23Part2(logger *slog.Logger, input string, part1Context any) string {
	bestSet := set.New()

	computers := set.New()
	for comp := range maps.Values(part1Context.(map[string]*d23Computer)) {
		computers.Insert(comp)
	}
	bronKerbosch(set.New(), computers, set.New(), &bestSet)

	finalSet := make([]string, 0, bestSet.Len())
	bestSet.Do(func(item any) {
		finalSet = append(finalSet, item.(*d23Computer).name)
	})
	slices.Sort(finalSet)
	var sb strings.Builder
	for _, name := range finalSet {
		sb.WriteString(name)
		sb.WriteRune(',')
	}
	answer := sb.String()
	answer = answer[:len(answer)-1]

	return answer
}

func bronKerbosch(r *set.Set, p *set.Set, x *set.Set, best **set.Set) {
	if p.Len() == 0 && x.Len() == 0 {
		// There's nothing more we could add, r is a maximal clique
		if r.Len() > (*best).Len() {
			// And it's bigger than any other we've found so far
			*best = r
		}
		return
	}
	if r.Len()+p.Len() <= (*best).Len() {
		// We don't have enough candidate vertices left to exceed the
		// biggest we've already found, give up here
		return
	}

	p.Do(func(item any) {
		itemSet := set.New(item)
		neighbourSet := item.(*d23Computer).connections
		bronKerbosch(r.Union(itemSet), p.Intersection(neighbourSet), x.Intersection(neighbourSet), best)
		p.Remove(item)
		x.Insert(item)
	})
}
