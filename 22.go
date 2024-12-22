package main

import (
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day22 = runner.DayImplementation{
	DayNumber:    22,
	ExecutePart1: Day22Part1,
	ExecutePart2: Day22Part2,
	ExampleInput: `1
10
100
2024`,
	ExamplePart1Answer: "37327623",
	ExamplePart2Answer: "23",
}

const PRUNE_BITS int = 0b111111111111111111111111

var freqMap sync.Map = sync.Map{}

type d22Buyer struct {
	originalSecret int
	currentSecret  int
	deltas         uint32
	seqToPrice     map[uint32]int8
}

func newBuyer(initialSecret int) d22Buyer {
	return d22Buyer{originalSecret: initialSecret, seqToPrice: make(map[uint32]int8)}
}

func (buyer *d22Buyer) generateAllSecrets(c chan int) {
	buyer.currentSecret = buyer.originalSecret
	for ix := range 2000 {
		newSecret := calcNextSecret(buyer.currentSecret)
		price := (int8)(newSecret % 10)
		delta := price - (int8)(buyer.currentSecret%10)
		buyer.deltas <<= 8
		buyer.deltas |= uint32(uint8(delta))
		if ix > 2 {
			_, found := buyer.seqToPrice[buyer.deltas]
			if !found {
				buyer.seqToPrice[buyer.deltas] = price
			}
			if price == 9 {
				var newCount uint32 = 1
				count, loaded := freqMap.LoadOrStore(buyer.deltas, &newCount)
				if loaded {
					atomic.AddUint32(count.(*uint32), 1)
				}
			}
		}
		buyer.currentSecret = newSecret
	}
	c <- buyer.currentSecret
}

func (buyer *d22Buyer) getPrice(key uint32) int {
	return int(buyer.seqToPrice[key])
}

func calcNextSecret(secret int) int {
	secret = ((secret << 6) ^ secret) & PRUNE_BITS
	secret = ((secret >> 5) ^ secret) & PRUNE_BITS
	secret = ((secret << 11) ^ secret) & PRUNE_BITS
	return secret
}

func Day22Part1(logger *slog.Logger, input string) (string, any) {
	lines := strings.Fields(input)
	buyers := make([]d22Buyer, 0, len(lines))

	c := make(chan int)
	for _, line := range lines {
		secret, _ := strconv.Atoi(line)
		buyer := newBuyer(secret)
		buyers = append(buyers, buyer)
		go buyer.generateAllSecrets(c)
	}

	sum := 0
	for range buyers {
		sum += <-c
	}

	return strconv.Itoa(sum), buyers
}

func Day22Part2(logger *slog.Logger, input string, part1Context any) string {
	buyers := part1Context.([]d22Buyer)
	numThreads := 0
	c := make(chan int)
	freqMap.Range(func(key any, _ any) bool {
		deltas := key.(uint32)
		numThreads++
		go func() {
			result := 0
			for _, buyer := range buyers {
				result += buyer.getPrice(deltas)
			}
			c <- result
		}()
		return true
	})

	bestResult := 0
	for range numThreads {
		result := <-c
		if result > bestResult {
			bestResult = result
		}
	}

	return strconv.Itoa(bestResult)
}
