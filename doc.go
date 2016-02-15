/*Package compiler contains code and manuscripts that works as lexer

Basics

Documentation for compiler implementation for practical assignment
for module Compiler Construction COMP30330 at UCD CSI

This package contains modules, due to restriction on folders imposed by
lecturer these modules cannot be made into subpackages.

	- lexer
	- scanner
	- datastructure - trie
	- utilities

Trie data structure

To process ASCII characters the Trie features a bounded tree with the bound
equals to 128, so that the look up process has a temporal complexity of O(d),
with the d presenting the depth of the tree, also as known the length of the
look up input string.
*/
package compiler
