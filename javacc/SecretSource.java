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
        //      and when 2nd A pushed onto stack, immediately pop 3 times
        //      and push the result back onto stack 
        // 4 is advancing with ternary conditional, push ? onto stack
        //      when 2 As has been pushed after the ? immediately pop 4 times
        //      eval the FILO (4th item) A, if true eval 2nd item, else eval 3rd
        //      item. push the result back onto stack, advance to 5
        // 5 is cleanup op. pop 3 times, discard the assignment op, store (3rd)
        //      LIFO with 1st poped item. reset phase back to 1.
        //      stack should be empty at this point.
        int state = 0;
        Map<String,String> identifiers = new HashMap<>();
        int varEnum = 0;
        Stack<String> stack = new Stack<>();
        private String newVar() {
                return "t"+ varEnum++;
        }
        public String top() {
                return stack.peek();
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
        
        public void generateA(String op, String src1, String src2) {
                
        }
}