package eval

import (
	"fmt"
	"github.com/cuiweixie/toylang/ir"
	"github.com/cuiweixie/toylang/syntax"
	"strconv"
)

func EvalFile(name string) error {
	file, err := syntax.ParseFile(name)
	if err != nil {
		return err
	}
	nodes := ir.GenAst(file)
	c := loadNodes(nodes)
	err = registGlobalBultin(c.Scope)
	if err != nil {
		return err
	}
	node := GetFuncByName(c, "main")
	EvalNode(c, node, nil)
	return nil
}


func GetFuncByName(c *EvalCtx, name string) ir.Node {
	f, ok := c.Scope.Def[name]
	if ok {
		if v, ok := f.(*Var); ok && v.Type == FUNC {
			return v.Func
		}
	}
	return nil
}

func NewEvalCtx(scope *Scope) *EvalCtx {
	c := new(EvalCtx)
	c.Scope = scope
	return c
}

func registGlobalBultin(scope *Scope) error {
	var printVar Var
	printVar.Type = FUNC
	printVar.BuiltIn = func(c *EvalCtx, args []*Var) {
		for _, arg := range args {
			PrintVar(arg)
		}
	}
	scope.Def["print"] = &printVar
	return nil
}


func PrintVar(v *Var) {
	switch v.Type {
	case NUM:
		fmt.Printf("%f", v.NumVal)
	case STRING:
		fmt.Printf("%s", v.StringVal)
	case FUNC:
		fmt.Printf("func[%s]", v.Func.(*ir.Func).FuncName)
	case BOOL:
		fmt.Printf("%v", v.BoolVal)
	case NIL:
		fmt.Print("nil")
	}
}

type Scope struct {
	Parent *Scope
	Def map[string]Def
}

type Def interface {
	aVar()
}

type Var struct {
	Name string
	NumVal float64
	StringVal string
	BoolVal bool
	Func ir.Node
	Type VarType
	BuiltIn func(c *EvalCtx, args []*Var)
	Def
}

type EvalCtx struct {
	Scope      *Scope
	Result     []*Var
	isReturn   bool
	isContinue bool
	isBreak    bool
}

type VarType int

const (
	_ VarType = iota
	BOOL
	NUM
	STRING
	NIL
	FUNC
)


func loadNodes(nodes []ir.Node) *EvalCtx {
	scope := &Scope{
		Parent: nil,
		Def:    make(map[string]Def),
	}
	c := NewEvalCtx(scope)
	for _, node := range nodes {
		switch node := node.(type) {
		case *ir.VarDecl:
			c = EvalNode(c, node, nil)
		case *ir.Func:
			c.Scope.Def[node.FuncName] = &Var{Name: node.FuncName, Type: FUNC, Func: node}
		default:
			panic("no support node")
		}
	}
	return c
}

func EvalNode(c *EvalCtx, node ir.Node, args []*Var) *EvalCtx {
	c.Result = nil
	if c.isReturn || c.isBreak || c.isContinue {
		return c
	}
	switch node := node.(type) {
	case *ir.Name:
		_varDef := c.LookupVar(node.Name)
		_var, ok := _varDef.(*Var)
		if !ok {
			panic(fmt.Sprintf("undefine var %s", node.Name))
		}
		c.Result = append(c.Result, _var)
	case *ir.VarDecl:
		for i := range node.Lhs {
			c = EvalNode(c, node.Rhs[i], nil)
			c.Scope.Def[node.Lhs[i]] = c.Result[0]
		}
	case *ir.Literal:
		var result Var
		switch node.Type {
		case syntax.TNUM:
			result.Type = NUM
			result.NumVal, _ = strconv.ParseFloat(node.Val, 64)
		case syntax.TSTRING:
			result.Type = STRING
			result.StringVal = node.Val
		}
		c.Result = []*Var{&result}
	case *ir.CallExpr:
		funcName := node.Name
		def := c.LookupVar(funcName)
		varItem, _ := def.(*Var)
		var args []*Var
		for i := 0; i<len(node.Args); i++ {
			c = EvalNode(c, node.Args[i], nil)
			for _, v := range c.Result {
				args = append(args, v)
			}
		}
		if varItem.BuiltIn != nil {
			varItem.BuiltIn(c, args)
		} else {
			c = EvalNode(c, varItem.Func, args)
		}
		c.isReturn = false
	case *ir.Func:
		c.PushScope()
		for _, bn := range node.Body {
			c = EvalNode(c, bn, nil)
			if c.isReturn {
				c.isReturn = false
				break
			}
		}
		c.PopScope()
	case *ir.ReturnStmt:
		var results []*Var
		for i := 0; i<len(node.Returns); i++ {
			c = EvalNode(c, node.Returns[i], nil)
			for _, v := range c.Result {
				results = append(results, v)
			}
		}
		c.Result = results
		c.isReturn = true
	case *ir.BlockStmt:
		for _, stmt := range node.Stmts {
			c = EvalNode(c, stmt, nil)
			if c.isReturn {
				return c
			}
		}
	case *ir.BinaryExpr:
		c = EvalNode(c, node.Lhs, nil)
		leftVar := c.Result[0]
		if len(c.Result) != 1 {
			panic("expr op != 1")
		}
		c = EvalNode(c, node.Rhs, nil)
		rightVar := c.Result[0]
		if len(c.Result) != 1 {
			panic("expr op != 1")
		}
		result := GetBinaryOpResult(node.Op, leftVar, rightVar)
		if result == nil {
			panic("get binary op result error")
		}
		c.Result = []*Var{result}
	case *ir.IfStmt:
		c = EvalNode(c, node.Cond, nil)
		if len(c.Result) != 1 {
			panic("condition result != 1")
		}
		if c.Result[0].Type != BOOL {
			panic("error result")
		}
		c.PushScope()
		if c.Result[0].BoolVal {
			c = EvalNode(c, node.Body, nil)
		} else {
			c = EvalNode(c, node.Else, nil)
		}
		c.PopScope()
	case *ir.ForStmt:
		c.PushScope()
		c = EvalNode(c, node.Init, nil)
		for {
			c = EvalNode(c, node.Cond, nil)
			if len(c.Result) != 1 {
				panic("condition result != 1")
			}
			if c.Result[0].Type != BOOL {
				panic("error result")
			}
			if !c.Result[0].BoolVal {
				break
			}
			c = EvalNode(c, node.Body, nil)
			if c.isBreak {
				c.isBreak = false
				break
			}
			if c.isContinue {
				c.isContinue = false
			}
			c = EvalNode(c, node.Post, nil)
		}
		c.PopScope()
	case *ir.AssignStmt:
		for i:=0; i<len(node.Lhs); i++ {
			_varDef := c.LookupVar(node.Lhs[i])
			if _varDef == nil {
				panic(fmt.Sprintf("undefine var %s", node.Lhs[i]))
			}
			_var, _ := _varDef.(*Var)
			c = EvalNode(c, node.Rhs[i], nil)
			*_var = *c.Result[0]
		}
	case *ir.ContinueStmt:
		c.isContinue = true
	case *ir.BreakStmt:
		c.isBreak = true
	}
	return c
}

func (c *EvalCtx) PushScope() {
	curScope := &Scope{
		Parent:	c.Scope,
		Def:    make(map[string]Def),
	}
	c.Scope = curScope
}

func (c *EvalCtx) PopScope() {
	c.Scope = c.Scope.Parent
}

func GetBinaryOpResult(op syntax.Op, leftVar *Var, rightVar *Var) *Var {
	switch op {
	case syntax.OpPLUS:
		if leftVar.Type != rightVar.Type {
			panic("type not equal")
		}
		switch leftVar.Type {
		case NUM:
			return &Var{
				NumVal:    leftVar.NumVal + rightVar.NumVal,
				Type:      NUM,
			}
		case STRING:
			return &Var{
				StringVal: leftVar.StringVal + rightVar.StringVal,
				Type:      STRING,
			}
		default:
			panic("error type +")
		}
	case syntax.OpMINUS:
		if leftVar.Type != rightVar.Type || leftVar.Type != NUM  {
			panic("type not equal in minus")
		}
		return &Var{
			NumVal: leftVar.NumVal - rightVar.NumVal,
			Type:      NUM,
		}
	case syntax.OpMUL:
		if leftVar.Type != rightVar.Type || leftVar.Type != NUM  {
			panic("type not equal mul")
		}
		return &Var{
			NumVal: leftVar.NumVal * rightVar.NumVal,
			Type:      NUM,
		}
	case syntax.OpDiv:
		if leftVar.Type != rightVar.Type || leftVar.Type != NUM  {
			panic("type not equal div")
		}
		return &Var{
			NumVal: leftVar.NumVal / rightVar.NumVal,
			Type:      NUM,
		}
	case syntax.OpEQ:
		if leftVar.Type != rightVar.Type || leftVar.Type != NUM  {
			panic("type not equal div")
		}
		return &Var{
			BoolVal: leftVar.NumVal == rightVar.NumVal,
			Type:      BOOL,
		}
	case syntax.OpLEQ:
		if leftVar.Type != rightVar.Type || leftVar.Type != NUM  {
			panic("type not equal div")
		}
		return &Var{
			BoolVal: leftVar.NumVal <= rightVar.NumVal,
			Type:      BOOL,
		}
	case syntax.OpLT:
		if leftVar.Type != rightVar.Type || leftVar.Type != NUM  {
			panic("type not equal div")
		}
		return &Var{
			BoolVal: leftVar.NumVal < rightVar.NumVal,
			Type:      BOOL,
		}
	case syntax.OpGEQ:
		if leftVar.Type != rightVar.Type || leftVar.Type != NUM  {
			panic("type not equal div")
		}
		return &Var{
			BoolVal: leftVar.NumVal >= rightVar.NumVal,
			Type:      BOOL,
		}
	case syntax.OpGT:
		if leftVar.Type != rightVar.Type || leftVar.Type != NUM  {
			panic("type not equal div")
		}
		return &Var{
			BoolVal: leftVar.NumVal > rightVar.NumVal,
			Type:      BOOL,
		}
	}
	return nil
}

func (c *EvalCtx) LookupVar(name string) Def {
	for scope := c.Scope; scope != nil; scope = scope.Parent {
		if v, ok := scope.Def[name]; ok {
			return v
		}
	}
	return nil
}