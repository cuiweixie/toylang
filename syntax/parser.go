package syntax

import (
	"fmt"
	"io/ioutil"
)

type Parser struct {
	*Scanner
	*File
}

func getFileContent(fileName string)([]byte, error) {
	bs, err := ioutil.ReadFile(fileName)
	return bs, err
}

func ParseFile(fileName string)(*File, error) {
	content, err := getFileContent(fileName)
	if err != nil {
		return &File{}, err
	}
	p := &Parser{}
	p.Scanner = NewScanner(fileName, content)
	p.File = &File{}
	err = p.fileOrNil()
	return p.File, err
}

func(p *Parser) Next() {
	p.Scanner.Next()
}

func (p *Parser) fileOrNil() error {
	for p.Scanner.tToken != EOF {
		p.Next()
		switch p.tToken {
		case _KVAR:
			d := p.VarDecl()
			p.Decl = append(p.Decl, d)
		case _KFUNC:
			d := p.funcDecl()
			p.Decl = append(p.Decl, d)
		case SEMICOLON:

		default:
			p.Error()
			return fmt.Errorf("error")
		}
	}
	return nil
}

func (p *Parser) VarDecl() Decl {
	var varDecl VarDecl
	for {
		p.Next()
		if p.Scanner.tToken != IDENT {
			break
		}
		varDecl.Lhs = append(varDecl.Lhs, p.Scanner.literal)
	}

	if !p.Want(ASSIGN) {
		panic(fmt.Sprintf("%v need assign op here", p.Scanner.Pos))
	}

	p.Next()
	for i := 0; i< len(varDecl.Lhs); i++ {
		expr := p.BinaryExpr(0)
		if i == len(varDecl.Lhs) - 1 {
			if !p.Want(SEMICOLON) {
				panic(fmt.Sprintf("%v need ;  here", p.Scanner.Pos))
			}
		} else {
			if !p.Want(COMMA) {
				panic(fmt.Sprintf("%v need , here", p.Scanner.Pos))
			}
		}
		varDecl.Rhs = append(varDecl.Rhs, expr)
	}

	return &varDecl
}

func getOpFromTToken(t TokenType) Op {
	switch t {
	case PLUS:
		return OpPLUS
	case MINUS:
		return OpMINUS
	case MUL:
		return OpMUL
	case DIV:
		return OpDiv
	case LT:
		return OpLT
	case LEQ:
		return OpLEQ
	case GT:
		return OpGT
	case GEQ:
		return OpGEQ
	case EQUAL:
		return OpEQ
	}
	return 0
}

func (p *Parser) BinaryExpr(prec Prec) Expr {
	x := p.UnaryExpr()
	for p.Scanner.isBinaryOp && p.Scanner.Prec > prec {
		op := getOpFromTToken(p.tToken)
		prec := p.Scanner.Prec
		p.Next()
		y := p.BinaryExpr(prec)
		be := &BinaryExpr{}
		be.Op = op
		be.Lhs = x
		be.Rhs = y
		x = be
	}

	return x
}

func (p *Parser) CallExpr(name string) Expr {
	var callExpr CallExpr
	callExpr.Name = name
	for {
		p.Next()
		if expr := p.BinaryExpr(0); expr != nil {
			callExpr.Args = append(callExpr.Args, expr)
		} else {
			break
		}
		if p.Scanner.tToken == RIGHTPAREN {
			break
		}
	}
	if !p.Want(RIGHTPAREN) {
		panic(fmt.Sprintf("%v: ) need here", p.Scanner.Pos))
	}
	p.Next()
	return &callExpr
}

func (p *Parser) UnaryExpr() Expr {
	if p.Scanner.tToken == IDENT {
		expr := &Name{Name:p.Scanner.literal}
		p.Next()
		if p.Scanner.tToken != LEFTPAREN {
			return expr
		}
		callExpr := p.CallExpr(expr.Name)
		return callExpr
	}
	if p.Scanner.tToken == LEFTPAREN {
		expr := p.BinaryExpr(0)
		if !p.Want(RIGHTPAREN) {
			panic(fmt.Sprintf("%v need ) here", p.Scanner.Pos))
		}
		p.Next()
		return expr
	}
	if p.Scanner.tToken == NUM {
		expr := Literal{
			Val:  p.Scanner.literal,
			Type: TNUM,
		}
		p.Next()
		return &expr
	}
	if p.Scanner.tToken == STRING {
		expr := Literal{
			Val:  p.Scanner.literal,
			Type: TSTRING,
		}
		p.Next()
		return &expr
	}

	if p.Scanner.tToken == SEMICOLON {
		p.Next()
		return nil
	}
	return nil
}


func (p *Parser)Want(t TokenType) bool {
	return p.Scanner.tToken == t
}

func (p *Parser) funcDecl() Decl {
	p.Next()
	if !p.Want(IDENT) {
		panic(fmt.Sprintf("%v: func name need here", p.Scanner.Pos))
	}
	var funcDecl FuncDecl
	funcDecl.FuncName = p.Scanner.literal
	p.Next()
	if !p.Want(LEFTPAREN) {
		panic(fmt.Sprintf("%v: ( need here", p.Scanner.Pos))
	}
	for {
		p.Next()
		if p.Scanner.tToken == IDENT {
			funcDecl.Args = append(funcDecl.Args, p.Scanner.literal)
		} else {
			break
		}
	}
	if !p.Want(RIGHTPAREN) {
		panic(fmt.Sprintf("%v: ) need here", p.Scanner.Pos))
	}
	p.Next()
	stmts := p.funcBody()
	funcDecl.Body = stmts
	return &funcDecl
}

func (p *Parser) funcBody() []Stmt {
	var stmts []Stmt
	if !p.Want(LEFTBRACE) {
		panic(fmt.Sprintf("%v: { need here", p.Scanner.Pos))
	}
	p.Next()
	for {
		for  p.Scanner.tToken == SEMICOLON {
			p.Next()
		}
		if p.Want(RIGHTBRACE) {
			break
		}
		stmts = append(stmts, p.Stmt())
		p.Next()
	}
	p.Next()
	return stmts
}


func (p *Parser) SimpleStmt(isFor bool) Stmt {
	switch p.Scanner.tToken {
	case _KVAR:
		decl := p.VarDecl()
		var declStmt DeclStmt
		declStmt.Decl = decl
		return &declStmt
	case IDENT:
		x := p.Scanner.literal
		p.Next()
		if p.Scanner.tToken == LEFTPAREN {
			return &CallStmt{
				Call: p.CallExpr(x),
			}
		}
		return p.AssignStmt(x, isFor)
	}
	return nil
}

func (p *Parser) Stmt() Stmt {
	switch p.Scanner.tToken {
	case _KVAR:
		return p.SimpleStmt(false)
	case IDENT:
		return p.SimpleStmt(false)
	case _KIF:
		return p.IfStmt()
	case _KFOR:
		return p.ForStmt()
	case _KRETURN:
		return p.ReturnStmt()
	case _KBREAK:
		return &BreakStmt{}
	case _KCONTINUE:
		return &ContinueStmt{}
	}
	return nil
}


func (p *Parser) ReturnStmt() Stmt {
	p.Next()
	var returnStmt ReturnStmt
	for {
		expr := p.BinaryExpr(0)
		if expr != nil {
			returnStmt.Returns = append(returnStmt.Returns, expr)
		} else {
			break
		}
		if p.Scanner.tToken == SEMICOLON {
			break
		}
		p.Next()
	}

	return &returnStmt
}

func (p *Parser) IfStmt() Stmt {
	var ifStmt IfStmt
	p.Next()
	expr := p.BinaryExpr(0)
	if !p.Want(LEFTBRACE) {
		panic(fmt.Sprintf("%v: need {", p.Scanner.Pos))
	}
	ifStmt.Cond = expr
	var stmts []Stmt
	for {
		p.Next()
		if p.Want(RIGHTBRACE) {
			break
		}
		stmts = append(stmts, p.Stmt())
	}
	p.Next()
	ifStmt.Body = &BlockStmt{
		Stmts: stmts,
	}
	if p.Scanner.tToken == _KELSE {
		p.Next()
		if p.Scanner.tToken == _KIF {
			ifStmt.Stmt = p.IfStmt()
		} else {
			if !p.Want(LEFTBRACE) {
				panic(fmt.Sprintf("%v: need {", p.Scanner.Pos))
			}
			ifStmt.Cond = expr
			var stmts []Stmt
			for {
				p.Next()
				if p.Want(RIGHTBRACE) {
					break
				}
				stmts = append(stmts, p.Stmt())
			}
			ifStmt.Else = &BlockStmt{
				Stmts: stmts,
				Stmt:  nil,
			}
			p.Next()
		}
	}
	return &ifStmt
}

func (p *Parser) ForStmt() Stmt {
	p.Next()
	var forStmt ForStmt
	init := p.SimpleStmt(false)
	p.Next()
	cond := p.BinaryExpr(0)
	p.Next()
	post := p.SimpleStmt(true)
	if !p.Want(LEFTBRACE) {
		panic(fmt.Sprintf("%v: need {", p.Scanner.Pos))
	}
	forStmt.Cond = cond
	forStmt.Init = init
	forStmt.Post = post
	var stmts []Stmt
	for {
		p.Next()
		if p.Want(RIGHTBRACE) {
			break
		}
		stmts = append(stmts, p.Stmt())
	}
	forStmt.Body = &BlockStmt{Stmts: stmts}
	p.Next()
	return &forStmt
}

func (p *Parser) AssignStmt(name string, isFor bool) Stmt {
	var assignStmt AssignStmt
	if name != "" {
		assignStmt.Lhs = append(assignStmt.Lhs, name)
	}
	for {
		if p.Scanner.tToken != IDENT {
			break
		}
		assignStmt.Lhs = append(assignStmt.Lhs, p.Scanner.literal)
		p.Next()
	}

	if !p.Want(ASSIGN) {
		panic(fmt.Sprintf("%v need assign op here", p.Scanner.Pos))
	}

	p.Next()
	for i := 0; i< len(assignStmt.Lhs); i++ {
		expr := p.BinaryExpr(0)
		if i == len(assignStmt.Lhs) - 1 {
			if !isFor{
				if !p.Want(SEMICOLON) {
					panic(fmt.Sprintf("%v need ;  here", p.Scanner.Pos))
				}
			} else {
				if !p.Want(LEFTBRACE) {
					panic(fmt.Sprintf("%v need {  here", p.Scanner.Pos))
				}
			}
		} else {
			if !p.Want(COMMA) {
				panic(fmt.Sprintf("%v need , here", p.Scanner.Pos))
			}
		}
		assignStmt.Rhs = append(assignStmt.Rhs, expr)
	}
	return &assignStmt
}

func(p *Parser) Error() {
	panic(fmt.Sprintf("%s:%d:%d %s", p.Scanner.Pos.fileName, p.Scanner.Pos.line, p.Scanner.Pos.col, p.err.msg))
}
