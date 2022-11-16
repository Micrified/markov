/*
 *******************************************************************************
 * Created: 14/11/2022                                                         *
 *                                                                             *
 * Programmer(s):                                                              *
 * - Charles Randolph                                                          *
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
	Table      map[string]*State;
	Prefix_len int;
	Err        error;
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
	                       prefix_len int,
	                       f func([]byte, bool)(int, []byte, error)) error {
	var err error;
	var i int;
	var scanner *bufio.Scanner;
	var prefix []string;
	var table map[string]*State;

	// Check args
	if nil == in || nil == *in {
		err = errors.New("invalid argument: nil reader or reader pointer");
	} else if prefix_len < 1 {
		err = errors.New("invalid argument: prefix out of bounds [1,inf)");
	} else if nil == f {
		err = errors.New("invalid argument: invalid split function");
	} else {
		scanner = bufio.NewScanner(*in)
		scanner.Split(f);
	}
	
	// Construct the prefix: requires Prefix_len words
	if nil == err {
		prefix = make([]string, prefix_len);
		for i = 0; scanner.Scan() && i < prefix_len; i++ {
			prefix[i] = scanner.Text();
		}
		if i < prefix_len {
			err = errors.New("prefix may not exceed input size");
		}
	}

	// Install suffix; continue with all other prefixes 
	if nil == err {
		table = make(map[string]*State);
		for scanner.Scan() {
			var state *State;
			key, suffix := g.keyFrom(prefix), scanner.Text();
			if _, ok := table[key]; ok {
				state = table[key];
			} else {
				state = &State{Prefix: prefix, Suffixes: []string{}}
				table[key] = state;
			}
			state.Suffixes = append(state.Suffixes, suffix);
			prefix = append(prefix[1:], suffix);
		}
	}

	// Set table; Install/return potential scanner errors
	*g = Generator{Table: table, Prefix_len: prefix_len, Err: err};
	return err;
}

// [External] Generates and returns a markov chain string; else non-nil error
func (g *Generator) DelimitedString (limit int, delim string) (string, error) {

	// Check args
	if nil != g.Err {
		return "", g.Err;
	}

	// Build the list of table keys
	keys := []string{};
	for key, _ := range g.Table {
		keys = append(keys, key);
	}

	// Select a random starting point
	start := keys[rand.Intn(len(keys))];

	// Create the string
	prefix, elements := g.Table[start].Prefix, []string{};
	for i := len(prefix); i < limit; i++ {
		if state, ok := g.Table[g.keyFrom(prefix)]; ok {
			suffix := state.Suffixes[rand.Intn(len(state.Suffixes))];
			elements = append(elements, suffix);
			prefix = append(prefix[1:], suffix);
		} else {
			break;
		}
	}
	elements = append(g.Table[start].Prefix, elements...);
	return strings.Join(elements, delim), nil;
}

// [External] Wraps the markov chain generator with a space delimiter
func (g *Generator) String (limit int) (string, error) {
	return g.DelimitedString(limit, " ");
}