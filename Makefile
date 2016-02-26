TEXPAD_BUILD_DIR=.texpadtmp/
GODOC_OBJ=godoc.html
GODOCPDF_OBJ=godoc.pdf
GODOCPDF_FLAG=--dpi=120 --latex-engine=xelatex
XELATEX_SRC=report.tex
LATEX_OBJS=$(BIBTEX_SRC).out $(BIBTEX_SRC).log $(BIBTEX_SRC).blg $(BIBTEX_SRC).bbl $(BIBTEX_SRC).aux $(BIBTEX_OBJ)
BIBTEX_SRC=report
BIBTEX_OBJ=reportNotes.bib
FINALOUTCOME=cpu.pdf mem.pdf report.pdf coverage.html godoc.html

test:
	go test -v -parallel=1

godoc:
	godoc -tabwidth 2 -html . > $(GODOC_OBJ)

report:
	cd $(TEXPAD_BUILD_DIR) && \
	xelatex ../$(XELATEX_SRC) && bibtex $(BIBTEX_SRC) && xelatex ../$(XELATEX_SRC) && xelatex ../$(XELATEX_SRC) && \
	mv $(BIBTEX_SRC).pdf ../ && \
	cd ..

bench: clean_test
	go test -cover -c && \
	./compiler.test -test.bench=. -test.cpuprofile=cpu.out -test.coverprofile=cover.out -test.memprofile=mem.out && \
	go tool cover -html cover.out -o coverage.html && \
	go tool pprof -pdf -output=mem.pdf compiler.test mem.out && \
	go tool pprof -pdf -output=cpu.pdf compiler.test cpu.out

clean:
	rm -f $(FINALOUTCOME) *.log *.out

clean_test:
	rm -f compiler.test cpu.out mem.out cover.out 
