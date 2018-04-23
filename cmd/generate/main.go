// Copyright 2016 Scott Mansfield
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

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/ScottMansfield/markov/graph"
)

func main() {
	var infilename, cpuprofile string
	var numwalks, length int

	flag.StringVar(&infilename, "if", "", "Serialized graph file")
	flag.IntVar(&numwalks, "n", 1, "Number of random walks to generate")
	flag.IntVar(&length, "l", 7, "Length of each random walk")

	flag.StringVar(&cpuprofile, "cpuprofile", "", "File path for CPU profile. If set, program is profiled.")
	flag.Parse()

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	infile, err := os.Open(infilename)
	if err != nil {
		fmt.Println("Error while opening input file:", err)
		os.Exit(2)
	}

	mg, err := graph.DeserializeMarkov(infile)
	if err != nil {
		panic(err)
	}

	sps := mg.StartingPoints()

	for i := 0; i < numwalks; i++ {
		word := sps[rand.Intn(len(sps))]
		walk := &strings.Builder{}

		for k := 0; k < length; k++ {
			walk.WriteString(word)

			if k < length-1 {
				walk.WriteRune(' ')
				word = pick(mg.Relations(word))

				// terminal word
				if word == "" {
					break
				}
			}
		}

		fmt.Println(walk.String())
	}
}

func pick(rels map[string]uint64) string {
	cdf := make([]uint64, 0, len(rels))
	words := make([]string, 0, len(rels))

	var acc uint64

	for k, v := range rels {
		acc += v
		cdf = append(cdf, acc)
		words = append(words, k)
	}

	sel := uint64(rand.Intn(int(acc)))

	for i, n := range cdf {
		if n > sel {
			return words[i]
		}
	}

	panic("WE SHOULDN'T BE HERE")
}
