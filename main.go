package main

import (
	"fmt"
)

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
