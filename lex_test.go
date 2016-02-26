package compiler

import "testing"

var legal *Lexer

var LexerLegalWordsTestSuite = []struct {
	Identifier []string
	Expected   []int
}{
	{[]string{"Protected", "Private", "Public", "Uplink"},
		[]int{2, 0, 1, -1}},
	{[]string{"Try", "Exception", "Static", "Secondary", "Primary"},
		[]int{7, 6, 3, -1, 4}},
	{[]string{"Fake", "Integer"},
		[]int{-1, 5}},
	{[]string{"Test", "Main"},
		[]int{-1, -1}},
	{[]string{"Whatev", "Private", "Private", "Private"},
		[]int{-1, 0, 0, 0}},
}

var OfficialTests = []struct {
	Input    string
	Expected []Token
}{
	{
		input1,
		[]Token{
			{TokenIdentifier, 0, "8"},
			{TokenIdentifier, 5, "9"},
			{TokenIdentifier, 12, "10"},
		},
	},
	{
		input2,
		[]Token{
			{TokenIdentifier, 0, "8"},
			{TokenIdentifier, 5, "9"},
			{TokenIdentifier, 9, "10"},
		},
	},
	{
		input3,
		[]Token{
			{TokenIdentifier, 0, "8"},
			{TokenIdentifier, 3, "0"},
			{TokenIdentifier, 11, "9"},
			{TokenIdentifier, 15, "10"},
			{TokenIdentifier, 19, "11"},
			{TokenIdentifier, 25, "8"},
			{TokenIdentifier, 28, "1"},
			{TokenIdentifier, 35, "9"},
			{TokenIdentifier, 39, "12"},
		},
	},
}

func setupBarebone(t *testing.T) {
	legal = lexLegalWordsTester(t)
	legal.buildLegalWords()
}

func TestParseLegalKeywords(t *testing.T) {
	for idx, testcase := range LexerLegalWordsTestSuite {
		setupBarebone(t)
		outputs := legal.legalRun(t, testcase.Identifier)
		if !testEqualIntArray(testcase.Expected, outputs) {
			t.Errorf("Result mismatch for testcase %v, "+
				"input: %v, expected: %v, but got output: %v",
				idx, testcase.Identifier, testcase.Expected, outputs)
		}
	}
}

func TestLexer(t *testing.T) {
	for id, testcase := range OfficialTests {
		outputs := collect(testcase.Input)
		for tokenID, token := range outputs {
			if token.P != testcase.Expected[tokenID].P ||
				token.T != testcase.Expected[tokenID].T ||
				token.V != testcase.Expected[tokenID].V {
				t.Errorf("discrepency in test %d id %d! expected:%v got %v", id, tokenID, testcase.Expected[tokenID], token)
			}
		}
	}
}
