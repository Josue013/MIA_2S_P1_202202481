package filesystem

import (
    "encoding/binary"
    "fmt"
    "os"
    "strings"
)

// Estructura para el comando RMGRP
type RMGRP struct {
    Name string
}

// Función para eliminar un grupo
func commandRmgrp(rmgrp *RMGRP) string {
    // Verificar que haya una sesión activa y que el usuario sea root
    if Usr_sesion.Uid == -1 {
        return "Error: Necesita iniciar sesión\n"
    }
    if Usr_sesion.Usr != "root" {
        return "Error: Solo el usuario root puede eliminar grupos\n"
    }

    // Verificar que el grupo exista en el archivo users.txt
    indice := VerificarParticionMontada(Usr_sesion.Pid)
    if indice == -1 {
        return "Error: La partición no está montada\n"
    }
    MountActual := particionesMontadas[indice]
    superBloque := NewSuperBlock()
    archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
    if err != nil {
        return "Error al abrir el archivo\n"
    }
    defer archivo.Close()
    archivo.Seek(int64(MountActual.Start), 0)
    err = binary.Read(archivo, binary.LittleEndian, &superBloque)
    if err != nil {
        return "Error al leer el superbloque\n"
    }

    ruta := "/users.txt"
    numeroInodo := BuscarInodo(ruta, MountActual, superBloque, archivo)
    if numeroInodo == -1 {
        return "Error: El archivo users.txt no existe\n"
    }

    contenido := LeerArchivo(numeroInodo, superBloque, archivo)
    if len(contenido) == 0 {
        return "Error: El archivo users.txt está vacío\n"
    }

    if !ExisteGrupo(contenido, rmgrp.Name) {
        return fmt.Sprintf("Error: No existe un grupo con el nombre \"%s\"\n", rmgrp.Name)
    }

    // Marcar el grupo como eliminado en el archivo users.txt
    nuevoContenido := MarcarGrupoEliminado(contenido, rmgrp.Name)
    if EscribirArchivo(numeroInodo, superBloque, archivo, nuevoContenido) {
        return fmt.Sprintf("Grupo \"%s\" eliminado con éxito\n", rmgrp.Name)
    }

    return "Error al eliminar el grupo\n"
}

func MarcarGrupoEliminado(contenido string, groupName string) string {
    lineas := strings.Split(contenido, "\n")
    for i, linea := range lineas {
        parametros := strings.Split(linea, ",")
        if len(parametros) > 2 && parametros[1] == "G" && parametros[2] == groupName {
            lineas[i] = "0," + linea[2:]
            break
        }
    }
    return strings.Join(lineas, "\n")
}

// Función para escribir el contenido actualizado en el archivo
func EscribirArchivo(numeroInodo int, superBloque SuperBlock, archivo *os.File, contenido string) bool {
    // Convertir el contenido a bytes
    contenidoBytes := []byte(contenido)

    // Obtener el inodo del archivo
    inodo := Inodes{}
    archivo.Seek(int64(superBloque.S_inode_start+int32(numeroInodo)*int32(binary.Size(Inodes{}))), 0)
    err := binary.Read(archivo, binary.LittleEndian, &inodo)
    if err != nil {
        fmt.Println("Error al leer el inodo")
        return false
    }

    // Escribir el contenido en los bloques del archivo
    bloques := len(contenidoBytes) / 64
    if len(contenidoBytes)%64 != 0 {
        bloques++
    }

    for i := 0; i < bloques; i++ {
        bloque := Fileblock{}
        start := i * 64
        end := start + 64
        if end > len(contenidoBytes) {
            end = len(contenidoBytes)
        }
        copy(bloque.B_content[:], contenidoBytes[start:end])

        // Escribir el bloque en el archivo
        archivo.Seek(int64(superBloque.S_block_start+int32(inodo.I_block[i])*int32(binary.Size(Fileblock{}))), 0)
        err = binary.Write(archivo, binary.LittleEndian, &bloque)
        if err != nil {
            fmt.Println("Error al escribir el bloque")
            return false
        }
    }

    return true
}