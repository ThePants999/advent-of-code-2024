package main

import (
	"log/slog"
	"os"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

func main() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelWarn)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl}))

	r := runner.NewRunner(logger, "2024", []runner.DayImplementation{Day1, Day2, Day3, Day4, Day5, Day6, Day7, Day8, Day9, Day10, Day11, Day12})
	r.Run()
}
