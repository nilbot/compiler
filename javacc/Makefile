ONE=TheGrammarLex
TWO=TheGrammarParse
THREE=TheGrammarAST
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
TARGET?=$(THREE)	

clean:
	mkdir -p build && cd build && rm -f *.java *.jj *.class *.m* *.jar

prep:
	cp *.java build/

build: clean prep
	$(ENV) && cd build && jjtree ../$(TARGET)$(JJT) && \
	javacc $(TARGET)$(JJ) && $(JAVA_COMPILE)

run: build
	cd build && java $(TARGET) < ../$(TARGET)$(IN)
