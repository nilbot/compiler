default:
	flex ../compiler.flex && gcc -o lex -g *.c && cat input_all.txt | ./lex > outputs.txt && cat outputs.txt
clean:
	rm -f lex outputs.txt *.c

