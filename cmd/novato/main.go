package main

import (
	"fmt"
	"novato/pkg/evaluator"
	"novato/pkg/lexer"
	"novato/pkg/parser"
	"os"
)

func main() {
	// 1. Validar que se haya pasado un archivo como argumento
	if len(os.Args) < 2 {
		fmt.Println("Error: Debes proporcionar un archivo .va")
		fmt.Println("Uso: novato <archivo.va>")
		os.Exit(1)
	}

	archivoPath := os.Args[1]

	// 2. Leer el archivo fuente
	contenidoBytes, err := os.ReadFile(archivoPath)
	if err != nil {
		fmt.Printf("Error: No se pudo leer el archivo '%s'. Detalles: %v\n", archivoPath, err)
		os.Exit(1)
	}

	codigoFuente := string(contenidoBytes)

	fmt.Printf("--- Compilando / Interpretando: %s ---\n", archivoPath)

	// 3. Procesar el código con nuestro Lexer
	procesarCodigo(codigoFuente)
}

func procesarCodigo(codigo string) {
	l := lexer.New(codigo)
	p := parser.New(l)
	programa := p.ParseProgram()

	env := evaluator.NewEnvironment()
	resultado := evaluator.Eval(programa, env)
	if resultado != nil {
		fmt.Printf("Resultado de ejecución: %s\n", resultado.Inspect())
	}

	fmt.Printf("--- Árbol de Sintaxis Abstracta (AST) ---\n")
	for _, stmt := range programa.Statements {
		if stmt.String() != "" {
			fmt.Printf("Sentencia detectada: %s\n", stmt.String())
		}
	}
}
