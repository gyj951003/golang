// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Generating random text: a Markov chain algorithm

Based on the program presented in the "Design and Implementation" chapter
of The Practice of Programming (Kernighan and Pike, Addison-Wesley 1999).
See also Computer Recreations, Scientific American 260, 122 - 125 (1989).

A Markov chain algorithm generates text by creating a statistical model of
potential textual suffixes for a given prefix. Consider this text:

	I am not a number! I am a free man!

Our Markov chain algorithm would arrange this text into this set of prefixes
and suffixes, or "chain": (This table assumes a prefix length of two words.)

	Prefix       Suffix

	"" ""        I
	"" I         am
	I am         a
	I am         not
	a free       man!
	am a         free
	am not       a
	a number!    I
	number! I    am
	not a        number!

To generate text using this table we select an initial prefix ("I am", for
example), choose one of the suffixes associated with that prefix at random
with probability determined by the input statistics ("a"),
and then create a new prefix by removing the first word from the prefix
and appending the suffix (making the new prefix is "am a"). Repeat this process
until we can't find any suffixes for the current prefix or we exceed the word
limit. (The word limit is necessary as the chain table may contain cycles.)

Our version of this program reads text from standard input, parsing it into a
Markov chain, and writes generated text to standard output.
The prefix and output lengths can be specified using the -prefix and -words
flags on the command-line.
*/
package main

import (
	"bufio"
	//"flag"
	"fmt"
	"io"
	//"math/rand"
	"os"
	"strings"
	"strconv"
	//"time"
)

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string


// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.

type Chain struct {
	chain map[string]map[string]int
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string]map[string]int), prefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	p := make(Prefix, c.prefixLen)

	for scanner.Scan() {
		s := scanner.Text()

		key := p.String()

		if wordMap, ok := c.chain[key]; ok {
			if freq, ok := wordMap[s]; ok {
				wordMap[s] = freq + 1
			} else {
				wordMap[s] = 1
			}
			c.chain[key] = wordMap
		} else { // p not in chain
			//fmt.Println("Not In Chain")
			wordMap := make(map[string]int)
			wordMap[s] = 1
			c.chain[key] = wordMap
		}

		p.Shift(s)
	}
	//fmt.Println("Finish Building chain", len(c.chain))
}

func (c *Chain) OutputChain(r io.Writer) {
	fmt.Fprintf(r, strconv.Itoa(c.prefixLen) + "\n")

	for key := range c.chain {
		line := ""
		keys := strings.Split(key, " ")
		for _, item := range(keys) {
			if item == "" {
				line = line + "\"\"" + " "
			} else {
				line = line + item + " "
			}
		}
		line = line[:len(line)-1]

		for word, freq := range c.chain[key] {
			line = line + " " + word + " " + strconv.Itoa(freq)
		}

		line = line + "\n"


		fmt.Fprint(r, line)
	}
}

func (c *Chain) ReadChainFile(r io.Reader) {
	scanner := bufio.NewScanner(r)
	j := 0
	for scanner.Scan() {
		j++
		if j == 1 {
			continue
		}

		s := scanner.Text()
		fmt.Println(s)
		line := strings.Split(s[:len(s)], " ")
		keys := make(Prefix, 0)
		for i := 0; i < c.prefixLen; i++ {
			if line[i] == "\"\"" {
				line[i] = ""
			}
			keys = append(keys, line[i])
		}
		key := keys.String()
		wordMap := make(map[string]int, 0)

		for i:= c.prefixLen; i < len(line); i +=2 {
			value, err := strconv.Atoi(line[i+1])
			if err != nil {
				fmt.Println("Error reading frequency!")
				os.Exit(1)
			}
			wordMap[line[i]] = value
		}
		c.chain[key] = wordMap
		//fmt.Println(key, wordMap)
		//fmt.Println(key, wordMap)
	}
	fmt.Println(j)
}

// Generate retźrns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	words := make([]string, 0)

	for i := 0; i < n; i++ {
		if _,ok := c.chain[p.String()]; !ok {
			fmt.Println("stopped at ", p.String())
			break
		} else {
			maxFreq := 0
			freqWord := ""

			for word, freq := range(c.chain[p.String()]) {
				if freq >= maxFreq {
					maxFreq = freq
					freqWord = word
					}
				}

				next := freqWord
				words = append(words, next)
				p.Shift(next)
		}
	}
	fmt.Println(len(words))
	return strings.Join(words, " ")
}

func main() {
	cmd := os.Args[1]
	if cmd == "read" {
		if len(os.Args) < 5 {
			fmt.Println("Error parsing the command line!")
			os.Exit(1)
		}

		prefixLen, err1 := strconv.Atoi(os.Args[2]) // the number of words

		if err1 != nil || prefixLen <= 0 {
			fmt.Println("Error parsing the number of words!")
			os.Exit(1)
		}

		outFile, err3 := os.Create(os.Args[3])
		if err3 != nil {
			fmt.Println("Error opening the file!")
			os.Exit(1)
		}

		inFiles := os.Args[4:len(os.Args)]
		c := NewChain(prefixLen)

		for _, inFile := range(inFiles){
			file, err2 := os.Open(inFile)

			if err2 != nil {
				fmt.Println("Error opening the file!")
				os.Exit(1)
			}

			c.Build(file)
			file.Close()
		}

		c.OutputChain(outFile)
		outFile.Close()

	} else if cmd == "generate" {
		if len(os.Args) < 4 {
			fmt.Println("Error parsing the command line!")
			os.Exit(1)
		}

		modeFileName := os.Args[2]
		modeFile, err1 := os.Open(modeFileName) // the number of words

		if err1 != nil {
			fmt.Println("Error opening the file!")
			os.Exit(1)
		}

		n, err2 := strconv.Atoi(os.Args[3]) //  number of words to be generated

		if err2 != nil || n <= 0 {
			fmt.Println("Error parsing the number of words to be generated!")
			os.Exit(1)
		}

		scanner := bufio.NewScanner(modeFile)
		var prefixLen int
		for scanner.Scan() {
			fmt.Sscanf(scanner.Text(), "%d\n", &prefixLen)
			break
		}
		modeFile.Close()

		modeFileR, _:= os.Open(modeFileName)
		c := NewChain(prefixLen)
		c.ReadChainFile(modeFileR)
		fmt.Println(c.Generate(n))

		modeFileR.Close()

	} else {
		fmt.Println("Error in Command!")
		os.Exit(1)
	}
}
