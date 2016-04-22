import java.util.*;
public class SecretSource {
        // 0 is declaration phase
        // 1 is assignment phase, it was advanced from 0 by seeing assignment
        //      class and will push op onto the stack. id in this phase will
        //      push the id onto the stack and advance to 2.
        // 2 is 1st exp logic. Only A() should be the recipient of this state.
        //      if top is arithmetic op (can't be else), pop twice and push
        //      result onto stack. else push current A to stack. 
        // 3 is advanced with comparator: push current spaceship op onto stack
        //      and when next state transition, pop 3 times and push the 
        //      result back onto stack 
        // 4 is advancing with ternary conditional, push ? onto stack
        //      when 2 As has been pushed after the ? immediately pop 4 times
        //      eval the FILO (4th item) A, if true eval 2nd item, else eval 3rd
        //      item. push the result back onto stack, advance to 5
        // 5 is cleanup op. pop 3 times, pop the assignment op, store (3rd)
        //      LIFO with 1st poped item. reset phase back to 1.
        //      stack should be empty at this point.
        int state = 0;
        Map<String,String> identifiers = new HashMap<>();
        int varEnum = 20;
        int lEnum = 0;
        Stack<String> stack = new Stack<>();
        Vector<Vector<String>> code = new Vector<>();
        private final String LI = "loadquantity";
        private final String LD = "load";
        private final String LABEL = "label";
        private final String NULL = "null";
        private final String J = "jump";
        private final String JE = "jumpequal";
        private final String JM = "jumpmore";
        private final String JL = "jumpless";

        private String newVar() {
                return "t"+ varEnum++;
        }
        private String newLabel(String l) {
                return "l"+l+ lEnum++;
        }
        public String top() {
                return stack.peek();
        }
        public String peek(int number) {
                Stack<String> backup = new Stack<>();
                while (number-- != 0) {
                        backup.push(stack.pop());
                }
                String rst = backup.peek();
                while(!backup.empty()) {
                        stack.push(backup.pop());
                }
                return rst;
        }
        public String pop() {
                return stack.pop();
        }
        public void push(String symbol) {
                stack.push(symbol);
        }
        public boolean empty() {
                return stack.empty();
        }
        public String getVarNameForId(String id) {
                if (identifiers.get(id) == null) {
                        String t = this.newVar();
                        identifiers.put(id,t);
                        return t;
                } else if (identifiers.get(id).equals("")) {
                        String t = this.newVar();
                        identifiers.put(id,t);
                        return t;                        
                } else {
                        return identifiers.get(id);
                }
        }

        public String debug() {
                StringBuffer sb = new StringBuffer();
                sb.append("[DEBUG] \n");
                sb.append("state: "+state+", stack:");
                sb.append(stack.toString());
                sb.append("\n");
                return sb.toString();
        }
        
        public String genA(String op, String src1, String src2) {
                try{
                        int n = Integer.parseInt(src1);
                        code.add(loadi("t1",n));
                } catch (NumberFormatException e) {
                        code.add(load("t1",src1));
                }
                try{
                        int m = Integer.parseInt(src2);
                        code.add(loadi("t2",m));
                } catch (NumberFormatException e) {
                        code.add(load("t2",src2));
                }

                Vector<String> line = new Vector<>();
                String oper = parseOp(op);
                line.add(oper);
                String tx = newVar();
                line.add(tx);
                line.add("t1");
                line.add("t2");
                code.add(line);
                return tx;
        }
        
        private String parseOp(String op) {
                switch(op) {
                        case "+":
                                return "add";
                        case "-":
                                return "subtract";
                        case "*":
                                return "multiply";
                }
                System.err.println(debug());
                throw new IllegalStateException("op is " + op);
        }
        
        public void genAssign() {
                String value = pop();
                String dest = pop();
                String op = pop();
                if (!op.equals(":=")) {
                        System.err.println(debug());
                        throw new IllegalStateException("bottom of stack"+
                        " is not assignment operator");
                }
                Vector<String> line = new Vector<>();
                code.add(load("t1",value));
                code.add(load(dest,"t1"));
        }
        
        public String genTernary() {
                //pop 4 times; LIFO: a2,a1,?,boolExp
                String a2 = stack.pop();
                String a1 = stack.pop();
                String q = stack.pop();
                String bExp = stack.pop();
                if (!q.equals("?")) {
                        System.err.println(bExp+" "+q+" "+a1+" "+a2);
                        System.err.println(debug());
                        throw new IllegalStateException("where is the ?");
                }
                String lElse = newLabel("Else");
                String lTrue = newLabel("True");
                String lEnd = newLabel("End");
                String dst = newVar();
                // t0 = 0;
                code.add(loadZero());
                // t1 = bExp;
                code.add(load("t1",bExp));
                // je Else t0,t1
                code.add(je(lElse,"t0","t1"));
                // True:
                code.add(label(lTrue));
                // return a1;
                code.add(load(dst,a1));
                // j End
                code.add(j(lEnd));
                // Else:
                code.add(label(lElse));
                // return a2;
                code.add(load(dst,a2));
                // End
                code.add(label(lEnd));
                return dst;
        }
        
        public String genCmp(String cmp, String a1, String a2) {
                // load(i) t1, a1
                try {
                        int n = Integer.parseInt(a1);
                        code.add(loadi("t1",n));
                } catch (NumberFormatException e) {
                        code.add(load("t1",a1));
                }
                // load(i) t2, a2
                try {
                        int m = Integer.parseInt(a2);
                        code.add(loadi("t2",m));
                } catch (NumberFormatException e) {
                        code.add(load("t2",a2));
                }
                String lTrue = newLabel("True");
                String lEnd = newLabel("End");
                String dst = newVar();
                switch(cmp) {
                        case "<":
                                code.add(jl(lTrue,"t1","t2"));
                        break;
                        case ">":
                                code.add(jm(lTrue,"t1","t2"));
                        break;
                        case "=":
                                code.add(je(lTrue,"t1","t2"));
                        break;
                }
                code.add(loadi(dst,0));
                code.add(j(lEnd));
                code.add(label(lTrue));
                code.add(loadi(dst,1));
                code.add(label(lEnd));
                return dst;
        }
        
        ////////////////////////////////////////////////////////////
        // loadZero load 0 into t0
        private Vector<String> loadZero() {
                Vector<String> rst = new Vector<>();
                rst.add(LI);
                rst.add("t0");
                rst.add("0");
                rst.add(NULL);
                return rst;
        }
        
        // load dest, src
        private Vector<String> load(String dst, String src) {
                Vector<String> rst = new Vector<>();
                rst.add(LD);
                rst.add(dst);
                rst.add(src);
                rst.add(NULL);
                return rst;
        }

        // load immediate dst, num
        private Vector<String> loadi(String dst, int num) {
                Vector<String> rst = new Vector<>();
                rst.add(LI);
                rst.add(dst);
                rst.add(num+"");
                rst.add(NULL);
                return rst;
        }

        private Vector<String> label(String label) {
                Vector<String> rst = new Vector<>();
                rst.add(LABEL);
                rst.add(label);
                rst.add(NULL);
                rst.add(NULL);
                return rst;
        }


        private Vector<String> j(String label) {
                Vector<String> rst = new Vector<>();
                rst.add(J);
                rst.add(label);
                rst.add(NULL);
                rst.add(NULL);
                return rst;
        }

        private Vector<String> je(String label, String src1, String src2) {
                Vector<String> rst = new Vector<>();
                rst.add(JE);
                rst.add(label);
                rst.add(src1);
                rst.add(src2);
                return rst;
        }

        private Vector<String> jl(String label, String src1, String src2) {
                Vector<String> rst = new Vector<>();
                rst.add(JL);
                rst.add(label);
                rst.add(src1);
                rst.add(src2);
                return rst;
        }

        private Vector<String> jm(String label, String src1, String src2) {
                Vector<String> rst = new Vector<>();
                rst.add(JM);
                rst.add(label);
                rst.add(src1);
                rst.add(src2);
                return rst;
        }
}