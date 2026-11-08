package object

import (
	"fmt"
	"math"
	"novato/pkg/ast"
	"strconv"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ  = "INTEGER"
	DECIMAL_OBJ  = "DECI"
	STRING_OBJ   = "STRING"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	FUNCTION_OBJ = "FUNCTION"
	BUILTIN_OBJ  = "BUILTIN"
	ERROR_OBJ    = "ERROR"

	BREAK_SIGNAL_OBJ    = "BREAK_SIGNAL"
	CONTINUE_SIGNAL_OBJ = "CONTINUE_SIGNAL"
	YIELD_SIGNAL_OBJ    = "YIELD_SIGNAL"
	RETURN_SIGNAL_OBJ   = "RETURN_SIGNAL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Entorno léxico mínimo que necesita una función para soportar closures.
type Environment interface {
	Get(name string) (Object, bool)
	Set(name string, val Object) Object
	Define(name string, val Object) Object
}

// Objeto Entero
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// Objeto Decimal
type Decimal struct {
	Value float64
}

func (d *Decimal) Type() ObjectType { return DECIMAL_OBJ }

func (d *Decimal) Inspect() string {
	s := strconv.FormatFloat(d.Value, 'f', -1, 64)
	if math.IsInf(d.Value, 0) || math.IsNaN(d.Value) {
		return s
	}
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	return s
}

// Objeto Texto (cadena)
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Objeto Booleano
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// Objeto Null (para cuando una expresión no retorna nada)
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "nulo" }

// Señal de control: romper (break)
type BreakSignal struct{}

func (b *BreakSignal) Type() ObjectType { return BREAK_SIGNAL_OBJ }
func (b *BreakSignal) Inspect() string  { return "romper" }

// Señal de control: continuar (continue)
type ContinueSignal struct{}

func (c *ContinueSignal) Type() ObjectType { return CONTINUE_SIGNAL_OBJ }
func (c *ContinueSignal) Inspect() string  { return "continuar" }

// Señal de control: ceder (yield), con valor opcional
type YieldSignal struct {
	Value Object
}

func (y *YieldSignal) Type() ObjectType { return YIELD_SIGNAL_OBJ }
func (y *YieldSignal) Inspect() string {
	if y.Value == nil {
		return "ceder"
	}
	return "ceder " + y.Value.Inspect()
}

// Señal de control: retornar (return), con valor opcional
type ReturnSignal struct {
	Value Object
}

func (r *ReturnSignal) Type() ObjectType { return RETURN_SIGNAL_OBJ }
func (r *ReturnSignal) Inspect() string {
	if r.Value == nil {
		return "retornar"
	}
	return "retornar " + r.Value.Inspect()
}

// Objeto Función (incluye el entorno léxico donde fue definida para closures)
type Function struct {
	Parameters     []*ast.Identifier
	ParameterTypes []string
	ReturnType     string
	Body           *ast.BlockStatement
	Env            Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	params := make([]string, 0, len(f.Parameters))
	for i, param := range f.Parameters {
		paramStr := param.Value
		if i < len(f.ParameterTypes) && f.ParameterTypes[i] != "" {
			paramStr += " " + f.ParameterTypes[i]
		}
		params = append(params, paramStr)
	}

	signature := "func(" + strings.Join(params, ", ") + ")"
	if f.ReturnType != "" {
		signature += " " + f.ReturnType
	}

	return signature + " { ... }"
}

// Objeto Función Nativa (envuelve una función Go)
type Builtin struct {
	Fn func(args ...Object) Object
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "función nativa" }

// Objeto Error (para errores en tiempo de ejecución)
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
