package main

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day24 = runner.DayImplementation{
	DayNumber:          24,
	ExecutePart1:       Day24Part1,
	ExecutePart2:       Day24Part2,
	ExampleInput:       ``,
	ExamplePart1Answer: "",
	ExamplePart2Answer: "",
}

type d24Operator int

const (
	D24_AND d24Operator = iota
	D24_OR
	D24_XOR
)

// What we have in here is a carefully structured set of
// structs and relationships between them that allows us
// to build the graph of the circuit and then achieve all
// of (a) traversing it in any direction, (b) modifying
// it by swapping gate outputs, and (c) efficiently
// recalculating the value being output on the Z wires
// after doing (b).
//
// That turned out to be very unnecessary ;-( but hey,
// it's still cool, right?

// Something that provides the input being sent down a
// wire.
type d24InputProvider interface {
	provide() bool
}

// A constant input being sent down a wire (one of the X
// or Y values).
type d24FixedInput struct {
	name   string
	value  bool
	output *d24Wire
}

func (input *d24FixedInput) provide() bool {
	// Fixed inputs provide a constant value.
	return input.value
}

// A wire holding a single input/output. A wire can be
// between a fixed input and one or more downstream gates,
// between an upstream gate and one or more downstream
// gates, or after an upstream gate as a Z-output.
type d24Wire struct {
	name            string
	provider        d24InputProvider
	downstreamGates []*d24Gate
}

func newWire(name string, upstream d24InputProvider) *d24Wire {
	wire := d24Wire{name, upstream, make([]*d24Gate, 0, 2)}
	return &wire
}

func (wire *d24Wire) addDownstreamGate(gate *d24Gate) {
	wire.downstreamGates = append(wire.downstreamGates, gate)
}

// A gate. Always has two input wires and an output wire.
type d24Gate struct {
	operator d24Operator
	input1N  string
	input2N  string
	outputN  string
	input1   *d24Wire
	input2   *d24Wire
	output   *d24Wire
}

func (gate *d24Gate) provide() bool {
	// Gates provide a value determined by executing
	// their operator over the values provided by their
	// inputs.
	switch gate.operator {
	case D24_AND:
		return gate.input1.provider.provide() && gate.input2.provider.provide()
	case D24_OR:
		return gate.input1.provider.provide() || gate.input2.provider.provide()
	default:
		// XOR
		input1, input2 := gate.input1.provider.provide(), gate.input2.provider.provide()
		return (input1 || input2) && !(input1 && input2)
	}
}

// Instantiate a new gate. Will also set up any of the corresponding
// wires that haven't already been created, as well as establishing
// the relationships between them.
func newGate(operator d24Operator, input1 string, input2 string, output string, wires map[string]*d24Wire) *d24Gate {
	gate := d24Gate{operator, input1, input2, output, nil, nil, nil}
	var found bool

	gate.input1, found = wires[input1]
	if !found {
		gate.input1 = newWire(input1, nil)
		wires[input1] = gate.input1
	}
	gate.input1.addDownstreamGate(&gate)

	gate.input2, found = wires[input2]
	if !found {
		gate.input2 = newWire(input2, nil)
		wires[input2] = gate.input2
	}
	gate.input2.addDownstreamGate(&gate)

	gate.output, found = wires[output]
	if !found {
		gate.output = newWire(output, &gate)
		wires[output] = gate.output
	} else {
		gate.output.provider = &gate
	}

	return &gate
}

// The entire circuit
type d24Circuit struct {
	xInputs     []*d24FixedInput
	yInputs     []*d24FixedInput
	outputWires []*d24Wire
	correctZ    int64
	originalZ   int64
}

// Instantiate the circuit. We don't need to store
// all the wires and gates - they continue to exist,
// and can be found, courtesy of the connectivity
// between all the elements. We just store the
// externally-facing elements to facilitate
// interaction.
func newCircuit(wires map[string]*d24Wire, xInputs []*d24FixedInput, yInputs []*d24FixedInput) *d24Circuit {
	outputWires := make([]*d24Wire, 64)
	var x, y int64
	var maxXWire, maxYWire int
	for ix := range 64 {
		onesRune := '0' + rune(ix%10)
		tensRune := '0' + rune(ix/10)
		outputWires[ix] = wires[string([]rune{'z', tensRune, onesRune})]

		xWire, found := wires[string([]rune{'x', tensRune, onesRune})]
		if found {
			maxXWire = ix
			if xWire.provider.provide() {
				x |= 1 << ix
			}
		}

		yWire, found := wires[string([]rune{'y', tensRune, onesRune})]
		if found {
			maxYWire = ix
			if yWire.provider.provide() {
				y |= 1 << ix
			}
		}
	}

	circuit := &d24Circuit{xInputs, yInputs, outputWires, x + y, 0}
	circuit.originalZ = circuit.zValue()

	// Quick double-check of some assumptions.
	if maxXWire != maxYWire {
		panic("Different numbers of X and Y inputs")
	}
	if len(xInputs) != maxXWire+1 {
		panic("Unexpected numbering of X inputs")
	}
	if len(yInputs) != maxYWire+1 {
		panic("Unexpected numbering of Y inputs")
	}
	_, found := wires[fmt.Sprintf("z%02d", len(xInputs))]
	if !found {
		panic("Not enough Z wires")
	}
	_, found = wires[fmt.Sprintf("z%02d", len(xInputs)+1)]
	if found {
		panic("Too many Z wires")
	}

	return circuit
}

// Calculate the value currently being output on the
// z wires.
func (circuit *d24Circuit) zValue() int64 {
	var z int64
	for ix := range len(circuit.xInputs) + 1 {
		if circuit.outputWires[ix] != nil && circuit.outputWires[ix].provider.provide() {
			z |= 1 << ix
		}
	}
	return z
}

// Construct the entire circuit from the day's input.
func parseD24Input(input string) *d24Circuit {
	lines := strings.Split(input, "\n")

	// Parse the first section - wires corresponding
	// to fixed inputs.
	wires := make(map[string]*d24Wire)
	xInputs := make([]*d24FixedInput, 0, 64)
	yInputs := make([]*d24FixedInput, 0, 64)
	for ix, line := range lines {
		if len(line) == 0 {
			lines = lines[ix+1:]
			break
		}

		val := false
		if line[5] == '1' {
			val = true
		}

		input := &d24FixedInput{line[0:3], val, nil}
		wire := newWire(line[0:3], input)
		wires[line[0:3]] = wire
		input.output = wire
		if line[0] == 'x' {
			xInputs = append(xInputs, input)
		} else {
			yInputs = append(yInputs, input)
		}
	}

	// Cope with blank line at end
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[0 : len(lines)-1]
	}

	// Parse the second section - gates.
	gates := make(map[string]*d24Gate)
	for _, line := range lines {
		var op d24Operator
		postOpIndex := 7
		switch line[4] {
		case 'X':
			op = D24_XOR
		case 'A':
			op = D24_AND
		case 'O':
			op = D24_OR
			postOpIndex--
		}
		gates[line[postOpIndex+8:]] = newGate(op, line[0:3], line[postOpIndex+1:postOpIndex+4], line[postOpIndex+8:], wires)
	}

	return newCircuit(wires, xInputs, yInputs)
}

func Day24Part1(logger *slog.Logger, input string) (string, any) {
	// Part 1 is trivial - after constructing the circuit,
	// just return the initial Z value.
	circuit := parseD24Input(input)
	return strconv.Itoa(int(circuit.originalZ)), circuit
}

func Day24Part2(logger *slog.Logger, input string, part1Context any) string {
	// Skip example input, which isn't useful for part 2.
	if len(input) < 100 {
		return ""
	}

	// The approach we're going to take here is to firmly
	// assume what we're told in the problem - "the system
	// you're simulating is trying to add two binary numbers"
	// - and, furthermore, to assume that takes a very
	// particular form. A binary adder consists of a series
	// of sub-circuits each designed to calculate the value
	// of a single output bit. The sub-circuit for bit K
	// looks like this:
	// - The input bits X(K) and Y(K) are fed into two gates,
	//   XOR (output 1) and AND (output 2).
	// - Output 1 is XORed with the overflow indicator from
	//   the previous bit to give us the actual output for
	//   this bit, Z(K).
	// - Output 1 is also ANDed (output 3) with the overflow
	//   indicator from the previous bit.
	// - Finally, output 2 is ORed with output 3 to provide
	//   the overflow indicator for the next bit (output 4).
	//
	// We're going to look at each sub-circuit in turn and
	// validate it against this pattern, identifying what's
	// wrong if it doesn't conform.

	circuit := part1Context.(*d24Circuit)
	gatesToSwap := make([]*d24Gate, 0, 8)

	for ix, xInput := range circuit.xInputs {
		// X and Y inputs are paired together and should connect to
		// the same downstream gates.
		yInput := circuit.yInputs[ix]
		if yInput.name[1:] != xInput.name[1:] {
			panic("X and Y inputs don't match up")
		}
		if len(xInput.output.downstreamGates) != 2 || len(yInput.output.downstreamGates) != 2 {
			panic("X or Y input not connected to 2 gates")
		}
		if (xInput.output.downstreamGates[0] != yInput.output.downstreamGates[0] && xInput.output.downstreamGates[0] != yInput.output.downstreamGates[1]) || (xInput.output.downstreamGates[1] != yInput.output.downstreamGates[0] && xInput.output.downstreamGates[1] != yInput.output.downstreamGates[1]) {
			panic("X and Y inputs not connected to same gates")
		}

		if ix == 0 {
			// The first input has a very different shape, because
			// there's no overflow indicator from the previous bit
			// - we're just going to assume the problem isn't there,
			// accepting we won't cope if it is.
			continue
		}

		// The inputs should always connect to precisely two gates,
		// one of which is XOR and the other of which is AND -
		// validate this assumption.
		var upperXor, upperAnd *d24Gate
		if xInput.output.downstreamGates[0].operator == D24_XOR {
			upperXor = xInput.output.downstreamGates[0]
			upperAnd = xInput.output.downstreamGates[1]
			if upperAnd.operator != D24_AND {
				panic("X input connected to XOR but not AND")
			}
		} else {
			upperXor = xInput.output.downstreamGates[1]
			upperAnd = xInput.output.downstreamGates[0]
			if upperXor.operator != D24_XOR {
				panic("X input not connected to XOR")
			}
			if upperAnd.operator != D24_AND {
				panic("X input not connected to AND")
			}
		}

		if len(upperXor.output.downstreamGates) != 2 {
			// This gate should be upstream of an XOR and an AND, so
			// it's miswired.
			gatesToSwap = append(gatesToSwap, upperXor)
		}
		if len(upperAnd.output.downstreamGates) != 1 {
			// This gate should be upstream of an OR, so it's miswired.
			gatesToSwap = append(gatesToSwap, upperAnd)
		}

		var lowerXor, lowerAnd, lowerOr *d24Gate
		lowerXor = findDownstream(upperXor, upperAnd, D24_XOR)
		lowerAnd = findDownstream(upperXor, upperAnd, D24_AND)
		lowerOr = findDownstream(upperXor, upperAnd, D24_OR)

		if lowerXor != nil && len(lowerXor.output.downstreamGates) > 0 {
			// Lower XOR should output just to a z wire.
			gatesToSwap = append(gatesToSwap, lowerXor)
		}

		if lowerAnd != nil && (len(lowerAnd.output.downstreamGates) != 1 || lowerAnd.output.downstreamGates[0].operator != D24_OR) {
			// Lower AND should output just to an OR gate.
			gatesToSwap = append(gatesToSwap, lowerAnd)
		}

		if lowerOr != nil {
			// Lower OR should output to an XOR and an AND gate.
			if len(lowerOr.output.downstreamGates) == 2 {
				if (lowerOr.output.downstreamGates[0].operator != D24_XOR || lowerOr.output.downstreamGates[1].operator != D24_AND) && (lowerOr.output.downstreamGates[1].operator != D24_XOR || lowerOr.output.downstreamGates[0].operator != D24_AND) {
					gatesToSwap = append(gatesToSwap, lowerOr)
				}
			} else if len(lowerOr.output.downstreamGates) == 0 {
				// There's one exception - the lower or from the
				// final pair of inputs outputs to the final
				// output.
				if ix < len(circuit.xInputs)-1 || lowerOr.outputN != fmt.Sprintf("z%02d", len(circuit.xInputs)) {
					gatesToSwap = append(gatesToSwap, lowerOr)
				}
			} else {
				gatesToSwap = append(gatesToSwap, lowerOr)
			}
		}

		if len(gatesToSwap)%2 != 0 {
			panic("Odd number of gates swapped from one sub-circuit")
		}
	}

	// Let's just check this all adds up before we
	// confidently proclaim it the final result!
	if len(gatesToSwap) < 8 {
		panic("Haven't found all gates to swap")
	}
	if len(gatesToSwap) > 8 {
		panic("Found too many gates to swap")
	}
	swapGates(gatesToSwap[0], gatesToSwap[1])
	swapGates(gatesToSwap[2], gatesToSwap[3])
	swapGates(gatesToSwap[4], gatesToSwap[5])
	swapGates(gatesToSwap[6], gatesToSwap[7])
	if circuit.zValue() != circuit.correctZ {
		panic("Circuit still not outputting correct value after swaps")
	}

	// Output in the correct format - comma-separated,
	// alphabetically ordered list of output names.
	namesToSwap := make([]string, 8)
	for ix := range 8 {
		namesToSwap[ix] = gatesToSwap[ix].outputN
	}
	slices.Sort(namesToSwap)
	return strings.Join(namesToSwap, ",")
}

func findDownstream(upstream1 *d24Gate, upstream2 *d24Gate, operator d24Operator) *d24Gate {
	if len(upstream1.output.downstreamGates) >= 1 && upstream1.output.downstreamGates[0].operator == operator {
		return upstream1.output.downstreamGates[0]
	}
	if len(upstream1.output.downstreamGates) >= 2 && upstream1.output.downstreamGates[1].operator == operator {
		return upstream1.output.downstreamGates[1]
	}
	if len(upstream2.output.downstreamGates) >= 1 && upstream2.output.downstreamGates[0].operator == operator {
		return upstream2.output.downstreamGates[0]
	}
	if len(upstream2.output.downstreamGates) >= 2 && upstream2.output.downstreamGates[1].operator == operator {
		return upstream2.output.downstreamGates[1]
	}
	return nil
}

func swapGates(gate1 *d24Gate, gate2 *d24Gate) {
	gate1.output.provider = gate2
	gate2.output.provider = gate1
}
