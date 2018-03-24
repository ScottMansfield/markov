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

package graph

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Markov is the main graph data structure
type Markov map[string]map[string]uint64

// NewMarkov creates a new Markov ready to be used
func NewMarkov() Markov {
	return make(Markov)
}

// AddRelation adds a relation between two words to the graph
// and keeps track of number of occurrences.
func (mg Markov) AddRelation(from, to string) {
	if conns, ok := mg[from]; ok {
		conns[to] = conns[to] + 1
	} else {
		mg[from] = map[string]uint64{
			to: 1,
		}
	}
}

// Serialize serializes the entire graph and sends the data to the given io.Writer
func (mg Markov) Serialize(w io.Writer) error {
	for from, tomap := range mg {
		for to, count := range tomap {
			if _, err := fmt.Fprintf(w, "%s %s %d\n", from, to, count); err != nil {
				return err
			}
		}
	}

	return nil
}

// Deserialize reads the serialized graph data from the given reader to recreate the graph
func Deserialize(r io.Reader) (Markov, error) {
	mg := make(Markov)

	s := bufio.NewScanner(r)
	for s.Scan() {
		raw := s.Text()
		segs := strings.Split(raw, " ")

		from := segs[0]
		to := segs[1]
		weightraw := segs[2]

		weight, err := strconv.ParseUint(weightraw, 10, 64)
		if err != nil {
			panic(err)
		}

		if _, ok := mg[from]; ok {
			mg[from][to] = weight
		} else {
			mg[from] = map[string]uint64{
				to: weight,
			}
		}
	}

	return mg, nil
}
