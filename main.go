package main

import (
	"fmt"
	"strconv"
)

type Expression interface {
	String() string
	Reducible() bool
	Reduce() Expression
}

type Boolean struct {
	Value bool
}

func (b Boolean) Reducible() bool {
	return false
}

func (b Boolean) Reduce() Expression {
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

func (lt LessThan) Reduce() Expression {
	if lt.Left.Reducible() {
		return LessThan{lt.Left.Reduce(), lt.Right}
	} else if lt.Right.Reducible() {
		return LessThan{lt.Left, lt.Right.Reduce()}
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

func (n Number) Reduce() Expression {
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

func (a Add) Reduce() Expression {
	if a.Left.Reducible() {
		return Add{a.Left.Reduce(), a.Right}
	} else if a.Right.Reducible() {
		return Add{a.Left, a.Right.Reduce()}
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

func (m Multiply) Reduce() Expression {
	if m.Left.Reducible() {
		return Multiply{m.Left.Reduce(), m.Right}
	} else if m.Right.Reducible() {
		return Multiply{m.Left, m.Right.Reduce()}
	} else {
		return Number{m.Left.(Number).Value * m.Right.(Number).Value}
	}
}

func (m Multiply) String() string {
	return fmt.Sprintf("%s * %s", m.Left, m.Right)
}

type Machine struct {
	Expression Expression
}

func (m *Machine) Step() {
	m.Expression = m.Expression.Reduce()
}

func (m *Machine) Run() {
	for m.Expression.Reducible() {
		fmt.Println(m.Expression)
		m.Step()
	}
	fmt.Println(m.Expression)
}

func main() {
	machine := Machine{
		Add{
			Multiply{Number{1}, Number{2}},
			Multiply{Number{3}, Number{4}},
		},
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
	}
	machine.Run()
}