public class ASTId extends SimpleNode {
        private String name;
        public ASTId(int id) { super(id); }
        public void setName(String n) { name = n; }
        public String toString() { return "Id: " + name; }
}