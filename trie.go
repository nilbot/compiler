package compiler

import (
	"errors"
	"unicode/utf8"
)

type FlagVar int

const (
	Dynamic FlagVar = iota
	Static
)

func NewTrieNode() *TrieNode {
	return &TrieNode{
		Children: make([]*TrieNode, 128, 128),
		HasWord:  false,
		Index:    0,
	}
}

type TrieNode struct {
	Children []*TrieNode
	HasWord  bool
	Index    int
}

// Process goes through the input string character by character and return the
// ID of the word in the SymbolTable should the input is a word of the
// SymbolTable, or a NotFound signal if not.
//
// Depending on the flag, the receiver Trie might or might not get updated
// during the process.
func (s *SymbolTable) Process(text string, flag FlagVar) (rst int, err error) {
	l := utf8.RuneCountInString(text)
	current := s.TrieHead

	for i, r := range text {
		if hasThisChildNode, err := current.has(r); err != nil {
			return -1, err // crash here? or recover?
		} else if !hasThisChildNode {
			switch flag {
			case Dynamic:
				node := NewTrieNode()
				if i == l-1 {
					node.HasWord = true
					node.Index = len(s.Table)
					s.Table = append(s.Table, text)
				}
				current.set(r, node)
				current = node
			case Static:
				return -1, nil
			default:
				return -1, errors.New("undefined flag")
			}
		} else {
			current = current.get(r)
		}
	}
	if current.HasWord {
		return current.Index, nil
	}
	return -1, nil
}

// ProcessIndentifier process the input text string return the filled Trie
// based SymbolTable.
// BUG(n): fix the signature
func ProcessIdentifier(table *SymbolTable, text string, flag FlagVar) (rst int, err error) {
	return table.Process(text, flag)
}

func validASCII(t int) bool {
	if t >= 0 && t < 128 {
		return true
	}
	return false
}

type SymbolTable struct {
	TrieHead *TrieNode
	Table    []string
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		TrieHead: NewTrieNode(),
		Table:    make([]string, 0),
	}
}

func (t *TrieNode) set(r rune, n *TrieNode) {
	idx := int(r)
	t.Children[idx] = n
}

func (t *TrieNode) has(r rune) (bool, error) {
	idx := int(r)
	if !validASCII(idx) {
		return false, errors.New("while checking nodes, character is not in valid range of ascii")
	}
	ret := t.Children[idx] != nil
	return ret, nil
}

func (t *TrieNode) get(r rune) *TrieNode {
	idx := int(r)
	return t.Children[idx]
}
