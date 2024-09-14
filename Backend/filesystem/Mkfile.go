package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	//"strings"
	"io/ioutil"
)

// Estructura para el comando MKFILE
type MKFILE struct {
	Path string
	Size int
	R    bool
	Cont string
}

// Función para crear un archivo
func commandMkfile(mkfile *MKFILE) error {


	// Validar que la ruta padre existe, o crearla si el parámetro -r está presente
	parentDir := filepath.Dir(mkfile.Path)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		if mkfile.R {
			// Crear las carpetas padres
			err := os.MkdirAll(parentDir, os.ModePerm)
			if err != nil {
				return fmt.Errorf("Error al crear directorios padres: %w", err)
			}
		} else {
			return fmt.Errorf("Error: las carpetas padres no existen y el parámetro -r no fue proporcionado")
		}
	}

	// Si el parámetro -cont está presente, leer el contenido del archivo especificado
	var content string
	if mkfile.Cont != "" {
		fileContent, err := ioutil.ReadFile(mkfile.Cont)
		if err != nil {
			return fmt.Errorf("Error al leer el archivo de contenido: %w", err)
		}
		content = string(fileContent)
	} else {
		// Generar contenido si no se proporcionó -cont
		content = generateContent(mkfile.Size)
	}

	// Crear o sobrescribir el archivo
	err := ioutil.WriteFile(mkfile.Path, []byte(content), 0664) // Permisos 664
	if err != nil {
		return fmt.Errorf("Error al crear el archivo: %w", err)
	}

	fmt.Println("Archivo creado con éxito:", mkfile.Path)
	return nil
}

// Función para generar contenido en base al tamaño
func generateContent(size int) string {
	if size <= 0 {
		return ""
	}
	content := ""
	for len(content) < size {
		content += "0123456789"
	}
	return content[:size] // Limitar el contenido al tamaño especificado
}
