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
	"io"

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

		mg.AddRelation(prev, word)

		prev = word
	}
}

var utol = [...]rune{
	'a', 'b', 'c', 'd', 'e',
	'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y',
	'z'}

// normalize normalizes the word by keeping only the characters (and hyphens) and lowercasing everything
func normalize(word string) string {
	norm := make([]rune, 0, len(word))

	// first pass: lowercase and filter
	for _, char := range word {
		if char >= 'A' && char <= 'Z' {
			norm = append(norm, utol[char-'A'])

		} else if (char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') ||
			char == '-' {
			norm = append(norm, char)
		}
	}

	return string(norm)
}
