package compiler

// Author: Ersi Ni

import (
	"testing"
)

func TestProcessIdentifier(t *testing.T) {
	table := NewSymbolTable()
	for idx, testcase := range TrieTestSuite {
		output, err := table.Process(testcase.Identifier, testcase.Flag)
		if err != nil {
			t.Errorf("Error when processing identifiers, "+
				"testcase number %v, input: %v, error msg: %v",
				idx, testcase.Identifier, err)
		}
		if output != testcase.Expected {
			t.Errorf("Result mismatch for testcase %v, "+
				"input: %v, expected: %v, but got output: %v",
				idx, testcase.Identifier, testcase.Expected, output)
		}
	}
}

var TrieTestSuite = []struct {
	Identifier string
	Flag       FlagVar
	Expected   int
}{
	{"Protected", Dynamic, 0},
	{"Public", Dynamic, 1},
	{"Fake", Static, -1},
}

func TestSymbolTablePointer(t *testing.T) {
	st := NewSymbolTable()
	current := st.TrieHead
	n := NewTrieNode()
	current.set('z', n)
	if yes, err := st.TrieHead.has('z'); err != nil {
		t.Errorf("character %q is not valid", 'z')
	} else if !yes {
		t.Errorf("expected Yes but got %v", yes)
	}
}

var tests = getMemoryBenchmarkTests()

func BenchmarkMemoryFootprintSymbolTable(b *testing.B) {
	TheTrie := NewSymbolTable()
	b.Logf("testsuite of length %d loaded...\n", len(tests))
	for _, t := range tests {
		TheTrie.Process(t.Input, Dynamic)
	}
	b.Logf("built %d words in symbol table.\n", len(TheTrie.Table))
}
