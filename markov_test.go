/*
 *******************************************************************************
 * Created: 14/11/2022                                                         *
 *                                                                             *
 * Programmer(s):                                                              *
 * - Charles Randolph                                                          *
 *                                                                             *
 * Description:                                                                *
 *  Simple markov chain generator test suite                                   *
 *                                                                             *
 *******************************************************************************
*/

package markov_test

import (
	"io"
	"bufio"
	"testing"
	"strings"
	"markov"
)

func TestMarkovTwoPrefixGenerator (t *testing.T) {
	input :=`Show your flowcharts and conceal your tables and I will be mystified.
	Show your tables and your flowcharts will be obvious.`;
	prefix, string_reader, limit := 2, strings.NewReader(input), 6;

	// Configure generator
	var g markov.Generator;
	reader := io.Reader(string_reader);
	err := g.Build(&reader, prefix, bufio.ScanWords);
	if nil != err {
		t.Errorf("Non-nil error on Build: %s", err.Error());
	}
	
	// Generate a bunch of markov with space delimiters
	for i := 0; i < 100; i++ {
		m, err := g.String(limit);
		if nil != err {
			t.Errorf("Non-nil error on DelimitedString: %s", err.Error());
		}
		ws := strings.Split(m, " ");
		if (len(ws) < (prefix+1) || len(ws) > limit) {
			t.Errorf("Output exceeds word boundaries: (%d,%d]", prefix, limit);
		}
	}

	// Generate a bunch of dash delimited markov chains
	for i := 0; i < 100; i++ {
		m, err := g.DelimitedString(limit, "-");
		if nil != err {
			t.Errorf("Non-nil error on DelimitedString: %s", err.Error());
		}
		ws := strings.Split(m, "-")
		if (len(ws) < (prefix+1) || len(ws) > limit) {
			t.Errorf("Output exceeds word boundaries: (%d,%d]", prefix, limit);
		}
	}
}

func TestMarkovIllegalPrefixLength (t *testing.T) {
	input :=`Too short`;
	prefix, string_reader := 3, strings.NewReader(input);

	// Configure generator
	var g markov.Generator;
	reader := io.Reader(string_reader);
	err := g.Build(&reader, prefix, bufio.ScanWords);
	if nil == err {
		t.Errorf("Build should fail if prefix is longer than input!");
	}
}