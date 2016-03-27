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
	testCases := parseTestCases(topDownParsingInput1)
	for _, testcase := range testCases {
		p := NewDFSParser(testcase, t, false)
		p.RunDFS()
	}
}

var topDownParsingInput1 = `2 1 99 0

7 0 8 0 99 0

2 9 5 0 2 33 6 0 2 1 6 0 8 0 99 0

2 1 6 0 3 0 7 0 3 0 2 9 4 0 6 0 3 0 2 33 4 0 4 0 99 0

1 7 11 0 1 3 5 0 7 0 1 8 10 0 1 11 99 0

1 7 11 0
1 3 5 0 7 0
1 8 10 0 1 11
99 0`
