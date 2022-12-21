package main

import (
	"bytes"
	"fmt"
	"log"
	"sort"

	"golang.org/x/tools/cover"
)

// Packages stores a collection of coverage profiles.
type Packages struct {
	profiles []*cover.Profile
}

// Add appends or merges a profile to the collection.
func (p *Packages) Add(profile *cover.Profile) {
	i := sort.Search(
		len(p.profiles),
		func(i int) bool {
			return p.profiles[i].FileName >= profile.FileName
		},
	)

	if i < len(p.profiles) && p.profiles[i].FileName == profile.FileName {
		p.merge(p.profiles[i], profile)
	} else {
		p.profiles = append(p.profiles, nil)
		copy(p.profiles[i+1:], p.profiles[i:])
		p.profiles[i] = profile
	}
}

// Dump returns a generated merged profile collection.
func (p *Packages) Dump() []byte {
	buffer := bytes.NewBufferString("")

	if len(p.profiles) == 0 {
		return buffer.Bytes()
	}

	fmt.Fprintf(
		buffer,
		"mode: %s\n",
		p.profiles[0].Mode,
	)

	for _, profile := range p.profiles {
		for _, block := range profile.Blocks {
			fmt.Fprintf(
				buffer,
				"%s:%d.%d,%d.%d %d %d\n",
				profile.FileName,
				block.StartLine,
				block.StartCol,
				block.EndLine,
				block.EndCol,
				block.NumStmt,
				block.Count,
			)
		}
	}

	return buffer.Bytes()
}

// merge simply merges all coverage blocks as required.
func (p *Packages) merge(initial, addition *cover.Profile) {
	start := 0

	for _, block := range addition.Blocks {
		sortFunc := func(inner int) bool {
			return initial.Blocks[inner+start].StartLine >= block.StartLine && (initial.Blocks[inner+start].StartLine != block.StartLine || initial.Blocks[inner+start].StartCol >= block.StartCol)
		}

		i := 0

		if !sortFunc(i) {
			i = sort.Search(len(initial.Blocks)-start, sortFunc)
		}

		i += start

		if i < len(initial.Blocks) && initial.Blocks[i].StartLine == block.StartLine && initial.Blocks[i].StartCol == block.StartCol {
			if initial.Blocks[i].EndLine != block.EndLine || initial.Blocks[i].EndCol != block.EndCol {
				log.Fatalf("overlap merge: %v %v %v", initial.FileName, initial.Blocks[i], block)
			}

			switch initial.Mode {
			case "set":
				initial.Blocks[i].Count |= block.Count
			case "count", "atomic":
				initial.Blocks[i].Count += block.Count
			default:
				log.Fatalf("unsupported covermode: %s", initial.Mode)
			}
		} else {
			if i > 0 {
				if initial.Blocks[i-1].EndLine >= block.EndLine && (initial.Blocks[i-1].EndLine != block.EndLine || initial.Blocks[i-1].EndCol > block.EndCol) {
					log.Fatalf("overlap merge: %v %v %v", initial.FileName, initial.Blocks[i-1], block)
				}
			}

			if i < len(initial.Blocks)-1 {
				if initial.Blocks[i+1].StartLine <= block.StartLine && (initial.Blocks[i+1].StartLine != block.StartLine || initial.Blocks[i+1].StartCol < block.StartCol) {
					log.Fatalf("overlap after: %v %v %v", initial.FileName, initial.Blocks[i+1], block)
				}
			}

			initial.Blocks = append(
				initial.Blocks,
				cover.ProfileBlock{},
			)

			copy(
				initial.Blocks[i+1:],
				initial.Blocks[i:],
			)

			initial.Blocks[i] = block
		}

		start = i + 1
	}
}
