// Copyright 2018 Scott Mansfield
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parse

import (
	"bufio"
	"compress/gzip"
	"io"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/ScottMansfield/markov/graph"
)

// Parse reads an input stream word by word and adds the appropriate relations to the given graph.Markov
func Parse(in io.Reader, mg graph.Markov) {
	r := bufio.NewReader(in)
	prev := ""

	for {
		word, err := r.ReadString(byte(' '))
		if err == nil {
			break
		}

		// normalize words
		word = normalize(word)
		if word == "" {
			continue
		}

		mg.IncRelation(prev, word)

		prev = word
	}
}

// GoogleNgrams parses a google ngram version 2 file
// http://storage.googleapis.com/books/ngrams/books/datasetsv2.html
func GoogleNgrams(in io.Reader, mg *graph.Markov) {
	// ngram TAB year TAB match_count TAB volume_count NEWLINE
	// In this case, we take the ngram and split it into separate words
	// For each pair of words (in order in the ngram) add the match_count to the graph

	gzr, err := gzip.NewReader(in)
	if err != nil {
		panic(err)
	}
	r := bufio.NewScanner(gzr)
	i := 0

	rawchan := make(chan string)

	// Line parsing and map writing goroutines
	pwg := &sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		pwg.Add(1)
		go func(rawchan chan string, wg *sync.WaitGroup, mg *graph.Markov) {
			for line := range rawchan {
				parts := strings.Split(line, "\t")

				ngparts := strings.Split(parts[0], " ")
				for i := range ngparts {
					ngparts[i] = normalize(ngparts[i])
				}

				matchcount, err := strconv.ParseUint(parts[2], 10, 64)
				if err != nil {
					panic(err)
				}

				for i := 1; i < len(ngparts); i++ {
					if ngparts[i-1] == "" || ngparts[i] == "" {
						continue
					}

					mg.IncRelationBy(ngparts[i-1], ngparts[i], matchcount)
				}
			}

			wg.Done()
		}(rawchan, pwg, mg)
	}

	// Scanner component
	for r.Scan() {
		i++
		if i%100000 == 0 {
			log.Println(i)
		}

		rawchan <- r.Text()
	}

	close(rawchan)
	pwg.Wait()
}

var stripsuffixes = []string{"_NOUN", "_CONJ", "_ADV", "_PRON", "_VERB", "_ADP", "_DET", "_ADJ", "_NUM"}
var skipsuffixes = []string{"_END_", "_PRT"}

// normalize normalizes the word by keeping only the characters (and hyphens) and lowercasing everything
func normalize(word string) string {
	for _, suf := range skipsuffixes {
		if strings.HasSuffix(word, suf) {
			return ""
		}
	}

	for _, suf := range stripsuffixes {
		if strings.HasSuffix(word, suf) {
			word = strings.TrimSuffix(word, suf)
			break
		}
	}

	norm := make([]rune, 0, len(word))

	// first pass: lowercase and filter
	for _, char := range word {
		if char >= 'A' && char <= 'Z' {
			norm = append(norm, char+32)

		} else if (char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') ||
			char == '-' {
			norm = append(norm, char)
		}
	}

	return string(norm)
}
