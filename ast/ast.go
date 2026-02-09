package ast

// Node is the interface all AST nodes implement.
type Node interface {
	nodeType() string
}

// Program is the root node — a list of statements.
type Program struct {
	Statements []Node
}

func (p *Program) nodeType() string { return "Program" }

// Assignment: name = expr
type Assignment struct {
	Name  string
	Value Node
}

func (a *Assignment) nodeType() string { return "Assignment" }

// StringLiteral: "hello"
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) nodeType() string { return "StringLiteral" }

// IntLiteral: 42
type IntLiteral struct {
	Value string
}

func (i *IntLiteral) nodeType() string { return "IntLiteral" }

// BoolLiteral: true, false
type BoolLiteral struct {
	Value bool
}

func (b *BoolLiteral) nodeType() string { return "BoolLiteral" }

// Identifier: name
type Identifier struct {
	Name string
}

func (id *Identifier) nodeType() string { return "Identifier" }

// FuncCall: print("hello"), exec("cmd")
type FuncCall struct {
	Name string
	Args []Node
}

func (f *FuncCall) nodeType() string { return "FuncCall" }

// OrExpr: expr or fallback
type OrExpr struct {
	Expr     Node
	Fallback Node
}

func (o *OrExpr) nodeType() string { return "OrExpr" }

// FuncDecl: fn name(params) -> returnType { body }
type FuncDecl struct {
	Name       string
	Params     []Param
	ReturnType string
	Body       []Node
}

func (f *FuncDecl) nodeType() string { return "FuncDecl" }

// Param: name: type
type Param struct {
	Name string
	Type string
}

// IfStmt: if cond { body } else { elseBody }
type IfStmt struct {
	Condition Node
	Body      []Node
	ElseBody  []Node
}

func (i *IfStmt) nodeType() string { return "IfStmt" }

// ForStmt: for item in collection { body }
type ForStmt struct {
	Var        string
	Collection Node
	Body       []Node
}

func (f *ForStmt) nodeType() string { return "ForStmt" }

// BinaryExpr: left op right
type BinaryExpr struct {
	Left  Node
	Op    string
	Right Node
}

func (b *BinaryExpr) nodeType() string { return "BinaryExpr" }

// UnaryExpr: !expr
type UnaryExpr struct {
	Op      string
	Operand Node
}

func (u *UnaryExpr) nodeType() string { return "UnaryExpr" }

// DotExpr: obj.field
type DotExpr struct {
	Object Node
	Field  string
}

func (d *DotExpr) nodeType() string { return "DotExpr" }

// ReturnStmt: return expr
type ReturnStmt struct {
	Value Node
}

func (r *ReturnStmt) nodeType() string { return "ReturnStmt" }

// ContinueStmt: continue
type ContinueStmt struct{}

func (c *ContinueStmt) nodeType() string { return "ContinueStmt" }

// ExitCall: exit(code)
type ExitCall struct {
	Code Node
}

func (e *ExitCall) nodeType() string { return "ExitCall" }

// ListLiteral: ["a", "b", "c"]
type ListLiteral struct {
	Elements []Node
}

func (l *ListLiteral) nodeType() string { return "ListLiteral" }

// MapLiteral: {key: value, ...}
type MapLiteral struct {
	Keys   []string
	Values []Node
}

func (m *MapLiteral) nodeType() string { return "MapLiteral" }

// BlockExpr: { stmts... lastExpr } — used in `or { ... }` blocks
type BlockExpr struct {
	Statements []Node
}

func (b *BlockExpr) nodeType() string { return "BlockExpr" }

// MatchStmt: match expr { cases }
type MatchStmt struct {
	Expr  Node
	Cases []MatchCase
}

func (m *MatchStmt) nodeType() string { return "MatchStmt" }

// MatchCase: pattern => body
type MatchCase struct {
	Pattern Node // nil means wildcard _
	Body    []Node
}

// WhileStmt: while condition { body }
type WhileStmt struct {
	Condition Node
	Body      []Node
}

func (w *WhileStmt) nodeType() string { return "WhileStmt" }

// BreakStmt: break
type BreakStmt struct{}

func (b *BreakStmt) nodeType() string { return "BreakStmt" }
