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
	seenByIndex int
	seenByBit   uint32
	secret      int
	deltas      uint32
}

func newBuyer(index int, initialSecret int) d22Buyer {
	return d22Buyer{
		seenByIndex: index / 32,
		seenByBit:   1 << (index % 32),
		secret:      initialSecret}
}

func (buyer *d22Buyer) updateSecretAndDelta() int {
	newSecret := calcNextSecret(buyer.secret)
	price := newSecret % 10
	delta := price - (buyer.secret % 10)
	buyer.secret = newSecret

	// We don't care about the actual value of recent
	// deltas, we just want a unique key from them.
	// (delta + 9) is in the range 0-18 so needs only
	// 5 bits to store, so we can store the last four
	// as a 20-bit number - low enough to use as an
	// array index.
	buyer.deltas <<= 5
	buyer.deltas &= DELTA_MASK
	buyer.deltas |= uint32(delta + 9)

	return price
}

func (buyer *d22Buyer) hasSeenCurrentDelta() bool {
	// Bit of shenanigans here. We want a data structure
	// that allows us to very efficiently answer the
	// question "has buyer X seen delta sequence Y before".
	// Maps are slower than I'd like. But because we store
	// delta sequences as 20-bit values, if we can store
	// the answer to this question as a single bit, then
	// we can do so for every possible delta sequence
	// (2^20) for 3200 buyers in 400MB RAM, which is
	// reasonable. So what we have here is an array of
	// 2^20 arrays of 100 uint32s, where each buyer
	// corresponds to one bit in one of those uint32s.
	return deltasSeenBy[buyer.deltas][buyer.seenByIndex]&buyer.seenByBit != 0
}

func (buyer *d22Buyer) recordCurrentPrice(price int) {
	// Use atomic operations to allow us to update these
	// from different threads without any locking.
	atomic.OrUint32(&deltasSeenBy[buyer.deltas][buyer.seenByIndex], buyer.seenByBit)
	atomic.AddUint32(&priceForDeltas[buyer.deltas], uint32(price))
}

func (buyer *d22Buyer) generateAllSecrets(c chan int) {
	// The procedure here is that we do go through each of the
	// 2000 secret number updates, but we figure everything out
	// in that single pass.
	//
	// We maintain the last four price deltas as a single
	// 20-bit number. We can then use that as a key into
	// an effectively per-buyer record of whether this buyer
	// has seen that sequence before, and a global record of
	// what the total price all buyers pay when they see that
	// sequence is.
	for ix := range 2000 {
		price := buyer.updateSecretAndDelta()

		if ix > 2 && !buyer.hasSeenCurrentDelta() {
			// This is the first time we've seen this delta sequence for this buyer.
			// Record that we've seen it, and add the current price to the total
			// price that you get for this delta sequence.
			buyer.recordCurrentPrice(price)
		}
	}
	c <- buyer.secret
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

	// Each buyer is completely independent. So as we parse
	// each one from the input, kick off a goroutine to start
	// modelling it.
	c := make(chan int)
	for ix, line := range lines {
		secret, _ := strconv.Atoi(line)
		buyer := newBuyer(ix, secret)
		buyers = append(buyers, buyer)
		go buyer.generateAllSecrets(c)
	}

	// Collate results.
	sum := 0
	for range buyers {
		sum += <-c
	}

	return strconv.Itoa(sum), nil
}

func Day22Part2(logger *slog.Logger, input string, part1Context any) string {
	// We already calculated, during part 1, the total
	// price buyers pay for every delta they see. So
	// all we need to do now is find the highest.
	var bestResult uint32
	for deltas := range MAX_DELTAS {
		if priceForDeltas[deltas] > bestResult {
			bestResult = priceForDeltas[deltas]
		}
	}

	return strconv.Itoa(int(bestResult))
}
