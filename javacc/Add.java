public class Add extends Exp {
        public Add(int id) {super(id);}
        private String op = "";
        public void setOp(String op) {this.op = op;}
        public String toString() {return this.op;}
}