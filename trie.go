package compiler

import "errors"

//FlagVar is enum type for flags
type FlagVar int

const (
	//Static Trie: No new nodes
	Static FlagVar = iota
	//Dynamic Trie: New node are added
	Dynamic
)

//NewTrieNode constructs a TrieNode and return the pointer to it.
//It creates 128 nil pointers and of children nodes and they will occupy 1KB
//(8 bytes * 128)
func NewTrieNode() *TrieNode {
	return &TrieNode{
		Children: make([]*TrieNode, 128, 128),
		HasWord:  false,
		Index:    0,
	}
}

//TrieNode states whether it is a finishing state (HasWord) and the index of
//the word from root to current node inside the symbol table. Each node has a
//fixed branch bound of 128 children nodes (range of ascii characters)
type TrieNode struct {
	Children []*TrieNode
	HasWord  bool
	Index    int
}

//Process scan the input string character by character, depending on the flag
//it updates the Trie based SymbolTable data structure and return the index of
//the input symbol in the SymbolTable when appropiate.
func (s *SymbolTable) Process(text string, flag FlagVar) (int, error) {
	current := s.TrieHead
	for _, r := range text {
		haz, err := current.has(r)
		if err != nil {
			return -1, err // crash here? or recover?
		}
		if !haz {
			switch flag {
			case Dynamic:
				node := NewTrieNode()
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
	// reached eof; if dynamic then assign word
	if !current.HasWord && flag == Dynamic {
		current.HasWord = true
		current.Index = len(s.Table)
		s.Table = append(s.Table, text)
		return current.Index, nil
	}
	if current.HasWord {
		return current.Index, nil
	}
	// reached eof, no word && flag = static
	return -1, nil
}

func validASCII(t int) bool {
	if t >= 0 && t < 128 {
		return true
	}
	return false
}

//SymbolTable contains a pointer to a Trie and maintains a slice of symbols
type SymbolTable struct {
	TrieHead *TrieNode
	Table    []string
}

//NewSymbolTable constructs a SymbolTable with a new Trie Head and an empty
//string slice with 0 length
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
		return false, errors.New("while checking nodes, " +
			"character is not in valid range of ascii")
	}
	ret := t.Children[idx] != nil
	return ret, nil
}

func (t *TrieNode) get(r rune) *TrieNode {
	idx := int(r)
	return t.Children[idx]
}
