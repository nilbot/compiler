public class BinaryExp extends Exp {
    String op;
    Exp left, right;
    BinaryExp(int id) {super(id);}
    void Set(String o, Exp l, Exp r) {op = o; left = l; right = r;}
    public String toString() {return "(" + op + " " + left + " " + right + ")";}
}