package compiler

import (
	"strings"

	"io/ioutil"
)

func equals(target, expected int) bool {
	return target == expected
}

var memoryBenchmarkTest = []struct {
	Input         string
	LengthOfInput bool
}{
	{"Apple", true},
	{"Banana", true},
	{"Testing", true},
	{"Facility", true},
	{"cAsE", true},
	{".#/po!@<>", false},
}

func getMemoryBenchmarkTests() []struct {
	Input         string
	LengthOfInput bool
} {
	content, err := ioutil.ReadFile("words.test")
	if err != nil {
		return memoryBenchmarkTest
	}
	lines := strings.Split(string(content), "\n")
	var rst []struct {
		Input         string
		LengthOfInput bool
	}
	for _, l := range lines {
		if l != "" {
			rst = append(rst, struct {
				Input         string
				LengthOfInput bool
			}{
				l, true, // TODO(n) support negative cases
			})
		}
	}
	return rst
}
