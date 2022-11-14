/*
 *******************************************************************************
 *                        (C) Copyright 2022 Fantasy Inc                       *
 * Created: 14/11/2022                                                         *
 *                                                                             *
 * Programmer(s):                                                              *
 * - Micrified                                                                 *
 *                                                                             *
 * Description:                                                                *
 *  Simple markov chain generator                                              *
 *                                                                             *
 *******************************************************************************
*/

package markov

import (
	"io"
	"bufio"
	"errors"
	"strings"
	"math/rand"
)

/*
 *******************************************************************************
 *                              Type definitions                               *
 *******************************************************************************
*/

// Markov state: Multi-word prefix, and suffix slice
type State struct {
	Prefix   []string;
	Suffixes []string;
}

// Markov generator: Prefix-to-state hashmap and prefix length
type Generator struct {
	Table   map[string]*State;
	Prefix_len int;
}

/*
 *******************************************************************************
 *                                   Methods                                   *
 *******************************************************************************
*/

// [Internal] Returns a key, given a list of strings
func (g *Generator) keyFrom (prefixes []string) string {
	return strings.Join(prefixes, "")
}

// [External] Builds prefix-state hashmap from word-based input stream
func (g *Generator) Build (in *io.Reader, 
	                       f func([]byte, bool)(int, []byte, error)) error {
	var i int;
	scanner := bufio.NewScanner(bufio.NewReader(*in))
	scanner.Split(f);
	
	// Construct the prefix: requires Prefix_len words
	prefix := make([]string, g.Prefix_len)
	for i = 0; scanner.Scan() && i < g.Prefix_len; i++ {
		prefix[i] = scanner.Text();
	}
	if i < g.Prefix_len {
		return errors.New("prefix may not exceed input size");
	}

	// Install suffix; continue with all other prefixes 
	for scanner.Scan() {
		var state *State;
		key, suffix := g.keyFrom(prefix), scanner.Text();
		if _, ok := g.Table[key]; ok {
			state = g.Table[key];
		} else {
			state = &State{Prefix: prefix, Suffixes: []string{}}
			g.Table[key] = state;
		}
		state.Suffixes = append(state.Suffixes, suffix);
		prefix = append(prefix[1:], suffix);
	}

	return nil;
}

// [External] Generates and returns a markov chain string; else non-nil error
func (g *Generator) DelimitedString (limit int, delim string) (string, error) {

	// Build the list of table keys
	keys := []string{}
	for key, _ := range g.Table {
		keys = append(keys, key)
	}

	// Select a random starting point
	start := keys[rand.Intn(len(keys))]

	// Create the string
	prefix, elements := g.Table[start].Prefix, []string{}
	for i := len(prefix); i < limit; i++ {
		if state, ok := g.Table[g.keyFrom(prefix)]; ok {
			suffix := state.Suffixes[rand.Intn(len(state.Suffixes))]
			elements = append(elements, suffix)
			prefix = append(prefix[1:], suffix)
		} else {
			break;
		}
	}
	elements = append(g.Table[start].Prefix, elements...)
	return strings.Join(elements, delim), nil
}

// [External] Wraps the markov chain generator with a space delimiter
func (g *Generator) String (limit int) (string, error) {
	return g.DelimitedString(limit, " ")
}