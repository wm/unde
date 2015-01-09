package main

import (
	"fmt"
	"strconv"
)

type Expression interface {
	String() string
	Reducible() bool
	Reduce(environment map[string]Expression) Expression
	Equal(Expression) bool
}

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

type Statement interface {
	String() string
	Reduce(map[string]Expression) (Statement, map[string]Expression)
	Reducible() bool
	Equal(Statement) bool
}

type DoNothing struct{}

func (dn DoNothing) String() string {
	return "do-nothing"
}

func (dn DoNothing) Reducible() bool {
	return false
}

func (dn DoNothing) Reduce(environment map[string]Expression) (Statement, map[string]Expression) {
	return dn, environment
}

func (dn DoNothing) Equal(other Statement) bool {
	_, ok := other.(DoNothing)
	return ok
}

type Assign struct {
	Name       string
	Expression Expression
}

func (a Assign) String() string {
	return fmt.Sprintf("%s = %s", a.Name, a.Expression)
}

func (a Assign) Reducible() bool {
	return true
}

func (a Assign) Equal(other Statement) bool {
	otherAssign, ok := other.(Assign)
	return ok && a.Name == otherAssign.Name && a.Expression == otherAssign.Expression
}

func (a Assign) Reduce(environment map[string]Expression) (Statement, map[string]Expression) {
	if a.Expression.Reducible() {
		return Assign{a.Name, a.Expression.Reduce(environment)}, environment
	} else {
		newEnv := copyEnvironment(environment)
		newEnv[a.Name] = a.Expression
		return DoNothing{}, newEnv
	}
}

func copyEnvironment(src map[string]Expression) map[string]Expression {
	cpy := make(map[string]Expression)

	for k, v := range src {
		cpy[k] = v
	}

	return cpy
}

type If struct {
	Condition   Expression
	Consequence Statement
	Alternative Statement
}

func (i If) String() string {
	return fmt.Sprintf("if %s { %s } else { %s }", i.Condition, i.Consequence, i.Alternative)
}

func (i If) Equal(other Statement) bool {
	otherIf, ok := other.(If)
	return ok && i.Condition == otherIf.Condition && i.Consequence == otherIf.Consequence && i.Alternative == otherIf.Alternative
}

func (i If) Reducible() bool {
	return true
}

func (i If) Reduce(environment map[string]Expression) (Statement, map[string]Expression) {
	switch {
	case i.Condition.Reducible():
		return If{i.Condition.Reduce(environment), i.Consequence, i.Alternative}, environment
	case i.Condition.Equal(Boolean{true}):
		return i.Consequence, environment
	case i.Condition.Equal(Boolean{false}):
		return i.Alternative, environment
	default:
		// Execution should never reach here but Golang complains.
		// I want to be explicit for true and false so I will keep both and
		// this.
		return i.Alternative, environment
	}
}

type Machine struct {
	Statement   Statement
	Environment map[string]Expression
}

func (m *Machine) Step() {
	m.Statement, m.Environment = m.Statement.Reduce(m.Environment)
}

func (m *Machine) Run() {
	for m.Statement.Reducible() {
		fmt.Println(m.Statement, m.Environment)
		m.Step()
	}
	fmt.Println(m.Statement, m.Environment)
}

func main() {
	machine := Machine{
		Assign{"x", Add{Variable{"x"}, Number{1}}},
		map[string]Expression{"x": Number{2}},
	}
	machine.Run()

	fmt.Println("\n")
	machine = Machine{
		If{Variable{"x"}, Assign{"y", Number{1}}, DoNothing{}},
		map[string]Expression{"x": Boolean{false}},
	}
	machine.Run()
}
