package syntax


type File struct {
	Decl []Decl
}

type Node interface {
	aNode()
}

type Decl interface {
	Node
	aDecl()
}

type VarDecl struct {
	Lhs []string
	Rhs []Expr
	Decl
}

type FuncDecl struct {
	FuncName string
	Args []string
	Body []Stmt
	Decl
}

type Expr interface {
	Node
	aExpr()
}

type Name struct {
	Name string
	Expr
}

type Literal struct {
	Val string
	Type LiteralType
	Expr
}

//go:generate stringer -type LiteralType -linecomment node.go
type LiteralType int
const (
	TNUM LiteralType = iota + 1
	TSTRING
	TNIL
)
type BinaryExpr struct {
	Op Op
	Lhs, Rhs Expr
	Expr
}
//go:generate stringer -type Op -linecomment node.go
type Op int

const (
	_ Op = iota
	OpPLUS
	OpMINUS
	OpMUL
	OpDiv
	OpEQ
	OpLT
	OpGT
	OpLEQ
	OpGEQ
)

type Stmt interface {
	Node
	aStmt()
}

type AssignStmt struct {
	Lhs [] string
	Rhs []Expr
	Stmt
}

type CallExpr struct {
	Name string
	Args []Expr
	Expr
}

type DeclStmt struct {
	Decl Decl
	Stmt
}

type CallStmt struct {
	Call Expr
	Stmt
}

type ReturnStmt struct {
	Returns [] Expr
	Stmt
}

type BlockStmt struct {
	Stmts []Stmt
	Stmt
}

type IfStmt struct {
	Cond Expr
	Body Stmt
	Else Stmt
	Stmt
}

type ForStmt struct {
	Init Stmt
	Cond Expr
	Post Stmt
	Body Stmt
	Stmt
}

type BreakStmt struct {
	Stmt
}

type ContinueStmt struct {
	Stmt
}