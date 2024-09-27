package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)


func getContent(files []string) {
	var functions []string

	for _, path := range files {
		file, err := os.Open(path)
		if err != nil {
			log.Fatalf("Error opening the file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var contentBuilder strings.Builder
 		inDocSection := false
    	functionOpenCount := 0

    	for scanner.Scan() {
        	line := scanner.Text()

	        // Comprobar si la línea contiene el patrón
	        if strings.Contains(line, "@api_generate_doc") {
	            inDocSection = true
	            contentBuilder.WriteString(line + "\n")
	            functionOpenCount++ // Contar la apertura de función
	            continue
	        }

	        // Si estamos dentro de la sección de documentación, almacenamos la línea
	        if inDocSection {
	            contentBuilder.WriteString(line + "\n")
	            // Comprobar si hemos llegado al final de la función
	            functionOpenCount += strings.Count(line, "{") - strings.Count(line, "}")
	            if functionOpenCount == 0 {
	                // Fin de la sección de documentación
	                break
	            }
	        }
	    }

	    if err := scanner.Err(); err != nil {
	        log.Fatalf("Error leyendo el archivo: %v", err)
	    }

	    // Mostrar y guardar el contenido extraído
	    extractedContent := contentBuilder.String()
	    if extractedContent != "" {
	        fmt.Println("Contenido extraído")

			functions = append(functions, extractedContent)
		}
	}
	fmt.Println(functions)
}