# Author: Ersi Ni
TEXPAD_BUILD_DIR=.texpadtmp/
GODOC_OBJ=godoc.html
GODOCPDF_OBJ=godoc.pdf
GODOCPDF_FLAG=--dpi=120 --latex-engine=xelatex
LEXER_XELATEX_SRC=latex/lexer_report.tex
PARSER_XELATEX_SRC=latex/parser_report.tex
LATEX_OBJS=*.out *.log *.blg *.bbl *.aux $(ALL_BIBTEX_OBJ)
LEXER_BIBTEX_OBJ=lexer_reportNotes.bib
PARSER_BIBTEX_OBJ=parser_reportNotes.bib
ALL_BIBTEX_OBJ=$(LEXER_BIBTEX_OBJ) $(PARSER_BIBTEX_OBJ)
DFS_PARSER_OUTPUT_TXT=parser.output.txt
FINALOUTCOME=cpu.pdf mem.pdf lexer_report.pdf coverage.html godoc.html

test:
	go test

uf:
	go test unionfind.go unionfind_test.go

parser:
	go test parse.go parse_test.go

st:
	go test trie.go trie_test.go

lexer:
	go test lex.go lex_test.go tests.go trie.go trie_test.go

print_parser_output: clean_parser_output
	go test parse*.go -v > $(DFS_PARSER_OUTPUT_TXT)

godo:
	godoc -tabwidth 2 -html . > $(GODOC_OBJ)

report:
	cd $(TEXPAD_BUILD_DIR) && \
	xelatex ../$(XELATEX_SRC) && bibtex $(BIBTEX_SRC) && \
	xelatex ../$(XELATEX_SRC) && xelatex ../$(XELATEX_SRC) && \
	mv $(BIBTEX_SRC).pdf ../ && \
	cd ..

bench: clean_test
	go test -cover -c && \
	./compiler.test -test.bench=. -test.cpuprofile=cpu.out \
	-test.coverprofile=cover.out -test.memprofile=mem.out && \
	go tool cover -html cover.out -o coverage.html && \
	go tool pprof -pdf -output=mem.pdf compiler.test mem.out && \
	go tool pprof -pdf -output=cpu.pdf compiler.test cpu.out

clean:
	rm -f $(FINALOUTCOME) *.log *.out

clean_test:
	rm -f compiler.test cpu.out mem.out cover.out 

clean_parser_output:
	rm -f $(DFS_PARSER_OUTPUT_TXT)
