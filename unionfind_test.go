package compiler

import (
	"bufio"
	"strconv"
	"strings"
	"testing"
)

func simpleNewUF(count int) *unionfind {
	p := make([]int, count)
	sz := make([]int, count)
	for i := 0; i != count; i++ {
		p[i] = i
		sz[i] = 1
	}
	return &unionfind{
		p,
		sz,
		count,
	}
}

var tinyUF = `10
4 3
3 8
6 5
9 4
2 1
8 9
5 0
7 2
6 1
1 0
6 7
`

func TestUnionFind(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader(tinyUF))
	scanner.Scan()
	n, _ := strconv.ParseInt(scanner.Text(), 10, 0)
	uf := simpleNewUF(int(n))
	for scanner.Scan() {
		ab := strings.Split(scanner.Text(), " ")
		a, _ := strconv.ParseInt(ab[0], 10, 0)
		b, _ := strconv.ParseInt(ab[1], 10, 0)
		uf.Union(int(a), int(b))
	}
	if uf.Count() != 2 {
		t.Errorf("expected 2 components but got %v", uf.Count())
	}
}
