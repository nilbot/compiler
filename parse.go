package compiler

import (
	"fmt"
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
	num                        // 0-9 	: 1
	id                         // id 	: 2
	lb                         // '('	: 3
	rb                         // ')'	: 4
	and                        // and op	: 5
	or                         // or op	: 6
	not                        // not op	: 7
	trueConst                  // true	: 8
	falseConst                 // false	: 9
	eq                         // '='	: 10
	lt                         // '<'	: 11
	gt                         // '>'	: 12
	TEnd                       // terminal end
	NTBegin                    // nonterminal begin
	S                          // sentence start
	BExp                       // expression
	BExp2                      // Bexp' to eliminate left recursion
	BTerm                      // (math) term
	BFactor                    // factor
	BFactorP                   // LL(1) grammar rule
	BConst                     // constant
	NTEnd                      // nonterminal end
	Epsilon                    // empty production
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

// T tells if the the symbol is terminal
func (s SymbolID) T() bool {
	return s < TEnd && s > TBegin
}

//NT tells if the symbol is nonterminal
func (s SymbolID) NT() bool {
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
		len(p.symbols),
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
			if symbol.T() {
				if p.symbols[pos].ID == symbol {
					p.markMatch(lhs, pIdx, sIdx, symbol)
					pos++
					continue
				} else {
					p.missMatch(lhs, pIdx, sIdx, symbol)
					p.discarded++
					goto OUTERLOOP_CONTINUE
				}
			} else if symbol.NT() {
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
			} else if symbol == Epsilon {
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

// NullableNT is NT nullable? avoid list contains ε?
func NullableNT(NT SymbolID, G Productions, avoid []SymbolID) bool {
	for _, p := range G[NT] {
		if noTerminals(p.RHS) && notContains(p.RHS, NT) &&
			noMustAvoid(p.RHS, avoid) &&
			eachNullable(p.RHS, G, plus(avoid, NT)) {
			return true
		}
	}
	return false
}

func noTerminals(rhs []SymbolID) bool {
	for _, s := range rhs {
		if s.T() {
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

// FirstSet returns the First Set for seq in respect to G
func FirstSet(seq []SymbolID, G Productions) []SymbolID {
	if len(seq) == 0 {
		return []SymbolID{Epsilon}
	} else if len(seq) == 1 && seq[0] == Epsilon {
		return []SymbolID{Epsilon}
	} else if seq[0].T() || seq[0] == End {
		return []SymbolID{seq[0]}
	}
	nt := seq[0]
	var f2 []SymbolID
	for _, p := range G[nt] {
		f2 = union(f2, FirstSet(p.RHS, G))
	}
	if notContains(f2, Epsilon) {
		return f2
	}
	return union(minus(f2, Epsilon), FirstSet(seq[1:], G))

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

// FollowSet returns the Follow Set
func FollowSet(NT SymbolID) []SymbolID {
	return FS[NT]
}

// FollowSetHashTable SymbolID -> []SymbolID
type FollowSetHashTable map[SymbolID][]SymbolID

// FS FollowSetHashTable
var FS FollowSetHashTable

// IN InheritorsHashTable
var IN FollowSetHashTable

func buildFollowSetMap(G Productions) {
	FS = make(map[SymbolID][]SymbolID)
	IN = make(map[SymbolID][]SymbolID)
	FS[BExp] = plus([]SymbolID{}, End)
	for _, k := range NonTerminals() {
		for _, p := range G[k] {
			r := p.RHS
			l := len(r)
			for i := 0; i < l; i++ {
				s := r[i]
				if s.NT() {
					for _, f := range FirstSet(r[i+1:l],
						G) {
						if f == Epsilon {
							IN[k] = plus(IN[k], s)
						} else {
							FS[s] = plus(FS[s], f)
						}
					}
				}
			}
		}
	}
	idle := false
	for !idle {
		idle = true
		for nt1, value := range IN {
			for _, nt2 := range value {
				idle = fillNT(nt1, nt2)
			}
		}
	}
}

func fillNT(nt1, nt2 SymbolID) bool {
	idle := true
	for _, item := range FS[nt1] {
		if item.T() || item == End {
			if notContains(FS[nt2], item) {
				FS[nt2] = plus(FS[nt2], item)
				idle = false
			}
		}
	}
	return idle
}

// PredictiveParsingTable is a (r,c) pair to production map
type PredictiveParsingTable map[rxc]Production

// rxc is row cross column
type rxc struct {
	row SymbolID
	col SymbolID
}

// M is the table
var M PredictiveParsingTable

func buildM(G Productions) {
	M = make(map[rxc]Production)

	Xi := NonTerminals()
	for _, nt := range Xi {
		for _, a := range G[nt] {
			fi := FirstSet(a.RHS, G)
			for _, t := range fi {
				rc := rxc{nt, t}
				M[rc] = a
			}
			if contains(fi, Epsilon) {
				fo := FollowSet(nt)
				for _, t := range fo {
					if t.T() {
						rc := rxc{nt, t}
						M[rc] = a
					}
					if t == End {
						rc := rxc{nt, t}
						M[rc] = a
					}
				}
			}
		}
	}
}

func contains(set []SymbolID, target SymbolID) bool {
	for _, s := range set {
		if s == target {
			return true
		}
	}
	return false
}

// Stack Go has no built-in Stack
type Stack struct {
	container []SymbolID
}

// Empty tells if the stack is empty
func (s *Stack) Empty() bool { return len(s.container) == 0 }

// Push item on top
func (s *Stack) Push(t SymbolID) { s.container = append(s.container, t) }

// Top grabs top without modifying index
func (s *Stack) Top() SymbolID { return s.container[len(s.container)-1] }

// Pop removes top item
func (s *Stack) Pop() SymbolID {
	r := s.container[len(s.container)-1]
	s.container = s.container[:len(s.container)-1]
	return r
}

func (p *Parser) predictive() bool {
	st := newStack()
	st.Push(End)
	st.Push(S)
	idx := 0
	token := p.symbols[idx]

	for X := st.Top(); X != End; X = st.Top() {
		p.tried++
		if X.T() || X == End {
			if X == token.ID {
				oldX := X
				X = st.Pop()
				idx++
				token = p.symbols[idx]
				if p.chatty {
					p.log.Logf("X matched token %v, "+
						"new popped X is %v, "+
						"new token is %v, "+
						"stack size %v.",
						oldX,
						X,
						token.ID,
						len(st.container))
				}
			} else {
				p.log.Logf("[Error] Parsing failed on miss T")
				return false
			}
		} else if X == Epsilon {
			if p.chatty {
				p.log.Logf("X is %v, current stack %v, "+
					"original input symobls are: %v, "+
					"current symbol idx %v",
					X,
					st.container,
					p.symbols,
					idx)
			}
			st.Pop()
		} else {
			pr := M[rxc{X, token.ID}]
			if p.chatty {
				p.log.Logf("old stack: %v",
					st.container)
				p.log.Logf("expanding production: %v", pr)
				p.log.Logf("X is %v, "+
					"token is %v, "+
					"stack size %v.",
					X,
					token.ID,
					len(st.container))
			}
			if len(pr.RHS) > 0 {
				X = st.Pop()
				for i := len(pr.RHS) - 1; i >= 0; i-- {
					st.Push(pr.RHS[i])
				}
				if p.chatty {
					p.log.Logf("new stack: %v",
						st.container)
				}
			} else {
				p.log.Logf("[Error] Empty Cell")
				if p.chatty {
					p.log.Logf(
						"[Error]\nM(%v,%v)"+
							" cell is empty\n",
						X, token.ID)
				}
				return false
			}
		}
	}
	return true
}

func newStack() *Stack {
	return &Stack{make([]SymbolID, 0)}
}

// RunPredictiveParsing wraps predictive parsing routine and do some logging
func (p *Parser) RunPredictiveParsing() bool {
	success := p.predictive()
	if success {
		p.log.Logf("\n==== Grammatical ====\n")
	} else {
		p.log.Logf("\n=== Ungrammatical ===\n")
		p.log.Logf("failed with input: %v\n", p.symbols)
	}
	p.log.Logf("made %v matching with no backtracking", p.tried)
	return success
}

// NewPredictiveParser constructs a pointer of parser using LL1 grammar and
// builds followset and predictivetable
func NewPredictiveParser(input []Symbol, mylogger Logger,
	verbosity bool) *Parser {
	pointer := &Parser{
		grammar: buildLL1Productions(),
		chatty:  verbosity,
		symbols: input,
		log:     mylogger,
	}
	buildFollowSetMap(pointer.grammar)
	buildM(pointer.grammar)
	return pointer
}

// Terminals returns list of terminals
func Terminals() []SymbolID {
	var rst []SymbolID
	for i := TBegin + 1; i < TEnd; i++ {
		rst = append(rst, SymbolID(i))
	}
	return rst
}

// NonTerminals return list of non terminals
func NonTerminals() []SymbolID {
	var rst []SymbolID
	for i := NTBegin + 1; i < NTEnd; i++ {
		rst = append(rst, SymbolID(i))
	}
	return rst
}

func (p PredictiveParsingTable) String() string {
	rows := NonTerminals()
	cols := Terminals()
	rst := "\n\t" + fmt.Sprintln(cols)
	for _, r := range rows {
		rst = rst + r.String() + ": "
		for _, c := range cols {
			if len(p[rxc{r, c}].RHS) != 0 {
				rst = rst + fmt.Sprint(p[rxc{r, c}].RHS, " ")
			} else {
				rst = rst + fmt.Sprint("[ X ] ")
			}
		}

		rst = rst + "\n\n"
	}
	return rst
}

//String er for SymbolID
func (s SymbolID) String() string {
	switch s {
	case S:
		return "S"
	case BExp:
		return "BExp"
	case BExp2:
		return "BExp2"
	case BTerm:
		return "BTerm"
	case BFactor:
		return "BFactor"
	case BFactorP:
		return "BFactorP"
	case BConst:
		return "BConst"
	case Epsilon:
		return "ε"
	case id:
		return "id"
	case lb:
		return "("
	case rb:
		return ")"
	case num:
		return "num"
	case eq:
		return "="
	case lt:
		return "<"
	case gt:
		return ">"
	case and:
		return "and"
	case or:
		return "or"
	case not:
		return "not"
	case trueConst:
		return "true"
	case falseConst:
		return "false"
	case End:
		return "$$$"
	default:
		return fmt.Sprintf("UNDEFINED: %2d", s)
	}
}

//String er for SymbolID
func (s Symbol) String() string {
	var prefix, postfix string
	postfix = " " + strconv.Itoa(s.Attribute)
	switch s.ID {
	case S:
		prefix = "S"
	case BExp:
		prefix = "BExp"
	case BExp2:
		prefix = "BExp2"
	case BTerm:
		prefix = "BTerm"
	case BFactor:
		prefix = "BFactor"
	case BFactorP:
		prefix = "BFactorP"
	case BConst:
		prefix = "BConst"
	case Epsilon:
		prefix = "ε"
	case id:
		prefix = "id"
	case lb:
		prefix = "("
	case rb:
		prefix = ")"
	case num:
		prefix = "num"
	case eq:
		prefix = "="
	case lt:
		prefix = "<"
	case gt:
		prefix = ">"
	case and:
		prefix = "and"
	case or:
		prefix = "or"
	case not:
		prefix = "not"
	case trueConst:
		prefix = "true"
	case falseConst:
		prefix = "false"
	case End:
		prefix = "$$$"
	default:
		panic("can't happen")
	}
	return prefix + postfix
}

func (r rxc) String() string {
	return "(" + r.row.String() + "," + r.col.String() + ")"
}
