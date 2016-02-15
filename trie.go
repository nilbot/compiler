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

func ProcessIdentifier(text string, flag FlagVar) (rst TrieResult, err error) {
	trie := NewTrieNode()
	return trie.Process(text, flag)
}
