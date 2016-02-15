package compiler

import (
	"testing"
)

func TestProcessIdentifier(t *testing.T) {
	for idx, testcase := range ProcessIdentifierTestSet {
		output, err := ProcessIdentifier(testcase.Input.Identifier,
			testcase.Input.Flag)
		if err != nil {
			t.Errorf("Error when processing identifiers, "+
				"testcase number %v, input: %v, error msg: %v",
				idx, testcase.Input, err)
		}
		if !equals(output, testcase.Expected) {
			t.Errorf("Result mismatch for testcase %v, "+
				"input: %v, expected: %v, but got output: %v",
				idx, testcase.Input, testcase.Expected, output)
		}
	}
}

var ProcessIdentifierTestSet = []TrieTestCase{
	{},
	{},
}
