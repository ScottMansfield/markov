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

package walk

import (
	"github.com/ScottMansfield/markov/graph"
)

type Walker struct {
	mg graph.Markov
}

func NewWalker(mg graph.Markov) *Walker {
	return &Walker{mg}
}

func (w *Walker) Walk(length int) []string {
	var cur string
	for k := range w.mg {
		cur = k
		break
	}

	ret := make([]string, 0, length)
	ret = append(ret, cur)

	for i := 1; i < length; i++ {
		cur = weightedRandomPick(w.mg.Relations(cur))
	}
}

func weightedRandomPick(choices map[string]uint64) string {
	//
}
