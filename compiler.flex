UC [A-Z]
LC [a-z]
DIGIT [0-9]
WS 	[\n\t ]
ILLEGAL [^a-zA-Z\d\s;\(\)] 

%{
	int testcase;
	#include "symboltable.h"
	void parse_int(const char *input, int length);
	int which_keywords(const char *input);
	int which_id(const char *input);
	char* which_delimiter(const char *input);
	char* escape(const char *input);
%}

%x STRING_EOF

%%
			/* cosmetic */
"//"			{printf("====\nTestCase %d\n====\n",++testcase);root_node = create_node();symboltable=8;}
			/* whitespaces */
{WS}*
			/* integers */
{DIGIT}+      parse_int(yytext,yyleng);
			/* identifiers */
{UC}{LC}*	{printf("[Id] %d\n",which_keywords(yytext));}
{LC}+		{printf("[Id] %d\n",which_id(yytext));}
			/* single tokens */
"("|")"|";"	{printf("[%s] 0\n",which_delimiter(yytext));}
			/* string with escaped char '~' */
"\""		BEGIN STRING_EOF;
<STRING_EOF>.
<STRING_EOF>\n
<STRING_EOF><<EOF>>	{printf("[Error] 0\n");return 0;}

"\""((~\")*|(\n|\r)*|[^"]*)*"\""	{printf("[String] %s\n",escape(yytext));}
			/* error */
{ILLEGAL}	{printf("[Error] 0\n");}
%%


void parse_int(const char *i, int l) {
	int rst = 0;
	int iterator;
	const int maxint = 65535;
	for (iterator = 0; iterator < l; iterator++) {
		int n;
		n = i[iterator]-'0';
		if (n>=0 && n<=9) {
			if (rst<=6553 && rst>-1) {
				rst*=10;
				if ((maxint-rst)<n) {
					rst = -1;
				} else {
					rst += n;
				}
			} else {
				rst = -1;
			}
		} else {
			break;
		}
	}
	printf("[Int] %d\n",rst);
}

int which_keywords(const char *i) {
	if(strcmp(i, "Private") == 0)
	return 0;
	if(strcmp(i, "Public") == 0)
	return 1;
	if(strcmp(i, "Protected") == 0)
	return 2;
	if(strcmp(i, "Static") == 0)
	return 3;
	if(strcmp(i, "Primary") == 0)
	return 4;
	if(strcmp(i, "Integer") == 0)
	return 5;
	if(strcmp(i, "Exception") == 0)
	return 6;
	if(strcmp(i, "Try") == 0)
	return 7;
	
	return -1;
}

int which_id(const char *input) {
	insert_node(input);
	return search(input);
}

char* which_delimiter(const char* i) {
	if (strcmp(i,"(") == 0) return "Lpar";
	if (strcmp(i,")") == 0) return "Rpar";
	if (strcmp(i,";") == 0) return "Semicolon";
	return "Error";
}

char* escape(const char* i) {
	char *rst = NULL;
	int size = strlen(i);
	rst = malloc(size);
	int ptr = 0; // for rst only
	for (int idx = 1; idx < size-1; idx++) {
		if (i[idx] == '~') {
			if (idx+1<size-1) {
				rst[ptr++] = i[++idx];
			}
			continue;
		}
		rst[ptr++] = i[idx];
	}
	return rst;
}

int yywrap(void) {
	return 1;
}

int main(int argc, char *argv[]) {
    yyin = fopen(argv[1], "r");
    yylex();
    fclose(yyin);
}