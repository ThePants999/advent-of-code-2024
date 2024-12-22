package main

import (
	"log/slog"
	"strconv"
	"strings"
	"sync/atomic"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day22 = runner.DayImplementation{
	DayNumber:    22,
	ExecutePart1: Day22Part1,
	ExecutePart2: Day22Part2,
	ExampleInput: `1
2
3
2024`,
	ExamplePart1Answer: "37990510",
	ExamplePart2Answer: "23",
}

const PRUNE_BITS int = 0b111111111111111111111111

const MAX_DELTAS int = 1 << 20
const DELTA_MASK uint32 = 0b11111111111111111111

var deltasSeenBy [MAX_DELTAS][100]uint32
var priceForDeltas [MAX_DELTAS]uint32

type d22Buyer struct {
	seenByIndex    int
	seenByBit      uint32
	originalSecret int
	currentSecret  int
	deltas         uint32
}

func newBuyer(index int, initialSecret int) d22Buyer {
	return d22Buyer{
		seenByIndex:    index / 32,
		seenByBit:      1 << (index % 32),
		originalSecret: initialSecret}
}

func (buyer *d22Buyer) generateAllSecrets(c chan int) {
	buyer.currentSecret = buyer.originalSecret
	for ix := range 2000 {
		newSecret := calcNextSecret(buyer.currentSecret)
		price := newSecret % 10
		delta := price - (buyer.currentSecret % 10)

		// We don't care about the actual value of recent
		// deltas, we just want a unique key from them.
		// (delta + 9) is in the range 0-18 so needs only
		// 5 bits to store, so we can store the last four
		// as a 20-bit number - low enough to use as an
		// array index.
		buyer.deltas <<= 5
		buyer.deltas &= DELTA_MASK
		buyer.deltas |= uint32(delta + 9)

		if ix > 2 {
			if deltasSeenBy[buyer.deltas][buyer.seenByIndex]&buyer.seenByBit == 0 {
				// This is the first time we've seen this delta sequence for this buyer.
				// Record that we've seen it, and add the current price to the total
				// price that you get for this delta sequence.
				atomic.OrUint32(&deltasSeenBy[buyer.deltas][buyer.seenByIndex], buyer.seenByBit)
				atomic.AddUint32(&priceForDeltas[buyer.deltas], uint32(price))
			}
		}
		buyer.currentSecret = newSecret
	}
	c <- buyer.currentSecret
}

func calcNextSecret(secret int) int {
	secret = ((secret << 6) ^ secret) & PRUNE_BITS
	secret = ((secret >> 5) ^ secret)
	secret = ((secret << 11) ^ secret) & PRUNE_BITS
	return secret
}

func Day22Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	buyers := make([]d22Buyer, 0, len(lines))

	c := make(chan int)
	for ix, line := range lines {
		secret, _ := strconv.Atoi(line)
		buyer := newBuyer(ix, secret)
		buyers = append(buyers, buyer)
		go buyer.generateAllSecrets(c)
	}

	sum := 0
	for range buyers {
		sum += <-c
	}

	return strconv.Itoa(sum), nil
}

func Day22Part2(logger *slog.Logger, input string, part1Context any) string {
	var bestResult uint32
	for deltas := range MAX_DELTAS {
		if priceForDeltas[deltas] > bestResult {
			bestResult = priceForDeltas[deltas]
		}
	}

	return strconv.Itoa(int(bestResult))
}
