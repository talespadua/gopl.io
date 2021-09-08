// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 10.
//!+

// Dup2 prints the count and text of lines that appear more than once
// in the input.  It reads from stdin or from a list of named files.
package main

import (
	"bufio"
	"fmt"
	"os"
)

type duplicateInstances struct {
	count int
	files []string
}

func main() {
	counts := make(map[string]duplicateInstances)

	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for dupString, dupInstance := range counts {
		if dupInstance.count > 1 {
			fmt.Printf("%d\t%s\t%v\n", dupInstance.count, dupString, dupInstance.files)
		}
	}
}

func countLines(f *os.File, counts map[string]duplicateInstances) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		if val, ok := counts[input.Text()]; ok {
			val.count++
			if !val.fileInSlice(f.Name()) {
				val.files = append(val.files, f.Name())
			}

			counts[input.Text()] = val
		} else {
			var newCount duplicateInstances
			newCount.count = 1
			newCount.files = append(newCount.files, f.Name())
			counts[input.Text()] = newCount
		}
	}
	// NOTE: ignoring potential errors from input.Err()
}

func (d *duplicateInstances) fileInSlice(fileName string) bool {
	for _, fn := range d.files {
		if fn == fileName {
			return true
		}
	}
	return false
}

//!-
