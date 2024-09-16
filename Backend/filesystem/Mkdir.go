package filesystem

import (
    "fmt"
    "os"
    "path/filepath"
    
)

// Estructura para el comando MKDIR
type MKDIR struct {
    Path string
    P    bool
}

// Función para crear un directorio
func commandMkdir(mkdir *MKDIR) error {
    // Verificar que haya una sesión activa
    if Usr_sesion.Uid == -1 {
        return fmt.Errorf("Error: Necesita iniciar sesión")
    }

    // Validar que la ruta padre existe, o crearla si el parámetro -p está presente
    parentDir := filepath.Dir(mkdir.Path)
    if _, err := os.Stat(parentDir); os.IsNotExist(err) {
        if mkdir.P {
            // Crear las carpetas padres
            err := os.MkdirAll(parentDir, os.ModePerm)
            if err != nil {
                return fmt.Errorf("Error al crear directorios padres: %w", err)
            }
        } else {
            return fmt.Errorf("Error: las carpetas padres no existen y el parámetro -p no fue proporcionado")
        }
    }

    // Crear el directorio
    err := os.Mkdir(mkdir.Path, 0664) // Permisos 664
    if err != nil {
        return fmt.Errorf("Error al crear el directorio: %w", err)
    }

    fmt.Println("Directorio creado con éxito:", mkdir.Path)
    return nil
}
