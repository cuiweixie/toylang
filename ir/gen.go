package ir

import "github.com/cuiweixie/toylang/syntax"


type irgen struct {

}


func (irgen *irgen) VarDecl(d *syntax.VarDecl) Node {
	declNode := new(VarDecl)
	declNode.Lhs = d.Lhs
	for i := 0; i<len(d.Lhs); i++ {
		declNode.Rhs = append(declNode.Rhs, irgen.Expr(d.Rhs[i]))
	}
	return declNode
}

func (irgen *irgen) Expr(e syntax.Expr) Node {
	switch e := e.(type) {
	case *syntax.Name:
		n := new(Name)
		n.Name = e.Name
		return n
	case *syntax.Literal:
		n := new(Literal)
		n.Type = e.Type
		n.Val = e.Val
		return n
	case *syntax.BinaryExpr:
		n := new(BinaryExpr)
		n.Op = e.Op
		n.Lhs = irgen.Expr(e.Lhs)
		n.Rhs = irgen.Expr(e.Rhs)
		return n
	case *syntax.CallExpr:
		n := new(CallExpr)
		n.Name = e.Name
		for _, expr := range e.Args {
			n.Args = append(n.Args, irgen.Expr(expr))
		}
		return n
	default:
		panic("unknown expr type")
	}
}


func (irgen *irgen) FuncDecl(f *syntax.FuncDecl) Node {
	funcNode := new(Func)
	funcNode.FuncName = f.FuncName
	funcNode.Args = f.Args
	for _, stmt := range f.Body {
		funcNode.Body = append(funcNode.Body, irgen.Stmt(stmt))
	}
	return funcNode
}

func(irgen *irgen) Stmt(stmt syntax.Stmt) Node {
	switch stmt := stmt.(type) {
	case *syntax.DeclStmt:
		node := irgen.VarDecl(stmt.Decl.(*syntax.VarDecl))
		return node
	case *syntax.AssignStmt:
		node := new(AssignStmt)
		node.Lhs = stmt.Lhs
		for _, expr := range stmt.Rhs {
			node.Rhs = append(node.Rhs, irgen.Expr(expr))
		}
		return node
	case *syntax.IfStmt:
		node := new(IfStmt)
		node.Cond = irgen.Expr(stmt.Cond)
		node.Body = irgen.Stmt(stmt.Body)
		node.Else = irgen.Stmt(stmt.Else)
		return node
	case *syntax.ForStmt:
		node := new(ForStmt)
		node.Init = irgen.Stmt(stmt.Init)
		node.Cond = irgen.Expr(stmt.Cond)
		node.Post = irgen.Stmt(stmt.Post)
		node.Body = irgen.Stmt(stmt.Body)
		return node
	case *syntax.BlockStmt:
		node := new(BlockStmt)
		for _, oneStmt := range stmt.Stmts {
			node.Stmts = append(node.Stmts, irgen.Stmt(oneStmt))
		}
		return node
	case *syntax.ReturnStmt:
		node := new(ReturnStmt)
		for _, oneExpr := range stmt.Returns {
			node.Returns = append(node.Returns, irgen.Expr(oneExpr))
		}
		return node
	case *syntax.CallStmt:
		node := new(CallExpr)
		callExpr := stmt.Call.(*syntax.CallExpr)
		node.Name = callExpr.Name
		for _, expr := range callExpr.Args {
			node.Args = append(node.Args, irgen.Expr(expr))
		}
		return node
	case *syntax.BreakStmt:
		node := new(BreakStmt)
		return node
	case *syntax.ContinueStmt:
		node := new(ContinueStmt)
		return node
	}
	return nil
}

func GenAst(file *syntax.File)[]Node {
	var nodes []Node
	var irgen irgen
	for _, decl := range file.Decl {
		switch d := decl.(type) {
		case *syntax.VarDecl:
			nodes = append(nodes, irgen.VarDecl(d))
		case *syntax.FuncDecl:
			nodes = append(nodes, irgen.FuncDecl(d))
		default:
			panic("unknown decl")
		}
	}
	return nodes
}