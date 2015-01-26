package main

import (
	"fmt"
	"strconv"
)

// An Expression is a combination of values and operators that is reducible to
// another expression or to a value. Its purpose is to be evaluated to produce
// another expression.
type Expression interface {
	String() string
	Reducible() bool
	Reduce(environment map[string]Expression) Expression
	Equal(Expression) bool
}

// Boolean is our bool that is a fully reduced Expression.
type Boolean struct {
	Value bool
}

func (b Boolean) Reducible() bool {
	return false
}

func (b Boolean) Reduce(environment map[string]Expression) Expression {
	return b
}

func (b Boolean) String() string {
	return strconv.FormatBool(b.Value)
}

func (b Boolean) Equal(other Expression) bool {
	otherBoolean, ok := other.(Boolean)
	return ok && b.Value == otherBoolean.Value
}

// LessThan is an Expression with a Left and Right Expression that is reducible.
type LessThan struct {
	Left  Expression
	Right Expression
}

func (lt LessThan) String() string {
	return ""
}

func (lt LessThan) Equal(other Expression) bool {
	otherLt, ok := other.(LessThan)
	return ok && lt.Left == otherLt.Left && lt.Right == otherLt.Right
}

func (lt LessThan) Reducible() bool {
	return true
}

func (lt LessThan) Reduce(environment map[string]Expression) Expression {
	if lt.Left.Reducible() {
		return LessThan{lt.Left.Reduce(environment), lt.Right}
	} else if lt.Right.Reducible() {
		return LessThan{lt.Left, lt.Right.Reduce(environment)}
	} else {
		return Boolean{lt.Left.(Number).Value < lt.Right.(Number).Value}
	}
}

// Number is our int that is a fully reduced Expression.
type Number struct {
	Value int
}

func (n Number) Equal(other Expression) bool {
	otherNum, ok := other.(Number)
	return ok && n.Value == otherNum.Value
}

func (n Number) Reducible() bool {
	return false
}

func (n Number) Reduce(environment map[string]Expression) Expression {
	return n
}

func (n Number) String() string {
	return strconv.Itoa(n.Value)
}

// Add is an Expression with a Left and Right Expression that is reducible.
type Add struct {
	Left  Expression
	Right Expression
}

func (a Add) Equal(other Expression) bool {
	otherAdd, ok := other.(Add)
	return ok && a.Left == otherAdd.Left && a.Right == otherAdd.Right
}

func (a Add) Reducible() bool {
	return true
}

func (a Add) Reduce(environment map[string]Expression) Expression {
	if a.Left.Reducible() {
		return Add{a.Left.Reduce(environment), a.Right}
	} else if a.Right.Reducible() {
		return Add{a.Left, a.Right.Reduce(environment)}
	} else {
		return Number{a.Left.(Number).Value + a.Right.(Number).Value}
	}
}

func (a Add) String() string {
	return fmt.Sprintf("%s + %s", a.Left, a.Right)
}

// Multiply is an Expression with a Left and Right Expression that is reducible.
type Multiply struct {
	Left  Expression
	Right Expression
}

func (a Multiply) Equal(other Expression) bool {
	otherMultiply, ok := other.(Multiply)
	return ok && a.Left == otherMultiply.Left && a.Right == otherMultiply.Right
}

func (a Multiply) Reducible() bool {
	return true
}

func (m Multiply) Reduce(environment map[string]Expression) Expression {
	if m.Left.Reducible() {
		return Multiply{m.Left.Reduce(environment), m.Right}
	} else if m.Right.Reducible() {
		return Multiply{m.Left, m.Right.Reduce(environment)}
	} else {
		return Number{m.Left.(Number).Value * m.Right.(Number).Value}
	}
}

func (m Multiply) String() string {
	return fmt.Sprintf("%s * %s", m.Left, m.Right)
}

// Variable is an Expression with a Name that represents an Expression in the Environment.
type Variable struct {
	Name string
}

func (v Variable) Equal(other Expression) bool {
	otherVariable, ok := other.(Variable)
	return ok && v.Name == otherVariable.Name
}

func (v Variable) String() string {
	return v.Name
}

func (v Variable) Reducible() bool {
	return true
}

func (v Variable) Reduce(environment map[string]Expression) Expression {
	return environment[v.Name]
}
