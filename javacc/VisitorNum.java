public class VisitorNum extends SimpleNode {
        private int value;
        public VisitorNum(int id) { super(id); }
        public void setValue(String str) {value = Integer.parseInt(str);}
        public int getValue() {return value;}
        public String toString() { return "Num: " + value; }
        public Object jjtAccept(TheGrandFinaleVisitor visitor, Object data) {
                return visitor.visit(this, data);
        }
}