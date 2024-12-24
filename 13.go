package main

import (
	"log/slog"
	"strconv"
	"strings"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

type d13machine struct {
	a_x     int
	a_y     int
	b_x     int
	b_y     int
	prize_x int
	prize_y int
}

var Day13 = runner.DayImplementation{
	DayNumber:    13,
	ExecutePart1: Day13Part1,
	ExecutePart2: Day13Part2,
	ExampleInput: `Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279`,
	ExamplePart1Answer: "480",
	ExamplePart2Answer: "875318608908",
}

func Day13Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Split(input, "\n")
	machines := make([]d13machine, 0, len(lines)/4+1)
	// Input parsing.
	// This is horrific, but it's fast ;-)
	for ix := 0; ix < len(lines); ix += 4 {
		machine := d13machine{}
		lineIx := 12
		for lines[ix][lineIx] != ',' {
			machine.a_x *= 10
			machine.a_x += int(lines[ix][lineIx] - '0')
			lineIx++
		}
		lineIx += 4
		for ; lineIx < len(lines[ix]); lineIx++ {
			machine.a_y *= 10
			machine.a_y += int(lines[ix][lineIx] - '0')
		}
		lineIx = 12
		for lines[ix+1][lineIx] != ',' {
			machine.b_x *= 10
			machine.b_x += int(lines[ix+1][lineIx] - '0')
			lineIx++
		}
		lineIx += 4
		for ; lineIx < len(lines[ix+1]); lineIx++ {
			machine.b_y *= 10
			machine.b_y += int(lines[ix+1][lineIx] - '0')
		}
		lineIx = 9
		for lines[ix+2][lineIx] != ',' {
			machine.prize_x *= 10
			machine.prize_x += int(lines[ix+2][lineIx] - '0')
			lineIx++
		}
		lineIx += 4
		for ; lineIx < len(lines[ix+2]); lineIx++ {
			machine.prize_y *= 10
			machine.prize_y += int(lines[ix+2][lineIx] - '0')
		}
		machines = append(machines, machine)
	}

	return strconv.Itoa(d13solve(machines)), machines
}

func d13solve(machines []d13machine) int {
	// Pretty trivial day tbh - the configuration
	// of each machine boils down to a pair of
	// simultaneous equations over a pair of variables,
	// which have a single unique solution.
	//
	// The algebra was done on paper, and here's the
	// result ;-)
	total := 0
	for _, machine := range machines {
		a_presses := ((machine.prize_x * machine.b_y) - (machine.prize_y * machine.b_x)) / ((machine.a_x * machine.b_y) - (machine.a_y * machine.b_x))
		b_presses := ((machine.prize_x * machine.a_y) - (machine.prize_y * machine.a_x)) / ((machine.b_x * machine.a_y) - (machine.b_y * machine.a_x))
		if a_presses*machine.a_x+b_presses*machine.b_x == machine.prize_x && a_presses*machine.a_y+b_presses*machine.b_y == machine.prize_y {
			total += b_presses + 3*a_presses
		}
	}
	return total
}

func Day13Part2(logger *slog.Logger, input string, part1Context any) string {
	machines := part1Context.([]d13machine)
	for ix := range machines {
		machines[ix].prize_x += 10000000000000
		machines[ix].prize_y += 10000000000000
	}
	return strconv.Itoa(d13solve(machines))
}
