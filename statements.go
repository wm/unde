package main

import "fmt"

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
