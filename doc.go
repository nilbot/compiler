/*Package compiler is assignment project that implements a compiler for a
simple language.

Basics

This Compiler implementation is the result of open specification of courseware
assignments for module Compiler Construction COMP30330 at UCD CSI, that any
programming language can be used as long as it implements the requirements.

The Compiler is implemented after chronilogical order the requirements from
the courseware were released.
	- Trie based Symbol table data structure
	- Lexer
To aid the testing and implementation, these components were also introduced:
	- helper utilities
	- test and benchmarking code

General Notes and Caveats

The requirement of the module asks for handling characters range within ASCII
table. However the notion of "characters" is ambiguous in data presentation
level: depending on encoding scheme the data input chooses to adopt, a
"unicde" labled input text might have amibigous code points scheme that
requires normalization before proceeding to text processing. Also because Go
source code is always in UTF-8, it just makes sense to normalize texts input
either automatically using helper text normalization routine or manually saved
in UTF-8 encoding.

Specifically, we want to eliminate all multi rune
characters in the input text because they mess with rune based processing and
they are not in the ASCII range anyway.
See
	https://en.wikipedia.org/wiki/Unicode_equivalence
for background technicalities of this mess.

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
and if it does, return the index of the word in the Symbol table.

Lexer

The main reason of choosing Go to write the compiler, especially the lexer was
concurrency support and syntax support for function pointers. (C can do all
this too, minus the practicality of enormous concurrency house keeping)
Rob Pike has given a splendid talk about writing templating system using Go,
particularly the Lexical Scanning:
	https://www.youtube.com/watch?v=HxaD_trXwRE
That talk has planted deep root for my implementation.

State Machine using function pointer

The iniatialized lexer will start running with inital state and keep updating
itself: the state function does whatever it's supposed to do within the state,
triggers transition upon fulfilled condition and return the next state as a
function pointer, which is to be re-assigned to the same lexer.state variable.

The state machine diagram can be found in the appendix. In short, {start}
is the main entry point, from there the immediate possible states, depending
on leading significant character, are {integer}, {identifier}, {keyword}
(static identifers), {text} (string) and {singleToken}. Using function pointer
the {error} state doesn't need to be specified, instead the function pointer
will be set to nil, so the lexing will terminate.

Upon successful lexing a token, the lexer emits this token to its message
channel, this channel can be consumed by a potential parser, a logging / print
mechanism or any interested party. One obvious feature of this implementation
is that the parsing and lexing can be run at the same time.

Tokens

Our Lexer deals with a language that has these token types:
	- id		:	identifiers
	- int		:	intergers
	- lpar		:	left parenthesis
	- rpar		:	right parenthesis
	- semicolon	:	;s
	- string	:	strings
	- error		:	errors
With whitespace characters deemed as natural delimiters, and at cases the
transition between <id> and <int> is also deemed as delimiter. The <lpar>,
<rpar> and <semicolon> are both delimiter and tokens.
The error token will effectively terminate consumption from client.




*/
package compiler
