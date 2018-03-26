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
	"sync"
)

type markovNode struct {
	lock *sync.Mutex
	data map[string]uint64
}

type markovData map[string]markovNode

// Markov is the main graph data structure
type Markov struct {
	lock *sync.Mutex
	data markovData
}

// NewMarkov creates a new Markov ready to be used
func NewMarkov() *Markov {
	return &Markov{
		lock: &sync.Mutex{},
		data: make(markovData),
	}
}

// IncRelation adds 1 to the count for a relation between two words.
func (mg Markov) IncRelation(from, to string) {
	mg.IncRelationBy(from, to, 1)
}

// IncRelationBy adds n to the count for a relation between two words.
func (mg Markov) IncRelationBy(from, to string, n uint64) {
	mg.lock.Lock()
	node, ok := mg.data[from]
	mg.lock.Unlock()

	if ok {
		node.lock.Lock()
		node.data[to] = node.data[to] + n
		node.lock.Unlock()
	} else {
		newnode := markovNode{
			lock: &sync.Mutex{},
			data: map[string]uint64{
				to: n,
			},
		}
		mg.lock.Lock()
		node, ok = mg.data[from]
		if ok {
			mg.lock.Unlock()
			node.lock.Lock()
			node.data[to] = node.data[to] + n
			node.lock.Unlock()
		} else {
			mg.data[from] = newnode
			mg.lock.Unlock()
		}
	}
}

// Serialize serializes the entire graph and sends the data to the given io.Writer
func (mg Markov) Serialize(w io.Writer) error {
	mg.lock.Lock()
	defer mg.lock.Unlock()
	for from, to := range mg.data {
		for to, count := range to.data {
			if _, err := fmt.Fprintf(w, "%s %s %d\n", from, to, count); err != nil {
				return err
			}
		}
	}

	return nil
}

// Deserialize reads the serialized graph data from the given reader to recreate the graph
func Deserialize(r io.Reader) (*Markov, error) {
	mg := NewMarkov()

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

		mg.IncRelationBy(from, to, weight)
	}

	return mg, nil
}
