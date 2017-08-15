package parser

import (
	"fmt"
	"sol/ast"
	"sol/lexer"
)

type ParseError struct {
	err   string
	stack []string
}

func err(n *lexer.LexNode, err string, parsing string) *ParseError {
	e := &ParseError{
		err:   err,
		stack: []string{},
	}
	e.Trace(n, parsing)
	return e
}

func (e *ParseError) Error() string {
	str := e.err
	for _, err := range e.stack {
		str += "\n" + err
	}
	return str
}

func (e *ParseError) Trace(n *lexer.LexNode, parsing string) {
	e.stack = append(e.stack, fmt.Sprintf("at %d:%d parsing %s",
		n.LineNumber, n.LineIndex, parsing))
}

type ParserScope struct {
	parent   *ParserScope
	declared map[string]bool
}

func (s *ParserScope) Declare(name string) {
	s.declared[name] = true
}

func (s *ParserScope) IsDeclared(name string) bool {
	if s.declared[name] {
		return true
	}
	if s.parent != nil {
		return s.parent.IsDeclared(name)
	}
	return false
}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(input string) (*ast.Program, error) {

	s := lexer.NewScanner(input)

	// Basic syntax
	lSomeSpace := lexer.Regex("[\n\t\r ]+", false)
	lAnySpace := lexer.Regex("[\n\t\r ]*", true)

	lComma := lexer.Atom(",")
	lParenOpen := lexer.Atom("(")
	lParenClose := lexer.Atom(")")
	lBraceOpen := lexer.Atom("{")
	lBraceClose := lexer.Atom("}")

	// Literals
	lIdent := lexer.Group("identifier", lexer.Regex("[a-zA-Z]+", false))
	lInteger := lexer.Group("integer", lexer.Regex("[0-9]+", false))
	lString := lexer.Regex(`"((?:\\\\|\\"|[^"])+)"`, false)
	lFalse := lexer.Atom("false")
	lTrue := lexer.Atom("true")
	lNil := lexer.Atom("nil")

	// ???
	lParamSep := lexer.And(
		lAnySpace,
		lComma,
		lAnySpace,
	)
	lParamList := lexer.And(
		lParenOpen,
		lexer.Group("params", lexer.Interlace(
			lIdent,
			lParamSep,
		)),
		lParenClose,
	)

	// Keywords
	lKeyLet := lexer.Atom("let")
	lKeyReturn := lexer.Atom("return")
	lKeyFunc := lexer.Atom("fn")

	// Expressions and statements
	var lExpr lexer.Lexer
	var lStmt lexer.Lexer
	lExpr = lexer.Future(&lExpr, "lExpr")
	lStmt = lexer.Future(&lStmt, "lStmt")

	// Expressions
	lExprListSep := lexer.And(
		lAnySpace,
		lComma,
		lAnySpace,
	)

	// (<expr>, <expr>, <expr>, ...)
	lExprList := lexer.And(
		lParenOpen,
		lexer.Group("args", lexer.Interlace(
			lExpr,
			lExprListSep,
		)),
		lParenClose,
	)

	// ( <expr> )
	lExprClosed := lexer.Group("exprClosed", lexer.And(
		lParenOpen,
		lAnySpace,
		lExpr,
		lAnySpace,
		lParenClose,
	))

	// Primitive expressions
	lExprPrimitive := lexer.Or(
		lInteger,    // 5
		lString,     // "string"
		lFalse,      // false
		lTrue,       // true
		lNil,        // nil
		lIdent,      // x
		lExprClosed, // ( <expr> )
	)

	// Define operators for all precedence types
	lAssignmentOperators := lexer.And(
		lAnySpace,
		lexer.Group("operator", lexer.Atom("=")),
		lAnySpace,
	)
	lUnaryOperators := lexer.Or(
		lexer.Atom("!"),
		lexer.Atom("-"),
	)
	lMultiplyOperators := lexer.And(
		lAnySpace,
		lexer.Group("operator", lexer.Or(
			lexer.Atom("/"),
			lexer.Atom("*"),
		)),
		lAnySpace,
	)
	lAddOperators := lexer.And(
		lAnySpace,
		lexer.Group("operator", lexer.Or(
			lexer.Atom("+"),
			lexer.Atom("-"),
		)),
		lAnySpace,
	)
	lCompareOperators := lexer.And(
		lAnySpace,
		lexer.Group("operator", lexer.Or(
			lexer.Atom(">="),
			lexer.Atom("<="),
			lexer.Atom(">"),
			lexer.Atom("<"),
		)),
		lAnySpace,
	)
	lEqualityOperators := lexer.And(
		lAnySpace,
		lexer.Group("operator", lexer.Or(
			lexer.Atom("!="),
			lexer.Atom("=="),
		)),
		lAnySpace,
	)

	// Declare all operators
	var lAssignment, lUnary, lMultiply, lAdd, lCompare, lEquality lexer.Lexer
	lAssignment = lexer.Future(&lAssignment, "assignment")
	lUnary = lexer.Future(&lUnary, "unary")
	lMultiply = lexer.Future(&lMultiply, "multiply")
	lAdd = lexer.Future(&lAdd, "add")
	lCompare = lexer.Future(&lCompare, "compare")
	lEquality = lexer.Future(&lEquality, "equality")

	// Define all operators
	lAssignment = lexer.And(
		lexer.Group("left", lExprPrimitive),
		lexer.Group("rest", lexer.Repeat(
			lexer.And(
				lAssignmentOperators,
				lexer.Group("right", lExprPrimitive),
			), 0, 1,
		)),
	)
	lUnary = lexer.Group("right", lexer.Or(
		lexer.And(
			lUnaryOperators,
			lUnary,
		),
		lAssignment,
	))
	lMultiply = lexer.And(
		lexer.Group("left", lUnary),
		lexer.Group("rest", lexer.Repeat(
			lexer.And(
				lMultiplyOperators,
				lexer.Group("right", lUnary),
			), 0, 1,
		)),
	)
	lAdd = lexer.And(
		lexer.Group("left", lMultiply),
		lexer.Group("rest", lexer.Repeat( // TODO: Net zoals hierboven overal de 'right' toevoegen
			lexer.And(
				lAddOperators,
				lexer.Group("right", lMultiply),
			), 0, 1,
		)),
	)
	lCompare = lexer.And(
		lexer.Group("left", lAdd),
		lexer.Group("rest", lexer.Repeat(
			lexer.And(
				lCompareOperators,
				lexer.Group("right", lAdd),
			), 0, 1,
		)),
	)
	lEquality = lexer.And(
		lexer.Group("left", lCompare),
		lexer.Group("rest", lexer.Repeat(
			lexer.And(
				lEqualityOperators,
				lexer.Group("right", lCompare),
			), 0, 1,
		)),
	)

	lExprFunc := lexer.Group("exprFunc", lexer.And(
		lKeyFunc,   // fn
		lAnySpace,  //
		lParamList, // (a, b, c)
		lAnySpace,  //
		lBraceOpen, // {
		lAnySpace,  //
		lexer.Group("body", lexer.Repeat(lStmt, 0, -1)), // ...
		lAnySpace,   //
		lBraceClose, // }
	))
	lExprCall := lexer.Group("exprCall", lexer.And(
		lIdent, // functionName
		lAnySpace,
		lExprList,
	))
	lExpr = lexer.Group("expression", lexer.Or(
		lExprFunc, // fn(a, b) { <stmts> }
		lExprCall, // fn(a, b)
		lEquality, // strings, bools, ints
	))

	// Statements
	lStmtDeclare := lexer.Group("stmtDeclare", lexer.And(
		lKeyLet,         // let
		lSomeSpace,      //
		lIdent,          // <identifier>
		lAnySpace,       //
		lexer.Atom("="), // =
		lAnySpace,       //
		lExpr,           // <expression>
	))
	lStmtReturn := lexer.Group("stmtReturn", lexer.And(
		lKeyReturn, // return
		lSomeSpace, //
		lExpr,      // <expression>
	))
	lStmtBlock := lexer.Group("stmtBlock", lexer.And(
		lBraceOpen,
		lAnySpace,
		lexer.Repeat(
			lexer.And(
				lStmt,
				lAnySpace,
			), 0, -1,
		),
		lBraceClose,
	))
	lStmtExpr := lexer.Group("stmtExpr", lexer.And(
		lExpr, // <expression>
	))
	lStmt = lexer.Group("statement", lexer.Or(
		lStmtBlock,
		lStmtDeclare,
		lStmtReturn,
		lStmtExpr,
	))

	// Program
	lProgram := lexer.And(
		lAnySpace,
		lexer.Interlace(lStmt, lSomeSpace),
		lAnySpace,
	)

	// Actually lex the input
	tree, lexErr := lProgram.Lex(s)
	if lexErr != nil {
		return nil, lexErr
	}
	// fmt.Println(tree.String(0))

	// Parse the tree
	prog, parseErr := p.parseProgram(tree)
	if parseErr != nil {
		return nil, error(parseErr)
	}
	return prog, nil
}

func (p *Parser) parseProgram(node *lexer.LexNode) (*ast.Program, *ParseError) {

	stmtNodes := node.GroupNodes("statement")
	stmts := []ast.Statement{}

	for _, stmtNode := range stmtNodes {
		stmt, err := p.parseStatement(stmtNode)
		if err != nil {
			err.Trace(node, "program")
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	return &ast.Program{
		Statements: stmts,
	}, nil
}

func (p *Parser) parseIdentifier(node *lexer.LexNode) (string, *ParseError) {
	return node.Value, nil
}

func (p *Parser) parseStatement(node *lexer.LexNode) (ast.Statement, *ParseError) {

	keyword := node.Children[0]
	var stmt ast.Statement
	var err *ParseError

	switch keyword.GroupName {
	case "stmtDeclare":
		stmt, err = p.parseDeclarationStatement(node)
	case "stmtReturn":
		stmt, err = p.parseReturnStatement(node)
	case "stmtBlock":
		stmt, err = p.parseBlockStatement(node)
	default:
		stmt, err = p.parseExpressionStatement(node)
	}

	if err != nil {
		err.Trace(node, "statement")
	}

	return stmt, err
}

func (p *Parser) parseDeclarationStatement(node *lexer.LexNode) (ast.Statement, *ParseError) {

	// let <identifier> = <expression>
	nodeIdent := node.GroupNode("identifier")
	nodeExpr := node.GroupNode("expression")

	// Parse identifier
	ident, err := p.parseIdentifier(nodeIdent)
	if err != nil {
		err.Trace(node, "declaration statement")
		return nil, err
	}

	// Parse expression
	expr, err := p.parseExpression(nodeExpr)
	if err != nil {
		err.Trace(node, "declaration statement")
		return nil, err
	}

	// Return a declaration statement
	return &ast.DeclarationStatement{
		Identifier: ident,
		Expression: expr,
	}, nil
}

func (p *Parser) parseReturnStatement(node *lexer.LexNode) (ast.Statement, *ParseError) {
	expr, err := p.parseExpression(node.GroupNode("expression"))
	if err != nil {
		return nil, err
	}

	return &ast.ReturnStatement{
		Expression: expr,
	}, nil
}

func (p *Parser) parseExpressionStatement(node *lexer.LexNode) (ast.Statement, *ParseError) {
	expr, err := p.parseExpression(node.GroupNode("expression"))
	if err != nil {
		return nil, err
	}

	return &ast.ExpressionStatement{
		Expression: expr,
	}, nil
}

func (p *Parser) parseBlockStatement(node *lexer.LexNode) (ast.Statement, *ParseError) {
	var stmts []ast.Statement
	for _, stmtNode := range node.GroupNodes("statement") {
		stmt, err := p.parseStatement(stmtNode)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	return &ast.BlockStatement{
		Statements: stmts,
	}, nil
}

func (p *Parser) parseExpression(node *lexer.LexNode) (ast.Expression, *ParseError) {

	// There's a few expression types we can encounter, including:
	// - Function definitions: fn(a,b) { <stmts> }
	// - Function calls: ident(a, b)
	// - Infix expressions: (a + b) - c, x - 4 * y
	// - Primitive expressions: c, 5

	// Handle function definitions and calls
	if node.GroupName == "exprFunc" {
		return p.parseFunctionExpression(node)
	} else if node.GroupName == "exprCall" {
		return p.parseCallExpression(node)
	}

	// Next, we want to move down to the primitive
	// parts of this expression.

	// Example expression, with two pointers indicating
	// the positions and length of 'left' and 'rest':
	//		5 + 5
	//      ^ ^^^

	var leftNode, restNode *lexer.LexNode
	leftNode = node
	restNode = nil

	// We need to go deeper!
	for (restNode == nil || restNode.Value == "") &&
		leftNode.GroupExists("left") &&
		leftNode.GroupExists("rest") {

		current := leftNode
		leftNode = current.GroupNode("left")
		restNode = current.GroupNode("rest")
	}

	// If, at this point, the rest node is blank,
	// then we must be dealing with a primitive.
	if restNode == nil || restNode.Value == "" {
		return p.parsePrimitiveExpression(leftNode)
	}

	// In all other cases, we are presumably dealing
	// with an infix expression.

	// Get the operator and right hand side expression if necessary
	operatorNode := restNode.GroupNode("operator")
	rightNode := restNode.GroupNode("right")

	var left ast.Expression
	var err *ParseError

	left, err = p.parseExpression(leftNode)
	if err != nil {
		err.Trace(node, "expression")
		return nil, err
	}

	right, err := p.parseExpression(rightNode)
	if err != nil {
		err.Trace(node, "expression")
		return nil, err
	}

	return &ast.InfixExpression{
		Left:     left,
		Operator: operatorNode.Value,
		Right:    right,
	}, nil
}

func (p *Parser) parsePrimitiveExpression(node *lexer.LexNode) (ast.Expression, *ParseError) {
	if node == nil {
		return nil, nil
	}

	// Integers
	if node.GroupNode("integer") != nil {
		return &ast.IntegerExpression{
			Literal: node.Value,
		}, nil
	}

	// Strings
	if node.GroupNode("string") != nil {
		// TODO: String parsing
	}

	// Identifiers
	if node.GroupNode("identifier") != nil {
		return &ast.IdentifierExpression{
			Literal: node.Value,
		}, nil
	}

	return nil, err(
		node,
		fmt.Sprintf("unknown primitive expression type: %s", node.String(1)),
		"primitive expression",
	)
}

func (p *Parser) parseIdentifierExpression(node *lexer.LexNode) (ast.Expression, *ParseError) {
	return &ast.IdentifierExpression{
		Literal: node.Value,
	}, nil
}

func (p *Parser) parseIntegerExpression(node *lexer.LexNode) (ast.Expression, *ParseError) {
	return &ast.IntegerExpression{
		Literal: node.Value,
	}, nil
}

func (p *Parser) parseInfixExpression(node *lexer.LexNode) (ast.Expression, *ParseError) {
	leftNode := node.GroupNode("left")
	rightNode := node.GroupNode("right")
	operatorNode := node.GroupNode("operator")

	left, err := p.parseExpression(leftNode)
	if err != nil {
		err.Trace(node, "infix")
		return nil, err
	}

	right, err := p.parseExpression(rightNode)
	if err != nil {
		err.Trace(node, "infix")
		return nil, err
	}

	return &ast.InfixExpression{
		Left:     left,
		Right:    right,
		Operator: operatorNode.Value,
	}, nil
}

func (p *Parser) parseClosedExpression(node *lexer.LexNode) (ast.Expression, *ParseError) {
	return p.parseExpression(node.GroupNode("expression"))
}

func (p *Parser) parseFunctionExpression(node *lexer.LexNode) (ast.Expression, *ParseError) {

	var params []string

	paramsNode := node.GroupNode("params")
	for _, paramNode := range paramsNode.GroupNodes("identifier") {
		params = append(params, paramNode.Value)
	}

	block, err := p.parseBlockStatement(node.GroupNode("stmtBlock"))
	if err != nil {
		err.Trace(node, "function")
		return nil, err
	}

	return &ast.FunctionExpression{
		Parameters: params,
		Body:       block,
	}, nil
}

func (p *Parser) parseCallExpression(node *lexer.LexNode) (ast.Expression, *ParseError) {

	identifier := node.GroupNode("identifier").Value

	argsNode := node.GroupNode("args")
	argsNodes := argsNode.GroupNodes("expression")
	var args []ast.Expression

	for _, argNode := range argsNodes {
		expr, err := p.parseExpression(argNode)
		if err != nil {
			err.Trace(node, "call")
			return nil, err
		}
		args = append(args, expr)
	}

	return &ast.CallExpression{
		Identifier: identifier,
		Arguments:  args,
	}, nil
}
