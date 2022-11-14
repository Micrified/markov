# Markov

A markov-chain algorithm takes as input a series of overlapping phrases, and divides each phrase into two parts:

1. A prefix, which may consist of multiple words
2. A suffix, which is typically a single word

The chain is formed by emitting output phrases by randomly choosing the suffix for a given prefix, and is driven by the statistics of the original text. This module provides methods for processing a body of text into a structure upon which markov-chains can be derived.

## Algorithm and data-structure

1. Let `w1`, `w2` be the first two words of the text
2. Print `w1`, `w2`
3. loop:
4.    let `w3` be a randomly chosen word from the set of suffixes(`w1`,`w2`)
5.    print `w3`
6.    set `w1`, `w2` := `w2`, `w3`

The algorithm in the module is adapted from that by Brian W. Kernighan and Rob Pike in The Practice of Programming, Chapter 3, section 1. It is written in Go purely for pedagogical purposes. 

The internal representation of the markov-chain generator is a prefix-to-state hashmap, in which prefixes (defined as an N-length string) are mapped to state structures in memory. Each prefix has a list of associated suffixes, which are drawn upon during chain generation. The generator performs the processing in a single pass, and strings are generated in another pass on demand. 
