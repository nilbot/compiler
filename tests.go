package compiler

type TrieTestCase struct {
	Input struct {
		Identifier string
		Flag       FlagVar
	}
	Expected TrieResult
}

func equals(target, expected TrieResult) bool {
	return target == expected
}
