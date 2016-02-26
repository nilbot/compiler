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
	{
		input4,
		[]Token{
			{TokenIdentifier, 0, "8"},
			{TokenIdentifier, 5, "9"},
			{TokenIdentifier, 14, "-1"},
			{TokenIdentifier, 24, "10"},
			{TokenLeftParenthesis, 44 - 17, "0"},
			{TokenInteger, 45 - 17, "12"},
			{TokenIdentifier, 48 - 17, "11"},
			{TokenIdentifier, 51 - 17, "12"},
			{TokenRightParenthesis, 55 - 17, "0"},
			{TokenIdentifier, 57 - 17, "13"},
			{TokenIdentifier, 60 - 17, "14"},
		},
	},
	{
		input5,
		[]Token{
			{TokenIdentifier, 0, "8"},
			{TokenIdentifier, 23 - 17, "9"},
			{TokenText, 27 - 17, "strings of things"},
			{TokenIdentifier, 47 - 17, "10"},
			{TokenIdentifier, 51 - 17, "8"},
			{TokenIdentifier, 57 - 17, "9"},
			{TokenText, 61 - 17, "chickens\nwith\nwings"},
		},
	},
	{
		input6,
		[]Token{
			{TokenIdentifier, 17 - 17, "8"},
			{TokenIdentifier, 22 - 17, "9"},
			{TokenIdentifier, 26 - 17, "8"},
			{TokenIdentifier, 31 - 17, "9"},
			{TokenIdentifier, 35 - 17, "8"},
			{TokenIdentifier, 40 - 17, "10"},
		},
	},
	{
		input7,
		[]Token{
			{TokenInteger, 0, "888"},
			{TokenIdentifier, 21 - 17, "8"},
			{TokenIdentifier, 24 - 17, "9"},
			{TokenIdentifier, 30 - 17, "10"},
			{TokenIdentifier, 33 - 17, "11"},
			{TokenInteger, 37 - 17, "-1"},
			{TokenIdentifier, 43 - 17, "8"},
			{TokenText, 46 - 17, "not ok"},
		},
	},
	{
		input8,
		[]Token{
			{TokenIdentifier, 17 - 17, "1"},
			{TokenIdentifier, 23 - 17, "-1"},
			{TokenIdentifier, 26 - 17, "0"},
			{TokenIdentifier, 33 - 17, "-1"},
			{TokenIdentifier, 39 - 17, "-1"},
			{TokenIdentifier, 42 - 17, "-1"},
			{TokenIdentifier, 44 - 17, "-1"},
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
