package evaluator

import "novato/pkg/object"

func NewEnvironment() *Environment {
	s := make(map[string]object.Object)
	return &Environment{store: s}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

type Environment struct {
	store map[string]object.Object
	outer *Environment
}

func (e *Environment) Get(name string) (object.Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val object.Object) object.Object {
	if _, ok := e.store[name]; !ok && e.outer != nil {
		if _, found := e.outer.Get(name); found {
			return e.outer.Set(name, val)
		}
	}
	e.store[name] = val
	return val
}

func (e *Environment) Define(name string, val object.Object) object.Object {
	e.store[name] = val
	return val
}
