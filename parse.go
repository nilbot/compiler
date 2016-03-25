package compiler

// Author: Ersi Ni

// SymbolID typedef int for constants
type SymbolID int

// Grammar def T and NT, here the tiny spec allows us to use couple numbers to
// represent them all.
const (
	TBegin     SymbolID = iota // ignore 0
	num                        // 0-9
	id                         // id
	lb                         // '('
	rb                         // ')'
	and                        // and op
	or                         // or op
	not                        // not op
	trueConst                  // true
	falseConst                 // false
	eq                         // '='
	lt                         // '<'
	gt                         // '>'
	TEnd                       // terminal threshold

	NTBegin SymbolID = iota + 50 // ignore 50
	S                            //sentence start
	Epsilon                      // empty production
	BExp                         // expression
	BExp2                        // Bexp' to eliminate left recursion
	BTerm                        // (math) term
	BFactor                      // factor
	BConst                       // constant
	NTEnd                        // NT end

	End SymbolID = iota + 99 // sentence end
)

// Symbol represent a symbol of ID and attribute value
type Symbol struct {
	ID        SymbolID
	Attribute int
}

// Production is a L->R mapping
type Production struct {
	RHS []SymbolID // right hand side
}

// IsTerminal tells if the the symbol is terminal
func (s SymbolID) IsTerminal() bool {
	return s < TEnd
}

// LREProduction left recursion eliminated production table
var LREProduction = buildLREProduction()

func buildLREProduction() map[SymbolID][]Production {
	rst := make(map[SymbolID][]Production)
	rst[S] = []Production{{[]SymbolID{BExp, End}}}
	rst[BExp] = []Production{{[]SymbolID{BTerm, BExp2}}}
	rst[BExp2] = []Production{
		{[]SymbolID{and, BTerm, BExp2}},
		{[]SymbolID{or, BTerm, BExp2}},
		{[]SymbolID{Epsilon}},
	}
	rst[BTerm] = []Production{
		{[]SymbolID{BFactor}},
		{[]SymbolID{not, BTerm}},
	}
	rst[BFactor] = []Production{
		{[]SymbolID{lb, BExp, rb}},
		{[]SymbolID{id}},
		{[]SymbolID{BConst}},
		{[]SymbolID{id, eq, num}},
		{[]SymbolID{id, gt, num}},
		{[]SymbolID{id, lt, num}},
		{[]SymbolID{num, eq, num}},
		{[]SymbolID{num, gt, num}},
		{[]SymbolID{num, lt, num}},
	}
	rst[BConst] = []Production{
		{[]SymbolID{trueConst}},
		{[]SymbolID{falseConst}},
	}
	return rst
}

// pathinfo contains info for stats of matched/unmatched attempt and when
// successfully matched, the tree traversal path for the symbol, depth etc.
type pathinfo struct {
}

// Parser for one sentence
type Parser struct {
	inputStream string
	symbols     []Symbol
	diagnostics chan pathinfo
}

func loadSymbols(sentence string) []Symbol {
	var rst []Symbol
	return rst
}

// RunDFS runs DFS (or fancy named recursive descent traversal)
func (p *Parser) RunDFS() {

}

// NewDFSParser construct a parser pointer which runs DFS to parse the input
func NewDFSParser(input string) *Parser {
	point := &Parser{
		inputStream: input,
		symbols:     loadSymbols(input),
		diagnostics: make(chan pathinfo),
	}
	go point.RunDFS()
	return point
}
