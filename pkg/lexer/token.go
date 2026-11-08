package lexer

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// Especiales
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	// Identificadores y Literales
	IDENT  = "IDENT"  // variables, funciones (ej: verificar_acceso, usuario)
	INT    = "INT"    // números enteros (ej: 3, 18)
	FLOAT  = "FLOAT"  // números decimales (ej: 3.14, 0.5)
	STRING = "STRING" // cadenas de texto (ej: "Acceso concedido")

	// Operadores y símbolos
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	RT     = ">"
	LTE    = "<="
	GTE    = ">="
	ARROW  = "->"

	COMMA    = ","
	COLON    = ":"
	DOT      = "."
	LPAREN   = "("
	RPAREN   = ")"
	LBRACKET = "["
	RBRACKET = "]"
	LBRACE   = "{" // <-- Agrega esto
	RBRACE   = "}" // <-- Agrega esto
)

// Palabras reservadas según tu diseño de NoVAto
var keywords = map[string]TokenType{
	"variable":  "VARIABLE",
	"var":       "VARIABLE",
	"constante": "CONSTANTE",
	"func":      "FUNC",
	"funcion":   "FUNC",
	"clase":     "CLASE",

	//if, elif y else
	"si":   "SI",
	"osi":  "OSI",
	"sino": "SINO",

	//match, case
	"comparar": "COMPARAR",
	"encaso":   "ENCASO",
	"defecto":  "DEFECTO",

	// bucles
	"mientras": "MIENTRAS",
	"para":     "PARA",
	"hasta":    "HASTA",

	//return
	"retornar": "RETORNAR",
	//yield
	"ceder": "CEDER",

	//from, import
	"importar": "IMPORTAR",
	"de":       "DE",

	//raise
	"romper": "ROMPER",
	//control de flujo
	"continuar": "CONTINUAR",

	// Tipos de datos
	"ent":  "T_ENT",
	"text": "T_TEXT",
	"bool": "T_BOOL",
	"deci": "T_DECI",
	// Booleanos literales
	"Verdadero": "VERDADERO",
	"Falso":     "FALSO",
}

// Verifica si un identificador es una palabra reservada
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
