package main

import (
	"fmt"
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

	// Aquí llamaremos próximamente a nuestro Lexer y Parser
	procesarCodigo(codigoFuente)
}

func procesarCodigo(codigo string) {
	// De momento, solo confirmamos que el texto llegó al motor
	fmt.Println("Código cargado con éxito en memoria. Longitud:", len(codigo), "caracteres.")
}
