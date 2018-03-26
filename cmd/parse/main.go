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
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/ScottMansfield/markov/graph"
	"github.com/ScottMansfield/markov/parse"
)

func main() {
	var infilename, outfilename, cpuprofile string

	flag.StringVar(&infilename, "if", "", "Google Ngram file (gzipped input corpus)")
	flag.StringVar(&outfilename, "of", "", "Output file (serialized graph)")
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

	outfile, err := os.Create(outfilename)
	if err != nil {
		fmt.Println("Error while opening output file:", err)
		os.Exit(2)
	}

	mg := graph.NewMarkov()
	parse.GoogleNgrams(infile, mg)

	w := bufio.NewWriter(outfile)
	mg.Serialize(w)
	w.Flush()
}

// TODO:
// - stemming(-ing, -es, -s, etc.),
// - capital first letters
// - which ones start sentences
// - normalize UTF-8 text
//   - https://godoc.org/golang.org/x/text/unicode/norm#Iter

// Can use google quint-grams to train
