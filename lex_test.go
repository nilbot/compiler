package compiler

import (
	"unicode/utf8"

	"testing"
)

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
	{
		input9,
		[]Token{
			{TokenIdentifier, 17 - 17, "8"},
			{TokenIdentifier, 22 - 17, "9"},
			{TokenError, 26 - 17, "no matching state for rune '~'"},
			{TokenIdentifier, 28 - 17, "10"},
			{TokenIdentifier, 35 - 17, "11"},
			{TokenIdentifier, 38 - 17, "12"},
			{TokenIdentifier, 40 - 17, "13"},
			{TokenError, 46 - 17, "no matching state for rune '?'"},
			{TokenText, 48 - 17, "Here And Th~~~ere"},
			{TokenIdentifier, 73 - 17, "14"},
			{TokenError, 77 - 17, "not matched \" found for string token, position 71"},
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
		e := 0
		for tokenID, token := range outputs {
			if token.P != testcase.Expected[tokenID].P ||
				token.T != testcase.Expected[tokenID].T ||
				token.V != testcase.Expected[tokenID].V {
				if tokenID+1 < len(testcase.Expected) {
					e = tokenID + 1
				}
				if e > tokenID {
					t.Errorf("error in set %d id %d! \n%v\n expected %v got %v",
						id,
						tokenID,
						testcase.Input[testcase.Expected[tokenID].P:testcase.Expected[e].P],
						testcase.Expected[tokenID], token)
				} else {
					t.Errorf("error in set %d id %d! \n%v\n expected %v got %v",
						id,
						tokenID,
						testcase.Input[testcase.Expected[tokenID].P:],
						testcase.Expected[tokenID], token)
				}
			}
		}
	}
}

func BenchmarkLexingShortSequences(b *testing.B) {
	for i, test := range OfficialTests {
		outputs := collect(test.Input)
		b.Logf("%d tokens for long-test %d of length %d", len(outputs), i, utf8.RuneCountInString(test.Input))
		ec := 0
		for _, t := range outputs {
			if t.T == TokenError && t.V != "" {
				ec++
			}
		}
		b.Logf("gathered %d error tokens", ec)
	}
}
func BenchmarkLexingLongSequences(b *testing.B) {
	for i, test := range LongSequences {
		outputs := collect(test)
		b.Logf("%d tokens for long-test %d of length %d", len(outputs), i, utf8.RuneCountInString(test))
		ec := 0
		for _, t := range outputs {
			if t.T == TokenError && t.V != "" {
				ec++
			}
		}
		b.Logf("gathered %d error tokens", ec)
	}
}

var LongSequences = generateFromDataset()
