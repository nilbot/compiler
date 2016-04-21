import java.util.*;
public class SecretSource {
        int phase = 0; // 0 is declaration phase, 1 is assignment phase
        Map<String,String> identifiers = new HashMap<>();
        int varEnum = 0;
        Stack<String> stack = new Stack<>();
        public String newVar() {
                return "t"+ ++varEnum;
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
}