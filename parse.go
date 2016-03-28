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
	BFactorP                   // LL(1) grammar rule
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
	}
	rst[BConst] = []Production{
		{[]SymbolID{trueConst}},
		{[]SymbolID{falseConst}},
	}
	return rst
}

// Parser for one sentence
type Parser struct {
	symbols   []Symbol
	chatty    bool
	discarded int
	tried     int
	log       Logger
	grammar   Productions
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
		p.log.Logf("\n==== Grammatical ====\n")
	} else {
		p.log.Logf("\n=== Ungrammatical ===\n")
	}
	p.log.Logf("original input:\n%v\n"+
		"# of symbols: %v;\n"+
		"tried %v; discarded %v; successfully matched %v.\n",
		p.symbols,
		len(p.symbols)-1,
		p.tried,
		p.discarded,
		p.tried-p.discarded)
}

func (p *Parser) markMatch(leftHandSide SymbolID, prodIdx, symIdx int,
	rightHandSideCurrentSymbol SymbolID) {
	if p.chatty {
		p.log.Logf("[partial success]"+
			"\n symbol %v inside LHS %v\n "+
			"production index of prods derived from LHS: %v\n"+
			"symbol index of current prod: %v\n",
			rightHandSideCurrentSymbol,
			leftHandSide,
			prodIdx,
			symIdx)
	}
}

func (p *Parser) missMatch(leftHandSide SymbolID, prodIdx, symIdx int,
	rightHandSideCurrentSymbol SymbolID) {
	if p.chatty {
		p.log.Logf("[discard]"+
			"\n miss matched symbol %v inside LHS %v\n"+
			"production index of prods derived from LHS: %v\n"+
			"symbol index of current prod: %v\n",
			rightHandSideCurrentSymbol,
			leftHandSide,
			prodIdx,
			symIdx)
	}
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
	for pIdx, prod := range p.grammar[lhs] {
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
				return true, pos
			}
			p.log.Errorf("\n!!!\nthere might be error in "+
				"your productions table: %v,"+
				"current pos %v, current lhs %v, "+
				"current rhs %v, current rhs sym %v\n!!!\n\n",
				p.grammar, pos, lhs, prod, symbol)
		}
		return true, pos
	OUTERLOOP_CONTINUE:
		pos = startPos // try new prod in rhs rules
	}
	if p.chatty {
		p.log.Logf("[failed]\n grammar not matched, symbols \n%v\n"+
			"failed at pos %v\n",
			p.symbols, pos)
	}
	return false, -1
}

// RunDFS runs DFS (or fancy named recursive descent traversal)
func (p *Parser) RunDFS() bool {
	if len(p.symbols) == 0 {
		p.log.Errorf("symbols input is empty array, this is " +
			"considered as error\nCheck your input.\n")
	}

	success, _ := p.dfs(S, 0)
	p.markFinish(success)
	return success
}

// NewDFSParser construct a parser pointer which runs DFS to parse the input
func NewDFSParser(syms []Symbol, mylogger Logger,
	gid GrammarID, verbosity bool) *Parser {
	var mygrammar map[SymbolID][]Production
	switch gid {
	case LRE:
		mygrammar = buildLREProduction()
	case LL1:
		mygrammar = buildLL1Productions()
	default:
		panic("this should never happen.")
	}
	point := &Parser{
		symbols: syms,
		log:     mylogger,
		chatty:  verbosity,
		grammar: mygrammar,
	}
	return point
}

// Logger custom logging mechnaism
type Logger interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// GrammarID notes what grammar version we are dealing with
type GrammarID int

// all grammar version listing
const (
	Original GrammarID = iota // original unmodified
	LRE                       // Left Recusion Eliminated
	LL1                       // LL(1) Grammar
)

// LL1Productions left recursion eliminated production table
func buildLL1Productions() map[SymbolID][]Production {
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
		{[]SymbolID{id, BFactorP}},
		{[]SymbolID{BConst}},
	}
	rst[BFactorP] = []Production{
		{[]SymbolID{Epsilon}},
		{[]SymbolID{eq, num}},
		{[]SymbolID{gt, num}},
		{[]SymbolID{lt, num}},
	}
	rst[BConst] = []Production{
		{[]SymbolID{trueConst}},
		{[]SymbolID{falseConst}},
	}
	return rst
}

// Productions synonym for ease of use
type Productions map[SymbolID][]Production

func NullableNT(NT SymbolID, G Productions, avoid []SymbolID) bool {
	for _, p := range G[NT] {
		if noTerminals(p.RHS) && noNT(NT, p.RHS) &&
			noMustAvoid(p.RHS, avoid) &&
			eachNullable(p.RHS, G, plus(avoid, NT)) {
			return true
		}
	}
	return false
}

func noTerminals(rhs []SymbolID) bool {
	for _, s := range rhs {
		if s.Terminal() {
			return false
		}
	}
	return true
}

func noNT(nt SymbolID, rhs []SymbolID) bool {
	for _, s := range rhs {
		if s == nt {
			return false
		}
	}
	return true
}

func noMustAvoid(rhs, avoid []SymbolID) bool {
	set := make(map[SymbolID]bool)
	for _, s := range avoid {
		set[s] = true
	}
	for _, s := range rhs {
		if set[s] {
			return false
		}
	}
	return true
}

func plus(list []SymbolID, target SymbolID) []SymbolID {
	for _, s := range list {
		if s == target {
			return list
		}
	}
	return append(list, target)
}

func minus(list []SymbolID, target SymbolID) []SymbolID {
	var rst []SymbolID
	for _, s := range list {
		if s != target {
			rst = append(rst, s)
		}
	}
	return rst
}

func eachNullable(rhs []SymbolID, gr Productions, avoid []SymbolID) bool {
	for _, s := range rhs {
		if !NullableNT(s, gr, avoid) {
			return false
		}
	}
	return true
}

func FirstSet(seq []SymbolID, G Productions) []SymbolID {
	if len(seq) == 0 {
		return []SymbolID{Epsilon}
	} else if seq[0].Terminal() {
		return []SymbolID{seq[0]}
	} else {
		nt := seq[0]
		var f2 []SymbolID
		for _, p := range G[nt] {
			f2 = union(f2, FirstSet(p.RHS, G))
		}
		if notContains(f2, Epsilon) {
			return f2
		}
		return union(minus(f2, Epsilon),
			FirstSet(seq[1:len(seq)], G))
	}
}

func union(s1 []SymbolID, s2 []SymbolID) []SymbolID {
	maps := make(map[SymbolID]bool)
	for _, s := range s1 {
		maps[s] = true
	}
	for _, s := range s2 {
		maps[s] = true
	}
	var rst []SymbolID
	for k := range maps {
		rst = append(rst, k)
	}
	return rst
}

func notContains(set []SymbolID, target SymbolID) bool {
	for _, s := range set {
		if s == target {
			return false
		}
	}
	return true
}
