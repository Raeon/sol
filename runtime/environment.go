package runtime

import (
	"fmt"
	"sol/ast"
	"strconv"
)

type Environment struct {
	scope *Scope
}

var ops map[string]func(Object, Object) Object

func NewEnv() *Environment {
	if ops == nil {
		ops = make(map[string]func(Object, Object) Object)
		ops["+"] = applyAdd
		ops["-"] = applySubtract
		ops["*"] = applyMultiply
		ops["/"] = applyDivide
	}
	return &Environment{
		scope: NewScope(),
	}
}

func (e *Environment) Evaluate(node ast.Node) Object {
	switch node.(type) {

	case *ast.Program:
		prog, _ := node.(*ast.Program)
		var last Object
		for _, stmt := range prog.Statements {
			last = e.Evaluate(stmt)
		}
		return last

	case *ast.BlockStatement:
		blockStmt, _ := node.(*ast.BlockStatement)
		var value Object
		value = &Nil{}
		for _, stmt := range blockStmt.Statements {
			value = e.Evaluate(stmt)

			// Check for return value
			if retVal, ok := value.(*ReturnValue); ok {
				value = retVal.Value
				break
			}
		}
		return value

	case *ast.ReturnStatement:
		retStmt, _ := node.(*ast.ReturnStatement)
		value := e.Evaluate(retStmt.Expression)
		return &ReturnValue{Value: value}

	case *ast.DeclarationStatement:
		decStmt, _ := node.(*ast.DeclarationStatement)
		value := e.Evaluate(decStmt.Expression)
		e.scope.Set(decStmt.Identifier, value)
		return value

	case *ast.ExpressionStatement:
		exprStmt, _ := node.(*ast.ExpressionStatement)
		return e.Evaluate(exprStmt.Expression)

	case *ast.IntegerExpression:
		num, err := strconv.Atoi(node.ToString())
		if err != nil {
			panic(err)
		}
		return &Number{Value: num}

	case *ast.IdentifierExpression:
		return e.scope.Get(node.ToString())

	case *ast.ClosedExpression:
		fmt.Println("closedExpr")
		expr, _ := node.(*ast.ClosedExpression)
		return e.Evaluate(expr.Expression)

	case *ast.InfixExpression:
		expr, _ := node.(*ast.InfixExpression)
		return e.applyOperator(expr.Operator, expr.Left, expr.Right)

	}
	panic(fmt.Sprintf("Uninterpreted AST node encountered: %s", node.ToString()))
}

func (e *Environment) applyOperator(op string, left, right ast.Node) Object {

	if op == "=" {
		return e.applyAssign(left, right)
	}

	fn, ok := ops[op]
	if ok {
		return fn(e.Evaluate(left), e.Evaluate(right))
	}
	return &Nil{}
}

func (e *Environment) applyAssign(left, right ast.Node) Object {
	ident, ok := left.(*ast.IdentifierExpression)
	if !ok {
		return &Exception{
			Message: "Can only assign to identifiers",
		}
	}

	return e.scope.Set(ident.Literal, e.Evaluate(right))
}

func applyAdd(left, right Object) Object {
	lt := left.TypeString()
	rt := right.TypeString()

	if lt == "number" && lt == rt {
		leftNum, _ := left.(*Number)
		rightNum, _ := right.(*Number)
		return &Number{Value: leftNum.Value + rightNum.Value}
	}

	return &Exception{Message: "Cannot add non-numbers"}
}

func applySubtract(left, right Object) Object {
	lt := left.TypeString()
	rt := right.TypeString()

	if lt == "number" && lt == rt {
		leftNum, _ := left.(*Number)
		rightNum, _ := right.(*Number)
		return &Number{Value: leftNum.Value - rightNum.Value}
	}

	return &Exception{Message: "Cannot subtract non-numbers"}
}

func applyMultiply(left, right Object) Object {
	lt := left.TypeString()
	rt := right.TypeString()

	if lt == "number" && lt == rt {
		leftNum, _ := left.(*Number)
		rightNum, _ := right.(*Number)
		return &Number{Value: leftNum.Value * rightNum.Value}
	}

	return &Exception{Message: "Cannot multiply non-numbers"}
}

func applyDivide(left, right Object) Object {
	lt := left.TypeString()
	rt := right.TypeString()

	if lt == "number" && lt == rt {
		leftNum, _ := left.(*Number)
		rightNum, _ := right.(*Number)
		return &Number{Value: leftNum.Value / rightNum.Value}
	}

	return &Exception{Message: "Cannot divide non-numbers"}
}
