package compiler

import (
	"strconv"
	"strings"
)

// Author: Ersi Ni

// SymbolID typedef int for constants
type SymbolID int

// Grammar def T and NT, here the tiny spec allows us to use couple numbers to
// represent them all.
const (
	TBegin     SymbolID = iota // Terminal begin
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
	TEnd                       // terminal end
	NTBegin                    // nonterminal begin
	S                          //sentence start
	Epsilon                    // empty production
	BExp                       // expression
	BExp2                      // Bexp' to eliminate left recursion
	BTerm                      // (math) term
	BFactor                    // factor
	BConst                     // constant
	NTEnd                      // nonterminal end
	End        = 99            // sentence end
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

// Terminal tells if the the symbol is terminal
func (s SymbolID) Terminal() bool {
	return s < TEnd && s > TBegin
}

//NonTerminal tells if the symbol is nonterminal
func (s SymbolID) NonTerminal() bool {
	return s < NTEnd && s > NTBegin
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
	rst[Epsilon] = []Production{}
	return rst
}

// pathinfo contains info for stats of matched/unmatched attempt and when
// successfully matched, the tree traversal path for the symbol, depth etc.
type pathinfo struct {
}

// Parser for one sentence
type Parser struct {
	symbols   []Symbol
	chatty    bool
	discarded int
	tried     int
	log       Logger
	// diagnostics chan pathinfo
}

func loadSymbols(sentence string) []Symbol {
	var rst []Symbol
	parts := strings.Fields(sentence)
	for i := 0; i < len(parts); i += 2 {
		id, _ := strconv.ParseInt(parts[i], 10, 0)
		att, _ := strconv.ParseInt(parts[i+1], 10, 0)
		symbol := Symbol{
			ID:        SymbolID(id),
			Attribute: int(att),
		}
		rst = append(rst, symbol)
	}
	return rst
}

func (p *Parser) markFinish(success bool) {
	if success {
		p.log.Logf("\n==== parsing successful ====\n"+
			"original input:\n%v\n tried %v; "+
			"discarded %v; successfully matched %v.\n"+
			"~#~#~#~#~#~#~#~#~#~#~#~#~#\n",
			p.symbols,
			p.tried,
			p.discarded,
			p.tried-p.discarded)
	} else {
		p.log.Logf("\n==== parsing ultimately failed ====\n"+
			"original input:\n%v\n tried %v; "+
			"discarded %v; successfully matched %v.\n"+
			"~@~@~@~@~@~@~@~@~@~@~@~@~@\n",
			p.symbols,
			p.tried,
			p.discarded,
			p.tried-p.discarded)
	}
}

func (p *Parser) markMatch(leftHandSide SymbolID, prodIdx, symIdx int,
	rightHandSideCurrentSymbol SymbolID) {
	p.log.Logf("[partial success]\n symbol %v inside LHS %v\n "+
		"production index of prods derived from LHS: %v\n"+
		"symbol index of current prod: %v\n",
		rightHandSideCurrentSymbol,
		leftHandSide,
		prodIdx,
		symIdx)
}

func (p *Parser) missMatch(leftHandSide SymbolID, prodIdx, symIdx int,
	rightHandSideCurrentSymbol SymbolID) {
	p.log.Logf("[discard]\n miss matched symbol %v inside LHS %v\n"+
		"production index of prods derived from LHS: %v\n"+
		"symbol index of current prod: %v\n",
		rightHandSideCurrentSymbol,
		leftHandSide,
		prodIdx,
		symIdx)
}

func (p *Parser) dfs(lhs SymbolID, startPos int) (match bool, pos int) {
	pos = startPos

	if pos == -1 {
		return false, -1
	}
	if pos == len(p.symbols)-1 && p.symbols[len(p.symbols)-1].ID == End {
		return true, pos
	}
	if lhs == Epsilon {
		return true, pos
	}
	for pIdx, prod := range LREProduction[lhs] {
		p.tried++

		for sIdx, symbol := range prod.RHS {
			if symbol.Terminal() {
				if p.symbols[pos].ID == symbol {
					p.markMatch(lhs, pIdx, sIdx, symbol)
					pos++

					continue
				} else {
					p.missMatch(lhs, pIdx, sIdx, symbol)
					p.discarded++

					goto OUTERLOOP_CONTINUE
				}
			} else if symbol.NonTerminal() {
				ok, position := p.dfs(symbol, pos)
				if ok {
					pos = position
					continue
				} else {
					p.discarded++
					goto OUTERLOOP_CONTINUE
				}
			} else if symbol == End {
				// p.markFinish()
				return true, pos
			}
			p.log.Errorf("\n!!!\nthere might be error in "+
				"your productions table: "+
				"current pos %v, current lhs %v, "+
				"current rhs %v, current rhs sym %v\n!!!\n\n",
				pos, lhs, prod, symbol)
		}
		return true, pos
	OUTERLOOP_CONTINUE:
		pos = startPos // try new prod in rhs rules

	}
	p.log.Logf("[failed]\n grammar not matched, symbols \n%v\n"+
		"failed at pos %v\n",
		p.symbols, pos)
	return false, -1
}

// RunDFS runs DFS (or fancy named recursive descent traversal)
func (p *Parser) RunDFS() {
	if len(p.symbols) == 0 {
		p.log.Errorf("symbols input is empty array, this is " +
			"considered as error\nCheck your input.\n")
	}
	if p.chatty {
		p.log.Logf("\n****\nsymbols: %v\n****\n", p.symbols)
	}
	success, _ := p.dfs(S, 0)
	p.markFinish(success)
	// close(p.diagnostics)
}

// NewDFSParser construct a parser pointer which runs DFS to parse the input
func NewDFSParser(syms []Symbol, mylogger Logger, verbosity bool) *Parser {
	point := &Parser{
		symbols: syms,
		log:     mylogger,
		chatty:  verbosity,
		// diagnostics: make(chan pathinfo),
	}
	// go point.RunDFS()
	return point
}

// Logger custom logging mechnaism
type Logger interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
