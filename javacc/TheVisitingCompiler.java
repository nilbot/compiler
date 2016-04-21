public class TheVisitingCompiler implements TheGrandFinaleVisitor {
        private int indent = 0;
        private String indentString() {
                StringBuffer sb = new StringBuffer();
                for (int i = 0; i < indent; ++i) {
                        sb.append(' ');
                }
                return sb.toString();
        }

        public Object visit(SimpleNode node, Object data) {
                System.out.println(indentString() + node +
                ": acceptor not unimplemented in subclass?");
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorBlock node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                if (data == null) data = new SecretSource();
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorAddition node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorSubtraction node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorMultiplication node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorId node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorNum node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorA node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorIntDecl node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorBoolDecl node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorTernaryIfElse node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorLtNext node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorEqNext node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorGtNext node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorAssignment node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }
}