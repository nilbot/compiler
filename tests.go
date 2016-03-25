package compiler

// Author: Ersi Ni

import (
	"strconv"
	"testing"

	"io/ioutil"
)

func testEqualIntArray(a, b []int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func generateFromDataset() []string {
	var rst []string
	content, err := ioutil.ReadFile("words.test")
	if err != nil {
		for _, v := range memoryBenchmarkTest {
			s := ""
			for i := 0; i < 10; i++ {
				s += v.Input
			}
			rst = append(rst, s)
		}
		return rst
	}
	l := len(content)
	apiece := l / 10
	for s := 0; s < l; s += apiece {
		offset := s + apiece
		if offset < l {
			rst = append(rst, string(content[s:offset]))
		} else {
			rst = append(rst, string(content[s:]))
		}
	}
	return rst
}

func lexLegalWordsTester(t *testing.T) *Lexer {
	pointer := &Lexer{
		S:               NewSymbolTable(),
		Tokens:          make(chan Token),
		State:           nil,
		StartPosition:   0,
		Width:           0,
		LastPosition:    0,
		CurrentPosition: 0,
		In:              "",
	}
	if pointer == nil {
		t.Errorf("built a nil barebone lexer\n")
	}
	return pointer
}

func (l *Lexer) injectInput(in string) {
	l.In = in
}

func (l *Lexer) legalRun(t *testing.T, keywords []string) []int {
	concat := ""
	for _, k := range keywords {
		concat += k + " "
	}
	l.injectInput(concat)
	if l.In != concat {
		t.Errorf("inject failed, got %v\n", l.In)
	}
	go l.Run()
	var rst []int
	for {
		token := l.ConsumeToken()
		if token.T == TokenError {
			break
		}
		i, e := strconv.ParseInt(token.V, 10, 0)
		if e != nil {
			t.Errorf("not a int, %v\n", token)
		}
		rst = append(rst, int(i))

	}
	return rst
}

func collect(input string) (rst []Token) {
	l := NewLexer(input)
	for {
		token := l.ConsumeToken()
		if token.T == TokenError && token.V == "" {
			l.Flush()
			break
		}
		rst = append(rst, token)

	}
	return
}

const input1 = `this little piggy`
const input2 = `what has happened`
const input3 = `in Private say one thing in Public say another`
const input4 = `when speaking Privately to (12 or more) be discreet`
const input5 = `there are "strings of things" and there are "chickens
with
wings"
`
const input6 = `over and over and over again`
const input7 = `888 is quite ok but 88888 is "not ok"`
const input8 = `PublicAndPrivateShouldNotBeConfused`
const input9 = `when can ~ appear in a string? "Here An~d~ Th~~~~~~ere"
but "some string
`
