package evaluator

import "novato/pkg/object"

func NewEnvironment() *Environment {
	s := make(map[string]object.Object)
	t := make(map[string]string)
	return &Environment{store: s, types: t}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewGlobalEnvironment() *Environment {
	env := NewEnvironment()

	for name, builtin := range builtins {
		env.Define(name, builtin)
	}

	return env
}

type Environment struct {
	store map[string]object.Object
	types map[string]string
	outer *Environment
}

func (e *Environment) Get(name string) (object.Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) GetType(name string) (string, bool) {
	typ, ok := e.types[name]
	if !ok && e.outer != nil {
		return e.outer.GetType(name)
	}
	return typ, ok
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
	delete(e.types, name)
	return val
}

func (e *Environment) DefineTyped(name string, val object.Object, typ string) object.Object {
	e.store[name] = val
	if typ != "" {
		e.types[name] = typ
	} else {
		delete(e.types, name)
	}
	return val
}
