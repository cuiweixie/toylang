package ir

import "github.com/cuiweixie/toylang/syntax"

type Node interface {
	aNode()
}

type VarDecl struct {
	Lhs []string
	Rhs []Node
	Node
}

type Func struct {
	FuncName string
	Args   []string
	Body   []Node
	Node
}

type Name struct {
	Name string
	Node
}

type Literal struct {
	Val  string
	Type syntax.LiteralType
	Node
}

type BinaryExpr struct {
	Op syntax.Op
	Lhs, Rhs Node
	Node
}

type AssignStmt struct {
	Lhs []string
	Rhs  []Node
	Node
}

type CallExpr struct {
	Name string
	Args []Node
	Node
}

type IfStmt struct {
	Cond Node
	Body Node
	Else Node
	Node
}

type ForStmt struct{
	Init Node
	Cond Node
	Post Node
	Body Node
	Node
}

type BreakStmt struct {
	Node
}

type ContinueStmt struct {
	Node
}

type BlockStmt struct {
	Stmts []Node
	Node
}

type ReturnStmt struct {
	Returns []Node
	Node
}