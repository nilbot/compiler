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

const input1 = `this little piggy`
const input2 = `what has happened`
const input3 = `in Private say one thing in Public say another`
const input4 = `when speaking Privately to (12 or more) be discreet`
const input5 = `there are "strings of things" and there are "chickens
with
wings"
`
const input6 = `over and over and over again`
const input7 = `888 is quite ok but 88888 is "not ok"`
const input8 = `PublicAndPrivateShouldNotBeConfused`
const input9 = `when can ~ appear in a string? "Here An~d~ Th~~~~~~ere"
but "some string
`
