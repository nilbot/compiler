ONE=TheGrammarLex
TWO=TheGrammarParse
THREE=TheGrammarAst
FOUR=TheGrammarVisitor
JJT=.jjt
JJ=.jj
IN=.in
JAR=.jar

ONE_TREE=$(ONE)$(JJT)
ONE_JJ=$(ONE)$(JJ)
ONE_IN=$(ONE)$(IN)

TWO_TREE=$(TWO)$(JJT)
TWO_JJ=$(TWO)$(JJ)
TWO_IN=$(TWO)$(IN)

THREE_TREE=$(THREE)$(JJT)
THREE_JJ=$(THREE)$(JJ)
THREE_IN=$(THREE)$(IN)

FOUR_TREE=$(FOUR)$(JJT)
FOUR_JJ=$(FOUR)$(JJ)
FOUR_IN=$(FOUR)$(IN)

JAVA_COMPILE=javac *.java
ENV=export PATH=$$PATH:$(abspath bin)
	

clean:
	cd build && rm -f *.java *.jj *.class *.mk *.mf *.jar

build_dir:
	mkdir -p build

one: clean build_dir
	$(ENV) && cd build && jjtree ../$(ONE_TREE) && \
	javacc $(ONE_JJ) && javac *.java && java $(ONE) < ../$(ONE_IN)

two: clean build_dir
	$(ENV) && cd build && jjtree ../$(TWO_TREE) && \
	javacc $(TWO_JJ) && javac *.java && java $(TWO) < ../$(TWO_IN)