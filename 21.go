package main

import (
	"log/slog"
	"strconv"
	"strings"
	"sync"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day21 = runner.DayImplementation{
	DayNumber:    21,
	ExecutePart1: Day21Part1,
	ExecutePart2: Day21Part2,
	ExampleInput: `029A
980A
179A
456A
379A`,
	ExamplePart1Answer: "126384",
	ExamplePart2Answer: "154115708116294",
}

type d21DirKeypadButton int

const (
	D21_UP d21DirKeypadButton = iota
	D21_LEFT
	D21_DOWN
	D21_RIGHT
	D21_PRESS
)

var directionDirections [5][5][]d21DirKeypadButton = [5][5][]d21DirKeypadButton{
	// From UP
	{
		{D21_PRESS},                      // UP
		{D21_DOWN, D21_LEFT, D21_PRESS},  // LEFT
		{D21_DOWN, D21_PRESS},            // DOWN
		{D21_DOWN, D21_RIGHT, D21_PRESS}, // RIGHT
		{D21_RIGHT, D21_PRESS},           // PRESS
	},
	// From LEFT
	{
		{D21_RIGHT, D21_UP, D21_PRESS},            // UP
		{D21_PRESS},                               // LEFT
		{D21_RIGHT, D21_PRESS},                    // DOWN
		{D21_RIGHT, D21_RIGHT, D21_PRESS},         // RIGHT
		{D21_RIGHT, D21_RIGHT, D21_UP, D21_PRESS}, // PRESS
	},
	// From DOWN
	{
		{D21_UP, D21_PRESS},            // UP
		{D21_LEFT, D21_PRESS},          // LEFT
		{D21_PRESS},                    // DOWN
		{D21_RIGHT, D21_PRESS},         // RIGHT
		{D21_UP, D21_RIGHT, D21_PRESS}, // PRESS
	},
	// From RIGHT
	{
		{D21_LEFT, D21_UP, D21_PRESS},   // UP
		{D21_LEFT, D21_LEFT, D21_PRESS}, // LEFT
		{D21_LEFT, D21_PRESS},           // DOWN
		{D21_PRESS},                     // RIGHT
		{D21_UP, D21_PRESS},             // PRESS
	},
	// From PRESS
	{
		{D21_LEFT, D21_PRESS},                     // UP
		{D21_DOWN, D21_LEFT, D21_LEFT, D21_PRESS}, // LEFT
		{D21_LEFT, D21_DOWN, D21_PRESS},           // DOWN
		{D21_DOWN, D21_PRESS},                     // RIGHT
		{D21_PRESS},                               // PRESS
	},
}

type numKeypadPress struct {
	button     d21DirKeypadButton
	numPresses int
}

const D21_PRESS_NUM int = 10

var numericDirections [11][11][]numKeypadPress = [11][11][]numKeypadPress{
	// From 0
	{
		{}, // 0
		{numKeypadPress{D21_UP, 1}, numKeypadPress{D21_LEFT, 1}},  // 1
		{numKeypadPress{D21_UP, 1}},                               // 2
		{numKeypadPress{D21_UP, 1}, numKeypadPress{D21_RIGHT, 1}}, // 3
		{numKeypadPress{D21_UP, 2}, numKeypadPress{D21_LEFT, 1}},  // 4
		{numKeypadPress{D21_UP, 2}},                               // 5
		{numKeypadPress{D21_UP, 2}, numKeypadPress{D21_RIGHT, 1}}, // 6
		{numKeypadPress{D21_UP, 3}, numKeypadPress{D21_LEFT, 1}},  // 7
		{numKeypadPress{D21_UP, 3}},                               // 8
		{numKeypadPress{D21_UP, 3}, numKeypadPress{D21_RIGHT, 1}}, // 9
		{numKeypadPress{D21_RIGHT, 1}},                            // PRESS
	},
	// From 1
	{
		{numKeypadPress{D21_RIGHT, 1}, numKeypadPress{D21_DOWN, 1}}, // 0
		{},                             // 1
		{numKeypadPress{D21_RIGHT, 1}}, // 2
		{numKeypadPress{D21_RIGHT, 2}}, // 3
		{numKeypadPress{D21_UP, 1}},    // 4
		{numKeypadPress{D21_UP, 1}, numKeypadPress{D21_RIGHT, 1}},   // 5
		{numKeypadPress{D21_UP, 1}, numKeypadPress{D21_RIGHT, 2}},   // 6
		{numKeypadPress{D21_UP, 2}},                                 // 7
		{numKeypadPress{D21_UP, 2}, numKeypadPress{D21_RIGHT, 1}},   // 8
		{numKeypadPress{D21_UP, 2}, numKeypadPress{D21_RIGHT, 2}},   // 9
		{numKeypadPress{D21_RIGHT, 2}, numKeypadPress{D21_DOWN, 1}}, // PRESS
	},
	// From 2
	{
		{numKeypadPress{D21_DOWN, 1}},  // 0
		{numKeypadPress{D21_LEFT, 1}},  // 1
		{},                             // 2
		{numKeypadPress{D21_RIGHT, 1}}, // 3
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_UP, 1}},    // 4
		{numKeypadPress{D21_UP, 1}},                                 // 5
		{numKeypadPress{D21_UP, 1}, numKeypadPress{D21_RIGHT, 1}},   // 6
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_UP, 2}},    // 7
		{numKeypadPress{D21_UP, 2}},                                 // 8
		{numKeypadPress{D21_UP, 2}, numKeypadPress{D21_RIGHT, 1}},   // 9
		{numKeypadPress{D21_DOWN, 1}, numKeypadPress{D21_RIGHT, 1}}, // PRESS
	},
	// From 3
	{
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_DOWN, 1}}, // 0
		{numKeypadPress{D21_LEFT, 2}},                              // 1
		{numKeypadPress{D21_LEFT, 1}},                              // 2
		{},                                                         // 3
		{numKeypadPress{D21_LEFT, 2}, numKeypadPress{D21_UP, 1}}, // 4
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_UP, 1}}, // 5
		{numKeypadPress{D21_UP, 1}},                              // 6
		{numKeypadPress{D21_LEFT, 2}, numKeypadPress{D21_UP, 2}}, // 7
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_UP, 2}}, // 8
		{numKeypadPress{D21_UP, 2}},                              // 9
		{numKeypadPress{D21_DOWN, 1}},                            // PRESS
	},
	// From 4
	{
		{numKeypadPress{D21_RIGHT, 1}, numKeypadPress{D21_DOWN, 2}}, // 0
		{numKeypadPress{D21_DOWN, 1}},                               // 1
		{numKeypadPress{D21_DOWN, 1}, numKeypadPress{D21_RIGHT, 1}}, // 2
		{numKeypadPress{D21_DOWN, 1}, numKeypadPress{D21_RIGHT, 2}}, // 3
		{},                             // 4
		{numKeypadPress{D21_RIGHT, 1}}, // 5
		{numKeypadPress{D21_RIGHT, 2}}, // 6
		{numKeypadPress{D21_UP, 1}},    // 7
		{numKeypadPress{D21_UP, 1}, numKeypadPress{D21_RIGHT, 1}},   // 8
		{numKeypadPress{D21_UP, 1}, numKeypadPress{D21_RIGHT, 2}},   // 9
		{numKeypadPress{D21_RIGHT, 2}, numKeypadPress{D21_DOWN, 2}}, // PRESS
	},
	// From 5
	{
		{numKeypadPress{D21_DOWN, 2}},                               // 0
		{numKeypadPress{D21_DOWN, 1}, numKeypadPress{D21_LEFT, 1}},  // 1
		{numKeypadPress{D21_DOWN, 1}},                               // 2
		{numKeypadPress{D21_DOWN, 1}, numKeypadPress{D21_RIGHT, 1}}, // 3
		{numKeypadPress{D21_LEFT, 1}},                               // 4
		{},                                                          // 5
		{numKeypadPress{D21_RIGHT, 1}},                              // 6
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_UP, 1}},    // 7
		{numKeypadPress{D21_UP, 1}},                                 // 8
		{numKeypadPress{D21_UP, 1}, numKeypadPress{D21_RIGHT, 1}},   // 9
		{numKeypadPress{D21_RIGHT, 1}, numKeypadPress{D21_DOWN, 2}}, // PRESS
	},
	// From 6
	{
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_DOWN, 2}}, // 0
		{numKeypadPress{D21_DOWN, 1}, numKeypadPress{D21_LEFT, 2}}, // 1
		{numKeypadPress{D21_DOWN, 1}, numKeypadPress{D21_LEFT, 1}}, // 2
		{numKeypadPress{D21_DOWN, 1}},                              // 3
		{numKeypadPress{D21_LEFT, 2}},                              // 4
		{numKeypadPress{D21_LEFT, 1}},                              // 5
		{},                                                         // 6
		{numKeypadPress{D21_LEFT, 2}, numKeypadPress{D21_UP, 1}},   // 7
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_UP, 1}},   // 8
		{numKeypadPress{D21_UP, 1}},                                // 9
		{numKeypadPress{D21_DOWN, 2}},                              // PRESS
	},
	// From 7
	{
		{numKeypadPress{D21_RIGHT, 1}, numKeypadPress{D21_DOWN, 3}}, // 0
		{numKeypadPress{D21_DOWN, 2}},                               // 1
		{numKeypadPress{D21_DOWN, 2}, numKeypadPress{D21_RIGHT, 1}}, // 2
		{numKeypadPress{D21_DOWN, 2}, numKeypadPress{D21_RIGHT, 2}}, // 3
		{numKeypadPress{D21_DOWN, 1}},                               // 4
		{numKeypadPress{D21_DOWN, 1}, numKeypadPress{D21_RIGHT, 1}}, // 5
		{numKeypadPress{D21_DOWN, 1}, numKeypadPress{D21_RIGHT, 2}}, // 6
		{},                             // 7
		{numKeypadPress{D21_RIGHT, 1}}, // 8
		{numKeypadPress{D21_RIGHT, 2}}, // 9
		{numKeypadPress{D21_RIGHT, 2}, numKeypadPress{D21_DOWN, 3}}, // PRESS
	},
	// From 8
	{
		{numKeypadPress{D21_DOWN, 3}},                               // 0
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_DOWN, 2}},  // 1
		{numKeypadPress{D21_DOWN, 2}},                               // 2
		{numKeypadPress{D21_DOWN, 2}, numKeypadPress{D21_RIGHT, 1}}, // 3
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_DOWN, 1}},  // 4
		{numKeypadPress{D21_DOWN, 1}},                               // 5
		{numKeypadPress{D21_DOWN, 1}, numKeypadPress{D21_RIGHT, 1}}, // 6
		{numKeypadPress{D21_LEFT, 1}},                               // 7
		{},                                                          // 8
		{numKeypadPress{D21_RIGHT, 1}},                              // 9
		{numKeypadPress{D21_DOWN, 3}, numKeypadPress{D21_RIGHT, 1}}, // PRESS
	},
	// From 9
	{
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_DOWN, 3}}, // 0
		{numKeypadPress{D21_LEFT, 2}, numKeypadPress{D21_DOWN, 2}}, // 1
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_DOWN, 2}}, // 2
		{numKeypadPress{D21_DOWN, 2}},                              // 3
		{numKeypadPress{D21_LEFT, 2}, numKeypadPress{D21_DOWN, 1}}, // 4
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_DOWN, 1}}, // 5
		{numKeypadPress{D21_DOWN, 1}},                              // 6
		{numKeypadPress{D21_LEFT, 2}},                              // 7
		{numKeypadPress{D21_LEFT, 1}},                              // 8
		{},                                                         // 9
		{numKeypadPress{D21_DOWN, 3}},                              // PRESS
	},
	// From PRESS
	{
		{numKeypadPress{D21_LEFT, 1}},                            // 0
		{numKeypadPress{D21_UP, 1}, numKeypadPress{D21_LEFT, 2}}, // 1
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_UP, 1}}, // 2
		{numKeypadPress{D21_UP, 1}},                              // 3
		{numKeypadPress{D21_UP, 2}, numKeypadPress{D21_LEFT, 2}}, // 4
		{numKeypadPress{D21_LEFT, 1}, numKeypadPress{D21_UP, 2}}, // 5
		{numKeypadPress{D21_UP, 2}},                              // 6
		{numKeypadPress{D21_UP, 3}, numKeypadPress{D21_LEFT, 2}}, // 7
		{numKeypadPress{D21_UP, 3}, numKeypadPress{D21_LEFT, 1}}, // 8
		{numKeypadPress{D21_UP, 3}},                              // 9
		{},                                                       // PRESS
	},
}

func findCodeComplexity(codeStr string, dirsLevels int, c chan int) {
	code := make([]int, len(codeStr)+1)
	code[0] = D21_PRESS_NUM
	code[len(code)-1] = D21_PRESS_NUM
	number := 0
	for ix := 0; ix < len(codeStr)-1; ix++ {
		code[ix+1] = (int)(codeStr[ix]) - '0'
		number *= 10
		number += code[ix+1]
	}

	sequenceLen := convertNumsToDirs(code, dirsLevels)

	c <- number * sequenceLen
}

func dirSeqToString(seq []d21DirKeypadButton) string {
	runes := make([]rune, len(seq))
	for ix, dir := range seq {
		switch dir {
		case D21_UP:
			runes[ix] = '^'
		case D21_DOWN:
			runes[ix] = 'v'
		case D21_LEFT:
			runes[ix] = '<'
		case D21_RIGHT:
			runes[ix] = '>'
		case D21_PRESS:
			runes[ix] = 'A'
		}
	}
	return string(runes)
}

func convertNumsToDirs(nums []int, dirsLevels int) int {
	// 6 is max presses to get from one num to another
	sequence := make([]d21DirKeypadButton, 1, len(nums)*6+1)

	// We always put PRESS at the start of a sequence because sequence
	// conversion depends on knowing what the previous button was, and
	// the robot arm starts off positioned facing PRESS.
	sequence[0] = D21_PRESS

	// Correspondingly, we don't convert the first number given to us,
	// as it's just telling us a starting point.
	for ix := 1; ix < len(nums); ix++ {
		presses := numericDirections[nums[ix-1]][nums[ix]]
		for _, press := range presses {
			for range press.numPresses {
				sequence = append(sequence, press.button)
			}
		}
		sequence = append(sequence, D21_PRESS)
	}

	return convertDirsToDirs(sequence, dirsLevels)
}

type d21CacheKey struct {
	dirs            string
	remainingLevels int
}

var d21Cache sync.Map = sync.Map{}

func convertDirsToDirs(dirs []d21DirKeypadButton, remainingLevels int) int {
	if remainingLevels == 0 {
		return len(dirs) - 1
	}

	// See if we've already answered this question.
	dirsStr := dirSeqToString(dirs)
	key := d21CacheKey{dirsStr, remainingLevels}
	val, found := d21Cache.Load(key)
	if found {
		return val.(int)
	}

	// This code will break if dirs doesn't both start and end with D21_PRESS,
	// but that should always be the case.
	totalLen := 0
	startIx := 0
	for startIx < len(dirs)-1 {
		// We're going to work in chunks that start and end with a PRESS.
		// startIx should be a PRESS, so find the next one.
		endIx := startIx + 1
		for dirs[endIx] != D21_PRESS {
			endIx++
		}

		// Convert this chunk to what the preceding robot needs to input.

		// It takes max 5 keypresses to get from one direction to another, so
		// max length of converted sequence is 5 times unconverted sequence
		// plus an opening PRESS.
		sequence := make([]d21DirKeypadButton, 1, (endIx-startIx)*5+1)
		sequence[0] = D21_PRESS
		for ; startIx < endIx; startIx++ {
			// Add to the sequence
			sequence = append(sequence, directionDirections[dirs[startIx]][dirs[startIx+1]]...)
		}

		// Pass this input to the next robot up the chain; it'll return the
		// length of the sequence at the end of the chain.
		totalLen += convertDirsToDirs(sequence, remainingLevels-1)
	}

	// Cache this result.
	d21Cache.Store(key, totalLen)
	return totalLen
}

func Day21Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	c := make(chan int)
	sum := 0
	for _, line := range lines {
		go findCodeComplexity(line, 2, c)
	}
	for range lines {
		sum += <-c
	}
	return strconv.Itoa(sum), lines
}

func Day21Part2(logger *slog.Logger, input string, part1Context any) string {
	lines := part1Context.([]string)
	c := make(chan int)
	sum := 0
	for _, line := range lines {
		go findCodeComplexity(line, 25, c)
	}
	for range lines {
		sum += <-c
	}
	return strconv.Itoa(sum)
}
