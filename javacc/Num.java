public class Num extends Exp {
        private int value;
        public Num(int id) { super(id); }
        public void setValue(String str) {value = Integer.parseInt(str);}
        public String toString() { return "Num: " + value; }
}