public class ASTNum extends SimpleNode {
        private int value;
        public ASTNum(int id) { super(id); }
        public void setValue(String str) {value = Integer.parseInt(str);}
        public String toString() { return "Num: " + value; }
}