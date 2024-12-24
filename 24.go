package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

var Day24 = runner.DayImplementation{
	DayNumber:    24,
	ExecutePart1: Day24Part1,
	ExecutePart2: Day24Part2,
	ExampleInput: `x00: 1
x01: 0
x02: 1
x03: 1
x04: 0
y00: 1
y01: 1
y02: 1
y03: 1
y04: 1

ntg XOR fgs -> mjb
y02 OR x01 -> tnw
kwq OR kpj -> z05
x00 OR x03 -> fst
tgd XOR rvg -> z01
vdt OR tnw -> bfw
bfw AND frj -> z10
ffh OR nrd -> bqk
y00 AND y03 -> djm
y03 OR y00 -> psh
bqk OR frj -> z08
tnw OR fst -> frj
gnj AND tgd -> z11
bfw XOR mjb -> z00
x03 OR x00 -> vdt
gnj AND wpb -> z02
x04 AND y00 -> kjc
djm OR pbm -> qhw
nrd AND vdt -> hwm
kjc AND fst -> rvg
y04 OR y02 -> fgs
y01 AND x02 -> pbm
ntg OR kjc -> kwq
psh XOR fgs -> tgd
qhw XOR tgd -> z09
pbm OR djm -> kpj
x03 XOR y03 -> ffh
x00 XOR y04 -> ntg
bfw OR bqk -> z06
nrd XOR fgs -> wpb
frj XOR qhw -> z04
bqk OR frj -> z07
y03 OR x01 -> nrd
hwm AND bqk -> z03
tgd XOR rvg -> z12
tnw OR pbm -> gnj`,
	ExamplePart1Answer: "2024",
	ExamplePart2Answer: "",
}

type d24Operator int

const (
	D24_AND d24Operator = iota
	D24_OR
	D24_XOR
)

type d24InputProvider interface {
	provide() bool
	drawGraph(*graphviz.Graph, *cgraph.Node)
}

type d24FixedInput struct {
	name  string
	value bool
	node  *cgraph.Node
}

func (input *d24FixedInput) provide() bool {
	return input.value
}

func (input *d24FixedInput) drawGraph(graph *graphviz.Graph, to *cgraph.Node) {
	var err error
	if input.node == nil {
		input.node, err = graph.CreateNodeByName(input.name)
		if err != nil {
			panic(err)
		}
	}
	var edge *cgraph.Edge
	edge, err = graph.CreateEdgeByName(input.name, input.node, to)
	if err != nil {
		panic(err)
	}
	edge.SetLabel("")
}

type d24Wire struct {
	name     string
	provider d24InputProvider
}

type d24Gate struct {
	operator d24Operator
	input1N  string
	input2N  string
	outputN  string
	input1   *d24Wire
	input2   *d24Wire
	output   *d24Wire
	level    int
	node     *cgraph.Node
}

func (gate *d24Gate) operatorString() string {
	switch gate.operator {
	case D24_AND:
		return "AND"
	case D24_OR:
		return "OR"
	default:
		return "XOR"
	}
}

func (gate *d24Gate) drawGraph(graph *graphviz.Graph, to *cgraph.Node) {
	var err error
	if gate.node == nil {
		gate.node, err = graph.CreateNodeByName(fmt.Sprintf("%s %s %s", gate.input1N, gate.operatorString(), gate.input2N))
		if err != nil {
			panic(err)
		}
		gate.node.SetLabel(gate.operatorString())
		//gate.node.SetPenWidth(3.0)
		gate.node.SetShape(cgraph.TriangleShape)
	}
	var edge *cgraph.Edge
	edge, err = graph.CreateEdgeByName(gate.outputN, gate.node, to)
	if err != nil {
		panic(err)
	}
	edge.SetLabel(gate.outputN)
	gate.input1.provider.drawGraph(graph, gate.node)
	gate.input2.provider.drawGraph(graph, gate.node)
}

func (gate *d24Gate) provide() bool {
	switch gate.operator {
	case D24_AND:
		return gate.input1.provider.provide() && gate.input2.provider.provide()
	case D24_OR:
		return gate.input1.provider.provide() || gate.input2.provider.provide()
	default:
		input1, input2 := gate.input1.provider.provide(), gate.input2.provider.provide()
		return (input1 || input2) && !(input1 && input2)
	}
}

func (gate *d24Gate) setLevel(level int) int {
	maxLevel := level
	gate.level = level
	gate1, ok := gate.input1.provider.(*d24Gate)
	if ok {
		maxLevel = gate1.setLevel(level + 1)
	}
	gate2, ok := gate.input2.provider.(*d24Gate)
	if ok {
		newMax := gate2.setLevel(level + 1)
		if newMax > maxLevel {
			maxLevel = newMax
		}
	}
	return maxLevel
}

func newGate(operator d24Operator, input1 string, input2 string, output string, wires map[string]*d24Wire) *d24Gate {
	gate := d24Gate{operator, input1, input2, output, nil, nil, nil, -1, nil}
	var found bool
	gate.input1, found = wires[input1]
	if !found {
		gate.input1 = &d24Wire{input1, nil}
		wires[input1] = gate.input1
	}
	gate.input2, found = wires[input2]
	if !found {
		gate.input2 = &d24Wire{input2, nil}
		wires[input2] = gate.input2
	}
	gate.output, found = wires[output]
	if !found {
		gate.output = &d24Wire{output, &gate}
		wires[output] = gate.output
	} else {
		gate.output.provider = &gate
	}
	return &gate
}

type d24Circuit struct {
	wires        map[string]*d24Wire
	gates        map[string]*d24Gate
	gatesByLevel [][]*d24Gate
	outputWires  []*d24Wire
	correctZ     int64
	originalZ    int64
}

func newCircuit(wires map[string]*d24Wire, gates map[string]*d24Gate) *d24Circuit {
	outputWires := make([]*d24Wire, 64)
	var x, y int64
	maxLevel := 0
	for ix := range 64 {
		onesRune := '0' + rune(ix%10)
		tensRune := '0' + rune(ix/10)
		outputWires[ix] = wires[string([]rune{'z', tensRune, onesRune})]

		if outputWires[ix] != nil {
			outputGate := outputWires[ix].provider.(*d24Gate)
			newMax := outputGate.setLevel(0)
			if newMax > maxLevel {
				maxLevel = newMax
			}
		}

		xWire, found := wires[string([]rune{'x', tensRune, onesRune})]
		if found && xWire.provider.provide() {
			x |= 1 << ix
		}

		yWire, found := wires[string([]rune{'y', tensRune, onesRune})]
		if found && yWire.provider.provide() {
			y |= 1 << ix
		}
	}

	gatesByLevel := make([][]*d24Gate, maxLevel+1)
	for ix := range maxLevel {
		gatesByLevel[ix] = make([]*d24Gate, 0)
	}
	for gate := range maps.Values(gates) {
		gatesByLevel[gate.level] = append(gatesByLevel[gate.level], gate)
	}

	return &d24Circuit{wires, gates, gatesByLevel, outputWires, x + y, 0}
}

func (circuit *d24Circuit) zValue() int64 {
	var z int64
	for ix := range 64 {
		if circuit.outputWires[ix] != nil && circuit.outputWires[ix].provider.provide() {
			z |= 1 << ix
		}
	}
	return z
}

func parseD24Input(input string) *d24Circuit {
	lines := strings.Split(input, "\n")

	// Parse the first section - wires
	wires := make(map[string]*d24Wire)
	for ix, line := range lines {
		if len(line) == 0 {
			lines = lines[ix+1:]
			break
		}

		val := false
		if line[5] == '1' {
			val = true
		}

		input := d24FixedInput{line[0:3], val, nil}
		wires[line[0:3]] = &d24Wire{line[0:3], &input}
	}

	// Cope with blank line at end
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[0 : len(lines)-1]
	}

	// Parse the second section - gates
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

	circuit := newCircuit(wires, gates)
	circuit.originalZ = circuit.zValue()
	return circuit
}

func Day24Part1(logger *slog.Logger, input string) (string, any) {
	circuit := parseD24Input(input)
	return strconv.Itoa(int(circuit.originalZ)), circuit
}

func Day24Part2(logger *slog.Logger, input string, part1Context any) string {
	// if len(input) < 100 {
	// 	return ""
	// }

	circuit := part1Context.(*d24Circuit)

	ctx := context.Background()
	g, err := graphviz.New(ctx)
	if err != nil {
		panic(err)
	}

	graph, err := g.Graph()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			panic(err)
		}
		g.Close()
	}()

	for _, wire := range circuit.outputWires {
		if wire != nil {
			node, err := graph.CreateNodeByName(wire.name)
			if err != nil {
				panic(err)
			}
			wire.provider.drawGraph(graph, node)
		}
	}

	// if err := g.RenderFilename(ctx, graph, graphviz.PNG, "./graph.png"); err != nil {
	// 	panic(err)
	// }
	var buf bytes.Buffer
	if err := g.Render(ctx, graph, "dot", &buf); err != nil {
		panic(err)
	}
	os.WriteFile("./graph.dot", buf.Bytes(), 0777)

	return ""
}

func resultStrictlyBetter(originalResult int64, newResult int64, correctResult int64) bool {
	originalWrong := correctResult ^ originalResult
	newWrong := correctResult ^ newResult
	return newWrong < originalWrong && (newWrong|originalWrong) == originalWrong
}

func trySwapGates(circuit *d24Circuit, gate1 *d24Gate, gate2 *d24Gate, gate3 *d24Gate, gate4 *d24Gate) bool {
	swapGates(gate1, gate2)
	swapGates(gate3, gate4)
	fixed := circuit.zValue() == circuit.correctZ
	swapGates(gate1, gate2)
	swapGates(gate3, gate4)
	return fixed
	//better := resultStrictlyBetter(circuit.originalZ, circuit.zValue(), circuit.correctZ)
	//return better
}

func swapGates(gate1 *d24Gate, gate2 *d24Gate) {
	gate1.output.provider = gate2
	gate2.output.provider = gate1
}
