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

        public Object visit(VisitorA node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        // handles phase 0; I don't care what type it is.
        public Object visit(VisitorDeclaration node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                ((SecretSource)data).state = 0; // ensure
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
        
        ////////////////////////////////////////////////////////
        //// ---- most important two handlers, i think ---- ////
        ////////////////////////////////////////////////////////

        public Object visit(VisitorId node, Object data) {
                System.out.println(indentString() + node);

                SecretSource src = (SecretSource)data;
                String id = node.getName();
                String var = src.getVarNameForId(id);
                switch (src.state) {
                        case 0:
                                // do nothing
                        break;
                        case 1:
                                // top() can only be assignment op
                                src.push(var);
                                src.state = 2;
                        break;
                        case 2:
                                if (arithmeticOp(src.top())) {
                                        String op = src.pop();
                                        String src1 = src.pop();
                                        
                                        src.generateA(op,src1,var);
                                } else {
                                       src.push(var); 
                                }
                        break;
                }


                return data;
        }

        public Object visit(VisitorNum node, Object data) {
                System.out.println(indentString() + node);
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }
        
        private boolean arithmeticOp(String op) {
                if (op.equals("+")||op.equals("-")||op.equals("*")) {
                        return true;
                }
                return false;
        }
}