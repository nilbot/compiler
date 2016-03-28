package compiler

import (
	"strconv"
	"strings"
	"testing"
)

func parseTestCases(s string) [][]Symbol {
	var rst [][]Symbol
	var chunk []Symbol
	nums := strings.Fields(s)
	for i := 0; i < len(nums); i += 2 {
		id, err := strconv.ParseInt(nums[i], 10, 0)
		if err != nil {
			panic("check your input text")
		}
		att, err := strconv.ParseInt(nums[i+1], 10, 0)
		if err != nil {
			panic("check your input text")
		}
		sym := Symbol{SymbolID(id), int(att)}
		chunk = append(chunk, sym)
		if sym.ID == End {
			rst = append(rst, chunk)
			chunk = make([]Symbol, 0)
		}
	}
	return rst
}

func TestDFSParser(t *testing.T) {
	var suc, fail int
	t.Logf("\n\n\n\nLRE\n\n\n\n\n")
	suc, fail = 0, 0
	testCases1 := parseTestCases(topDownParsingInput1)
	for _, testcase := range testCases1 {
		p := NewDFSParser(testcase, t, LRE, false)
		s := p.RunDFS()
		if s {
			suc++
		} else {
			fail++
		}
	}
	t.Logf("\nLRE input1 # of grammatical: %v\n", suc)
	t.Logf("\nLRE input1 # of ungrammatical: %v\n\n\n", fail)
	suc, fail = 0, 0
	testCases2 := parseTestCases(topDownParsingInput2)
	for _, testcase := range testCases2 {
		p := NewDFSParser(testcase, t, LRE, false)
		s := p.RunDFS()
		if s {
			suc++
		} else {
			fail++
		}
	}
	t.Logf("\nLRE input2 # of grammatical: %v\n", suc)
	t.Logf("\nLRE input2 # of ungrammatical: %v\n\n\n", fail)
	suc, fail = 0, 0
	t.Logf("\n\n\n\nLL1\n\n\n\n\n")
	for _, testcase := range testCases1 {
		p := NewDFSParser(testcase, t, LL1, false)
		s := p.RunDFS()
		if s {
			suc++
		} else {
			fail++
		}
	}
	t.Logf("\nLL1 input1 # of grammatical: %v\n", suc)
	t.Logf("\nLL1 input1 # of ungrammatical: %v\n\n\n", fail)
	suc, fail = 0, 0
	for _, testcase := range testCases2 {
		p := NewDFSParser(testcase, t, LL1, false)
		s := p.RunDFS()
		if s {
			suc++
		} else {
			fail++
		}
	}
	t.Logf("\nLL1 input2 # of grammatical: %v\n", suc)
	t.Logf("\nLL1 input2 # of ungrammatical: %v\n\n\n", fail)
}

// tom $$$
//
// not true $$$
//
// dick and harry or tom or true $$$
//
// tom or ( not ( dick ) or ( harry ) ) $$$
//
// 7 > 3 and not 8 = 11 $$$
//
// 7 >
// 3 and not
// 8 = 11
// $$$
var topDownParsingInput1 = `2 1 99 0

7 0 8 0 99 0

2 9 5 0 2 33 6 0 2 1 6 0 8 0 99 0

2 1 6 0 3 0 7 0 3 0 2 9 4 0 6 0 3 0 2 33 4 0 4 0 99 0

1 7 11 0 1 3 5 0 7 0 1 8 10 0 1 11 99 0

1 7 11 0
1 3 5 0 7 0
1 8 10 0 1 11
99 0`

// tom = 33 $$$
// 2 2 10 0 1 33 99 0
// not bill or not harry $$$
// 7 0 2 3 6 0 7 0 2 4 99 0
// ( bill and tom and harry and tom and mary = 66 ) $$$
// 3 0 2 3 5 0 2 2 5 0 2 4 5 0 2 2 5 0 2 5 10 0 1 66 4 0 99 0
// ( not not tom < harry ) $$$
// 3 0 7 0 7 0 2 2 11 0 2 4 4 0 99 0
// tom < harry = dick and mary $$$
// 2 2 11 0 2 4 10 0 2 6 5 0 2 5 99 0
// ( tom and ( not harry ) $$$
// 3 0 2 2 5 0 3 0 7 0 2 4 4 0 99 0
// true ) $$$
// 8 0 4 0 99 0
// 123 > > harry $$$
// 1 123 12 0 12 0 2 4 99 0
// false friend $$$
// 9 0 2 7 99 0
// false = not true $$$
// 9 0 10 0 7 0 8 0 99 0
var topDownParsingInput2 = `2 2 10 0 1 33 99 0 

7 0 2 3 6 0 7 0 2 4 99 0 

3 0 2 3 5 0 2 2 5 0 2 4 5 0 2 2 5 0 2 5 10 0 1 66 4 0 99 0 

3 0 7 0 7 0 2 2 11 0 2 4 4 0 99 0 

2 2 11 0 2 4 10 0 2 6 5 0 2 5 99 0 

3 0 2 2 5 0 3 0 7 0 2 4 4 0 99 0 

8 0 4 0 99 0 

1 123 12 0 12 0 2 4 99 0 

9 0 2 7 99 0 

9 0 10 0 7 0 8 0 99 0 `
