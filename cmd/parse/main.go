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
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/ScottMansfield/markov/graph"
	"github.com/ScottMansfield/markov/parse"
)

func main() {
	var infilename, outfilename, cpuprofile string

	flag.StringVar(&infilename, "if", "", "Folder containing Google Ngram files (gzipped input corpus)")
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

	var infilenames []string

	filepath.Walk(infilename, func(path string, f os.FileInfo, _ error) error {
		if f.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".gz" {
			infilenames = append(infilenames, path)
		}

		return nil
	})

	if len(infilenames) == 0 {
		log.Fatal("Could not find any ngram files")
	}

	outfile, err := os.Create(outfilename)
	if err != nil {
		log.Fatal("Error while opening output file:", err)
	}

	mg := graph.NewMarkov()

	for _, infilename := range infilenames {
		infile, err := os.Open(infilename)
		if err != nil {
			log.Fatal("Error while opening input file:", err)
		}

		parse.GoogleNgrams(infile, mg)

		infile.Close()
	}

	w := bufio.NewWriter(outfile)
	mg.Serialize(w)
	w.Flush()
	outfile.Close()
}

// TODO:
// - stemming(-ing, -es, -s, etc.),
// - capital first letters
// - which ones start sentences
// - normalize UTF-8 text
//   - https://godoc.org/golang.org/x/text/unicode/norm#Iter

// Can use google quint-grams to train
