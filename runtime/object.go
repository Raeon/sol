package runtime

import (
	"fmt"
	"sol/ast"
)

type Object interface {
	ToString() string
	IsEqual(other Object) bool
	TypeString() string
}

type Nil struct{}

func (n *Nil) ToString() string {
	return "nil"
}

func (n *Nil) IsEqual(other Object) bool {
	_, ok := other.(*Nil)
	return ok
}

func (n *Nil) TypeString() string {
	return n.ToString()
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) ToString() string {
	return r.Value.ToString()
}

func (r *ReturnValue) IsEqual(other Object) bool {
	return r.Value.IsEqual(other)
}

func (r *ReturnValue) TypeString() string {
	return r.Value.TypeString()
}

type Boolean struct {
	Value bool
}

func (b *Boolean) ToString() string {
	return fmt.Sprintf("%v", b.Value)
}

func (b *Boolean) IsEqual(other Object) bool {
	bl, ok := other.(*Boolean)
	if ok {
		return b.Value == bl.Value
	}
	return false
}

func (b *Boolean) TypeString() string {
	return "boolean"
}

type Number struct {
	Value int
}

func (n *Number) ToString() string {
	return fmt.Sprintf("%d", n.Value)
}

func (n *Number) IsEqual(other Object) bool {
	num, ok := other.(*Number)
	return ok && n.Value == num.Value
}

func (n *Number) TypeString() string {
	return "number"
}

type Exception struct {
	Message string
}

func (e *Exception) ToString() string {
	return fmt.Sprintf("%s", e.Message)
}

func (e *Exception) IsEqual(other Object) bool {
	ex, ok := other.(*Exception)
	return ok && e.Message == ex.Message
}

func (e *Exception) TypeString() string {
	return "exception"
}

type Function struct {
	Parameters []string
	Body       *ast.BlockStatement
}

func (f *Function) ToString() string {
	str := "fn("
	for i, param := range f.Parameters {
		if i != 0 {
			str += ", "
		}
		str += param
	}
	str += ") {"
	str += f.Body.ToString()
	return str + "}"
}

func (f *Function) IsEqual(other Object) bool {
	fn, ok := other.(*Function)
	if ok {
		return fn == f
	}
	return false
}

func (f *Function) TypeString() string {
	return "function"
}
