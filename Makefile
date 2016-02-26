DOC_BUILD_DIR=.texpadtmp/
GODOC_OBJ=godoc.html
GODOCPDF_OBJ=godoc.pdf
GODOCPDF_FLAG=--dpi=120 --latex-engine=xelatex
XELATEX_SRC=report.tex
OBJS=*.out *.log *.blg *.bbl *.aux *.bib
BIBTEX_SRC=report
BIBTEX_OBJ=reportNotes.bib
BUILT=*.pdf
doc:
	godoc -tabwidth 2 -html . > $(DOC_BUILD_DIR)$(GODOC_OBJ)
	# pandoc $(DOC_BUILD_DIR)$(GODOC_OBJ) -o $(GODOCPDF_OBJ) $(GODOCPDF_FLAG)

bench:
	go test -cover -c && ./compiler.test -test.bench=. -test.cpuprofile=cpu.out -test.coverprofile=cover.out -test.memprofile=mem.out -test.outputdir=profiling.test
	go tool cover -html profiling.test/cover.out -o .texpadtmp/coverage.html
	# pandoc .texpadtmp/coverage.html -o coverage.pdf $(GODOCPDF_FLAG)
	go tool pprof -pdf -output=mem.pdf compiler.test profiling.test/mem.out
	go tool pprof -pdf -output=cpu.pdf compiler.test profiling.test/cpu.out
clean:
	rm $(BUILT)
