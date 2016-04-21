public class VisitorId extends SimpleNode {
        private String name;
        public VisitorId(int id) {super(id);}
        public Object jjtAccept(TheGrandFinaleVisitor visitor, Object data) {
                return visitor.visit(this, data);
        }
        public void setName(String n) { name = n; }
        public String toString() { return "Id: " + name; }
}