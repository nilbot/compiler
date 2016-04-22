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
                SecretSource src = (SecretSource)data;  
                src.push("+");  
                ++indent;            
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorSubtraction node, Object data) {
                System.out.println(indentString() + node);
                SecretSource src = (SecretSource)data;
                src.push("-");
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorMultiplication node, Object data) {
                System.out.println(indentString() + node);
                SecretSource src = (SecretSource)data;
                src.push("*");
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

        // advance from 3 to 4
        public Object visit(VisitorTernaryIfElse node, Object data) {
                System.out.println(indentString() + node);
                SecretSource src = (SecretSource)data;
                if (src.state == 3) { finishCmp(src); }
                src.state = 4;
                src.push("?");
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                String rst = src.genTernary();
                src.push(rst);
                src.state = 5;
                return data;
        }

        // advance from 2 to 3
        public Object visit(VisitorLtNext node, Object data) {
                System.out.println(indentString() + node);
                SecretSource src = (SecretSource)data;
                src.state = 3;
                src.push("<");
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        // advance from 2 to 3
        public Object visit(VisitorEqNext node, Object data) {
                System.out.println(indentString() + node);
                SecretSource src = (SecretSource)data;
                src.state = 3;
                src.push("=");
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        // advance from 2 to 3
        public Object visit(VisitorGtNext node, Object data) {
                System.out.println(indentString() + node);
                SecretSource src = (SecretSource)data;
                src.state = 3;
                src.push(">");
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                return data;
        }

        public Object visit(VisitorAssignment node, Object data) {
                System.out.println(indentString() + node);
                SecretSource src = (SecretSource)data;
                src.state = 1;
                src.push(":=");
                ++indent;
                data = node.childrenAccept(this, data);
                --indent;
                if (src.state == 3) { // 2, 3, 5 are finishing states. just fyi
                        finishCmp(src);
                }
                src.genAssign();
                src.state = 1;
                if (!src.empty()) {
                        throw new IllegalStateException("stack not empty after"+
                        " finishing one assignment");
                }
                return data;
        }
        
        ////////////////////////////////////////////////////////
        //// ---- most important two handlers, i think ---- ////
        ////////////////////////////////////////////////////////

        public Object visit(VisitorId node, Object data) {
                System.out.println(indentString() + node);

                SecretSource src = (SecretSource)data;
                String id = node.getName();
                switch (src.state) {
                        case 0:
                                // do nothing
                        break;
                        case 1:
                                // top() can only be assignment op
                                src.push(id);
                                src.state = 2;
                        break;
                        case 2:
                        case 3:
                        case 4:
                                if (arithmeticOp(src.top())) {
                                        String op = src.pop();
                                        String src1 = src.pop();
                                        
                                        String t = src.genA(op,src1,id);
                                        src.push(t);
                                } else {
                                       src.push(id); 
                                }
                        break;
                }


                return data;
        }

        public Object visit(VisitorNum node, Object data) {
                System.out.println(indentString() + node);

                SecretSource src = (SecretSource)data;
                int val = node.getValue();
                switch (src.state) {
                        case 0:
                                // error
                                throw new IllegalStateException("should not happen");
                        case 1:
                                // error
                                throw new IllegalStateException("should not happen");
                        case 2:
                        case 3:
                        case 4:
                                if (arithmeticOp(src.top())) {
                                        String op = src.pop();
                                        String src1 = src.pop();
                                        
                                        String t = src.genA(op,src1,val+"");
                                        src.push(t);
                                } else {
                                       src.push(val+""); 
                                }
                        break;
                }

                return data;
        }
        
        private boolean arithmeticOp(String op) {
                if (op.equals("+")||op.equals("-")||op.equals("*")) {
                        return true;
                }
                return false;
        }
        
        private boolean spaceship(String cmp) {
                if (cmp.equals("<")||cmp.equals("=")||cmp.equals(">")) {
                        return true;
                }
                return false;
        }
        
        private void finishCmp(Object data) {
                SecretSource src = (SecretSource)data;
                if (spaceship(src.peek(2))) {
                        String a2 = src.pop();
                        String cmp = src.pop();
                        String a1 = src.pop();
                        String rst = src.genCmp(cmp,a1,a2);
                        src.push(rst);
                } else {
                        System.err.println(src.debug());
                        throw new IllegalStateException(
                                "state3: 2 As after spaceship on stack"
                        );
                }
        }
}