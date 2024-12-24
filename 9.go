package main

import (
	"log/slog"
	"slices"
	"strconv"

	runner "github.com/ThePants999/advent-of-code-go-runner"
)

var Day9 = runner.DayImplementation{
	DayNumber:          9,
	ExecutePart1:       Day9Part1,
	ExecutePart2:       Day9Part2,
	ExampleInput:       "2333133121414131402",
	ExamplePart1Answer: "1928",
	ExamplePart2Answer: "2858",
}

// Terrifically efficient part 1 implementation in both memory and processing.
// Pity it's 100% useless for part 2.

func Day9Part1(logger *slog.Logger, input string) (string, any) {
	files := make([]int, 0, len(input)/2)
	gaps := make([]int, 0, len(input)/2)
	for ix, char := range input {
		length := int(char - '0')
		if ix%2 == 0 {
			files = append(files, length)
		} else {
			gaps = append(gaps, length)
		}
	}

	ix, left_file_ix, right_file_ix, right_file_subix, gap_ix := 0, 0, len(files)-1, 0, 0
	checksum := 0
	for {
		// Fully append the next file on the left.
		for i := 0; i < files[left_file_ix]; i++ {
			checksum += (left_file_ix * ix)
			ix++
		}
		left_file_ix++

		// If the next file on the right is the one we
		// just appended, we're done.
		if left_file_ix > right_file_ix {
			break
		}

		if right_file_ix == left_file_ix {
			// The current right-hand file is the last
			// file. Finish appending it and then we're
			// done.
			for ; right_file_subix < files[right_file_ix]; right_file_subix++ {
				checksum += (right_file_ix * ix)
				ix++
			}
			break
		}

		// Fill the next gap with right-hand files.
		for i := 0; i < gaps[gap_ix]; i++ {
			checksum += (right_file_ix * ix)
			ix++
			right_file_subix++
			if right_file_subix == files[right_file_ix] {
				right_file_ix--
				right_file_subix = 0
			}
		}
		gap_ix++
	}

	return strconv.Itoa(checksum), nil
}

type diskElement struct {
	disk *disk
	prev *diskElement
	next *diskElement
	file bool
	id   int
	len  int
	pos  int
}

func (element *diskElement) addAtEnd(disk *disk) {
	element.disk = disk
	element.prev = disk.last
	if disk.last == nil {
		disk.first = element
	} else {
		disk.last.next = element
	}
	disk.last = element
}

func (element *diskElement) insertAfter(prev *diskElement) {
	element.next = prev.next
	if prev.next != nil {
		prev.next.prev = element
	} else {
		prev.disk.last = element
	}
	prev.next = element
	element.prev = prev
}

func (element *diskElement) replaceWithGap() {
	gap := diskElement{}
	gap.disk = element.disk
	gap.len = element.len
	gap.pos = element.pos

	// We know we won't call this function on the first element so
	// I'm gonna be lazy
	prev := element.prev
	element.remove()
	gap.insertAfter(prev)
}

func (element *diskElement) remove() {
	if element.prev != nil {
		element.prev.next = element.next
	}
	if element.next != nil {
		element.next.prev = element.prev
	}
	if element.disk.first == element {
		element.disk.first = element.next
	}
	if element.disk.last == element {
		element.disk.last = element.prev
	}
}

type disk struct {
	first *diskElement
	last  *diskElement
}

func Day9Part2(logger *slog.Logger, input string, part1Context any) string {
	gaps := make([][]*diskElement, 10)
	for size := 1; size < 10; size++ {
		gaps[size] = make([]*diskElement, 0, len(input)/5)
	}

	disk := disk{}
	files := make([]*diskElement, 0, len(input)/2)
	pos := 0

	for ix, char := range input {
		len := int(char - '0')
		if len > 0 {
			element := diskElement{}
			element.addAtEnd(&disk)

			element.len = len
			if ix%2 == 0 {
				element.file = true
				element.id = ix / 2
				files = append(files, &element)
			} else {
				for gapSize := 1; gapSize <= element.len; gapSize++ {
					gaps[gapSize] = append(gaps[gapSize], &element)
				}
			}

			element.pos = pos
			pos += element.len
		}
	}

	// Perform compaction.
	for file_ix := len(files) - 1; file_ix >= 0; file_ix-- {
		file := files[file_ix]
		if len(gaps[file.len]) > 0 {
			// There's a gap that can fit this file.
			gap := gaps[file.len][0]
			if gap.pos < file.pos {
				// The gap is left of the file. Move the file into
				// the gap.
				//
				// Firstly, replace the file with a new gap. We're not
				// going to bother doing anything clever in terms of
				// merging with gaps before and after because the "try
				// each file once" nature of the problem means that
				// we're never going to try to move anything into this
				// new gap, so we only need it for the purposes of
				// checksum calculation.
				file.replaceWithGap()

				file.insertAfter(gap.prev)
				newSize := gap.len - file.len
				shrinkGap(gap, gaps, newSize)
				if newSize == 0 {
					gap.remove()
				}
			}
		}
	}

	ix, checksum := 0, 0
	for element := disk.first; element != nil; element = element.next {
		if element.file {
			for ; element.len > 0; element.len-- {
				checksum += (element.id * ix)
				ix++
			}
		} else {
			ix += element.len
		}
	}

	return strconv.Itoa(checksum)
}

func shrinkGap(gap *diskElement, gaps [][]*diskElement, newSize int) {
	// Remove this gap from the lists of gaps larger than its new size.
	for gapLen := gap.len; gapLen > newSize; gapLen-- {
		gapIx := slices.Index(gaps[gapLen], gap)
		if gapIx < len(gaps[gapLen])/2 {
			// This gap is less than half way along the list - it's
			// less work to shift preceding gaps forwards.
			for ix := gapIx; ix > 0; ix-- {
				gaps[gapLen][ix] = gaps[gapLen][ix-1]
			}
			gaps[gapLen] = gaps[gapLen][1:]
		} else {
			// This is in the second half of gaps at this size - easier
			// to shift the rest of the slice backwards.
			gaps[gapLen] = slices.Delete(gaps[gapLen], gapIx, gapIx+1)
		}
	}
	gap.pos += (gap.len - newSize)
	gap.len = newSize
}
