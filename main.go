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

	r := runner.NewRunner(logger, "2024", []runner.DayImplementation{Day1, Day2, Day3})
	r.Run()
}
