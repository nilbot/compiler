/*Compiler is a project that implements a compiler for a simple language.

Basics

This Compiler implementation is the result of open requirement of courseware
assignments for module Compiler Construction COMP30330 at UCD CSI, that any
programming language can be used as long as it implements the requirements.

The Compiler is implemented after chronilogical order the requirements from
the courseware were released.
	- Trie based Symbol table data structure
	- scanner / lexer
	- grammar parser
	- code translation
To aid the testing and implementation, these components were also introduced:
    - helper utilities
    - test and benchmarking code

General Notes and Caveats

The requirement of the module asks for handling characters range within ASCII
table. However the notion of `characters` is ambiguous in data presentation
level: depending on encoding scheme the data input chooses to adopt, a
`unicde` labled input text might have amibigous code points scheme that
requires normalization before proceeding to text processing. Also because Go
source code is always in UTF-8, it just makes sense to normalize texts input
either automatically using helper text normalization routine or manually saved
in UTF-8 encoding.

Specifically, we want to eliminate all multi rune
characters in the input text because they mess with rune based processing and
they are not in the ASCII range anyway.
See https://en.wikipedia.org/wiki/Unicode_equivalence for
background technicalities of this mess.

Trie based Symbol Table Data Structure

The Symbol table is a lookup table for known Symbols, or identifiers. To
construct this lookup table a prefix tree (Trie) data structure is used to
process the input string.

To process ASCII characters the Trie features a bounded tree with the bound
equals to 128, so that the look up process has a temporal complexity of O(d),
with the d representing the depth of the tree, also as known as the length of
the input string.

The Trie stores information about the character and their meaning along the
path from the root: whether it completes word to form an Identifier or not;
and if it does, the index of the word in the Symbol table.

*/
package compiler
