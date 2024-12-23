package main

import (
	"iter"
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

type ComputerSet map[string]*d23Computer

func (set ComputerSet) Has(key string) bool {
	_, found := set[key]
	return found
}

func (first ComputerSet) Difference(second ComputerSet) ComputerSet {
	set := make(ComputerSet, len(first))
	for k, v := range first {
		_, found := second[k]
		if !found {
			set[k] = v
		}
	}
	return set
}

func (first ComputerSet) Intersection(second ComputerSet) ComputerSet {
	var set ComputerSet
	if len(first) > len(second) {
		set = make(ComputerSet, len(second))
		for k, v := range second {
			_, found := first[k]
			if found {
				set[k] = v
			}
		}
	} else {
		set = make(ComputerSet, len(first))
		for k, v := range first {
			_, found := second[k]
			if found {
				set[k] = v
			}
		}
	}
	return set
}

func (set ComputerSet) findOrAdd(name string) *d23Computer {
	computer, found := set[name]
	if !found {
		comp := newComputer(name)
		computer = &comp
		set[name] = computer
	}
	return computer
}

type d23Computer struct {
	name        string
	connections ComputerSet
}

func newComputer(name string) d23Computer {
	return d23Computer{name, make(ComputerSet)}
}

func Day23Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	computers := make(ComputerSet)
	for _, line := range lines {
		comp1 := computers.findOrAdd(line[0:2])
		comp2 := computers.findOrAdd(line[3:5])
		comp1.connections[comp2.name] = comp2
		comp2.connections[comp1.name] = comp1
	}

	sum := 0
	handled := set.New()
	for comp1Name, comp1 := range computers {
		// Only interested in computers starting with T
		if comp1Name[0] != 't' {
			continue
		}

		connsSlice := slices.Collect(maps.Values(comp1.connections))
		for ix, comp2 := range connsSlice {
			if handled.Has(comp2.name) {
				// This computer has already been comp1 in the
				// past so we've counted all its connections
				// already
				continue
			}

			// Now see if any of the others are also connected to
			// this one
			for ix2 := ix + 1; ix2 < len(connsSlice); ix2++ {
				comp3 := connsSlice[ix2]
				if handled.Has(comp3.name) {
					// This computer has already been comp1 in the
					// past so we've counted all its connections
					// already
					continue
				}

				if comp2.connections.Has(comp3.name) {
					// comp2 and comp3 are connected to each other
					// as well as comp1
					sum++
				}
			}
		}

		handled.Insert(comp1.name)
	}

	return strconv.Itoa(sum), computers
}

func Day23Part2(logger *slog.Logger, input string, part1Context any) string {
	bestSet := make(ComputerSet)
	computers := part1Context.(ComputerSet)
	bronKerbosch(make(ComputerSet), computers, make(ComputerSet), &bestSet)

	finalSet := slices.Collect(maps.Keys(bestSet))
	slices.Sort(finalSet)
	return strings.Join(finalSet, ",")
}

func bronKerbosch(r ComputerSet, p ComputerSet, x ComputerSet, best *ComputerSet) {
	if len(r)+len(p) <= len(*best) {
		// We don't have enough candidate vertices left to exceed the
		// biggest we've already found, give up here
		return
	}

	if len(p) == 0 && len(x) == 0 {
		// There's nothing more we could add, r is a maximal clique
		if len(r) > len(*best) {
			// And it's bigger than any other we've found so far
			*best = r
		}
		return
	}

	next, stop := iter.Pull(maps.Values(p))
	pivot, ok := next()
	stop()
	if !ok {
		next, stop = iter.Pull(maps.Values(x))
		pivot, _ = next()
		stop()
	}

	for k, v := range p.Difference(pivot.connections) {
		newR := maps.Clone(r)
		newR[k] = v
		pIntersectNeighbours := p.Intersection(v.connections)
		xIntersectNeighbours := x.Intersection(v.connections)
		bronKerbosch(newR, pIntersectNeighbours, xIntersectNeighbours, best)
		delete(p, k)
		x[k] = v
	}
}
