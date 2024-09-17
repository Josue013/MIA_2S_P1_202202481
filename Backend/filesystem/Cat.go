package filesystem

import (
    "encoding/binary"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
)

// Función para leer el contenido de archivos en la partición o en el sistema de archivos local
func Cat(fileValores []string) string {
    var respuesta string
    for _, filePath := range fileValores {
        // Verificar si el archivo existe en el sistema de archivos local
        if _, err := os.Stat(filePath); err == nil {
            // Leer el contenido del archivo local
            contenido, err := ioutil.ReadFile(filePath)
            if err != nil {
                respuesta += fmt.Sprintf("Error al leer el archivo %s\n", filePath)
                continue
            }
            respuesta += fmt.Sprintf("Contenido del archivo %s:\n%s\n", filePath, string(contenido))
            continue
        }

        // Si el archivo no existe localmente, verificar la partición montada
        indice := VerificarParticionMontada(Usr_sesion.Pid)
        if indice == -1 {
            respuesta += "Error: La partición no está montada\n"
            continue
        }
        MountActual := particionesMontadas[indice]
        superBloque := NewSuperBlock()
        archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
        if err != nil {
            respuesta += "Error al abrir el archivo\n"
            continue
        }
        defer archivo.Close()
        archivo.Seek(int64(MountActual.Start), 0)
        err = binary.Read(archivo, binary.LittleEndian, &superBloque)
        if err != nil {
            respuesta += "Error al leer el superbloque\n"
            continue
        }

        // Buscar el inodo del archivo en la partición
        numeroInodo := BuscarInode(filePath, MountActual, superBloque, archivo)
        if numeroInodo == -1 {
            respuesta += fmt.Sprintf("Error: El archivo %s no existe\n", filePath)
            continue
        }

        // Leer el contenido del archivo desde la partición
        contenido := ReadFile(numeroInodo, superBloque, archivo)
        if len(contenido) == 0 {
            respuesta += fmt.Sprintf("El archivo %s está vacío\n", filePath)
            continue
        }

        // Mostrar el contenido del archivo
        respuesta += fmt.Sprintf("Contenido del archivo %s:\n%s\n", filePath, contenido)
    }
    return respuesta
}

// Función que busca un inodo en el sistema de archivos
func BuscarInode(ruta string, MountActual Mount, superBloque SuperBlock, archivo *os.File) int {
    pathSplit := strings.Split(ruta, "/")
    var newPath []string
    for _, s := range pathSplit {
        if s != "" {
            newPath = append(newPath, s)
        }
    }

    pathSplit = newPath
    // Leer el inodo raíz
    inodoRaiz := NewInodes()
    archivo.Seek(int64(superBloque.S_inode_start), 0)
    err := binary.Read(archivo, binary.LittleEndian, &inodoRaiz)
    if err != nil {
        fmt.Println("Error al leer el inodo raíz")
        return -1
    }

    // Buscar el número de inodo del archivo
    numeroInodo := BuscarInodeRec(inodoRaiz, pathSplit, superBloque, archivo)
    return numeroInodo
}

// Función recursiva que busca un inodo en el sistema de archivos
func BuscarInodeRec(inodo Inodes, pathSplit []string, superBloque SuperBlock, archivo *os.File) int {
    if len(pathSplit) == 0 {
        return -1
    }
    actual := pathSplit[0]
    path := pathSplit[1:]
    for _, i := range inodo.I_block {
        if i != -1 {
            Desplazamiento := (superBloque.S_block_start) + (int32(i) * int32(binary.Size(FolderBlock{})))
            archivo.Seek(int64(Desplazamiento), 0)
            var folder FolderBlock
            err := binary.Read(archivo, binary.LittleEndian, &folder)
            if err != nil {
                fmt.Println("Error al leer el bloque")
                return -1
            }
            for _, j := range folder.B_content {
                if j.B_inodo != -1 && strings.Contains(string(j.B_name[:]), actual) {
                    if len(path) == 0 {
                        return int(j.B_inodo)
                    }
                    // Buscar el siguiente inodo
                    inodoSiguiente := NewInodes()
                    archivo.Seek(int64(superBloque.S_inode_start)+int64(j.B_inodo*int32(binary.Size(Inodes{}))), 0)
                    err := binary.Read(archivo, binary.LittleEndian, &inodoSiguiente)
                    if err != nil {
                        fmt.Println("Error al leer el inodo")
                        return -1
                    }
                    return BuscarInodeRec(inodoSiguiente, path, superBloque, archivo)
                }
            }
        }
    }
    return -1
}

// Función para leer el contenido de un archivo desde la partición
func ReadFile(numeroInodo int, superBloque SuperBlock, archivo *os.File) string {
    var respuesta string
    inodo := NewInodes()
    archivo.Seek(int64(superBloque.S_inode_start+int32(numeroInodo)*int32(binary.Size(Inodes{}))), 0)
    binary.Read(archivo, binary.LittleEndian, &inodo)

    if inodo.I_size == 0 {
        respuesta += "El archivo está vacío\n"
        return respuesta
    }

    // Buscar el inodo del archivo
    for _, i := range inodo.I_block {
        if i != -1 {
            Desplazamiento := superBloque.S_block_start + (int32(i) * int32(binary.Size(Fileblock{})))
            var fileBlock Fileblock
            archivo.Seek(int64(Desplazamiento), 0)
            binary.Read(archivo, binary.LittleEndian, &fileBlock)
            lectura := strings.TrimRight(string(fileBlock.B_content[:]), string(rune(0)))
            respuesta += lectura
        }
    }
    return respuesta
}