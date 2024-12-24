# Advent of Code 2024 - Go
A complete set of solutions to the Advent of Code 2024 puzzles, written in Go.

I used AoC this year to learn Go - I hadn't written a line of it before December 1st. So this isn't the repository to look at if you're after beautiful examples of idiomatic Go. I also haven't got the absolute best algorithm for every day; programmers who really know their algorithms could no doubt significantly improve on what's here.

What I *do* have here is consistently Pretty Decent™ solutions. I set out with the goal of having a total cumulative runtime for all 25 days together of under a second, but I knocked that out of the park - I've actually managed a total runtime of under 70ms on my machine (Ryzen 5600X), with every single day executing under 10ms, and half under 1ms. There's also enough commenting that you should be able to follow what's going on. So hopefully this will be a useful repository for anyone wanting inspiration for what a Pretty Decent™ solution looks like for any given day.

## Framework

It also showcases the value in having a framework. If you take a glance at the code in here, you'll see that upwards of 99% is directly solving AoC problems, as everything like downloading/caching inputs or timing execution is taken care of by my framework - https://github.com/ThePants999/advent-of-code-go-runner. It's not currently robust or well-documented, so it won't impress you, but it might inspire you to create your own if you don't have one. Or, if you don't care about that robustness malarkey and you're writing Go yourself, you're free to use it.

## Running the code

If you want to execute these solutions yourself, it's pretty simple, assuming you have Go installed:

```sh
git clone git@github.com:ThePants999/advent-of-code-2024.git
cd advent-of-code-2024
go build
./advent-of-code-2024 -a
```

You'll be prompted for your session cookie so that it can download your inputs for you.

## Execution times

Averages over a thousand executions.

| Day | Median | Mean |
| ------- | ------- | ------- |
| 1 | 95µs | 114µs |
| 2 | 168µs | 198µs |
| 3 | 181µs | 208µs |
| 4 | 843µs | 769µs |
| 5 | 942µs | 727µs |
| 6 | 7.034ms | 7.336ms |
| 7 | 6.622ms | 7.475ms |
| 8 | 257µs | 276µs |
| 9 | 3.960ms | 3.664ms |
| 10 | 1.44ms | 1.471ms |
| 11 | 7.867ms | 8.201ms |
| 12 | | |
| 13 | 18µs | 26µs |
| 14 | 592µs | 684µs |
| 15 | 603µs | 704µs |
| 16 | 6.994ms | 5.489ms |
| 17 | 27µs | 32µs |
| 18 | 543µs | 696µs |
| 19 | 2.416ms | 2.66ms |
| 20 | 2.574ms | 2.912ms |
| 21 | 25µs | 24µs |
| 22 | 7.107ms | 9.096ms |
| 23 | 2.076ms | 2.066ms |
| 24 | 80µs | 107µs |