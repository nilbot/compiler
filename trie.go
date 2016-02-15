package compiler

type FlagVar int

const (
	Dynamic FlagVar = iota
	Static
)

func NewTrieNode() *TrieNode {
	return &TrieNode{
		HasWord: false,
		Index:   0,
	}
}

type TrieNode struct {
	Children [128]*TrieNode
	HasWord  bool
	Index    uint
}

type TrieResult int

// Process goes through the input string character by character and return the
// ID of the word in the SymbolTable should the input is a word of the
// SymbolTable, or a NotFound signal if not.
//
// Depending on the flag, the receiver Trie might or might not get updated
// during the process.
func (t *TrieNode) Process(text string, flag FlagVar) (rst TrieResult, err error) {

	for i, l := 0, len(text); i < l; i++ {
		switch flag {
		case Dynamic:
		case Static:
		default:
		}
	}
	return
}

// ProcessIndentifier process the input text string return the filled Trie
// based SymbolTable.
// BUG(n): fix the signature
func ProcessIdentifier(text string, flag FlagVar) (rst TrieResult, err error) {
	trie := NewTrieNode()
	return trie.Process(text, flag)
}
