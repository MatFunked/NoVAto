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
	"constante": "CONSTANTE",
	"func":      "FUNC",
	"clase":     "CLASE",
	"si":        "SI",
	"encaso":    "ENCASO",
	"sino":      "SINO",
	"retornar":  "RETORNAR",
	"importar":  "IMPORTAR",
	"de":        "DE",
	"romper":    "ROMPER",
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
