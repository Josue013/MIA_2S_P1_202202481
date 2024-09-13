package filesystem

import (
	"fmt"
	"os"
	
)

func Cat(fileValores []string) string {
	var respuesta string
	for _, filePath := range fileValores {
			// Abrir el archivo
			file, err := os.Open(filePath)
			if err != nil {
					fmt.Println("Error: No se pudo abrir el archivo", filePath)
					respuesta += "Error: No se pudo abrir el archivo " + filePath + "\n"
					continue
			}
			defer file.Close()

			// Leer el contenido del archivo
			fileInfo, err := file.Stat()
			if err != nil {
					fmt.Println("Error: No se pudo obtener la información del archivo", filePath)
					respuesta += "Error: No se pudo obtener la información del archivo " + filePath + "\n"
					continue
			}

			fileSize := fileInfo.Size()
			buffer := make([]byte, fileSize)
			_, err = file.Read(buffer)
			if err != nil {
					fmt.Println("Error: No se pudo leer el archivo", filePath)
					respuesta += "Error: No se pudo leer el archivo " + filePath + "\n"
					continue
			}

			// Mostrar el contenido del archivo
			fmt.Println("Contenido del archivo", filePath, ":\n", string(buffer))
			respuesta += "Contenido del archivo " + filePath + ":\n" + string(buffer) + "\n"
	}
	return respuesta
}