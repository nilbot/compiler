public class Id extends Exp {
        private String name;
        public Id(int id) { super(id); }
        public void setName(String n) { name = n; }
        public String toString() { return "Id: " + name; }
}