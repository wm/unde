package main

import (
	"fmt"
	"strconv"
)

type Expression interface {
	String() string
	Reducible() bool
	Reduce(environment map[string]Expression) Expression
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

type LessThan struct {
	Left  Expression
	Right Expression
}

func (lt LessThan) String() string {
	return ""
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

func (v Variable) String() string {
	return v.Name
}

func (v Variable) Reducible() bool {
	return true
}

func (v Variable) Reduce(environment map[string]Expression) Expression {
	return environment[v.Name]
}

type Machine struct {
	Expression  Expression
	Environment map[string]Expression
}

func (m *Machine) Step() {
	m.Expression = m.Expression.Reduce(m.Environment)
}

func (m *Machine) Run() {
	for m.Expression.Reducible() {
		fmt.Println(m.Expression)
		m.Step()
	}
	fmt.Println(m.Expression)
}

func main() {
	environment := map[string]Expression{}

	machine := Machine{
		Add{
			Multiply{Number{1}, Number{2}},
			Multiply{Number{3}, Number{4}},
		},
		environment,
	}
	machine.Run()

	machine = Machine{
		LessThan{
			Number{5},
			Add{
				Number{2},
				Number{2},
			},
		},
		environment,
	}
	machine.Run()

	environment = map[string]Expression{
		"x": Number{3},
		"y": Number{4},
	}

	machine = Machine{
		Add{
			Variable{"x"},
			Variable{"y"},
		},
		environment,
	}
	machine.Run()
}
