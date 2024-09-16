package filesystem

import (
    "encoding/binary"
    "fmt"
    "os"
    "strings"
)

// Estructura para el comando MKGRP
type MKGRP struct {
    Name string
}

// Función para crear un grupo
func commandMkgrp(mkgrp *MKGRP) string {
    // Verificar que haya una sesión activa y que el usuario sea root
    if Usr_sesion.Uid == -1 {
        return "Error: Necesita iniciar sesión\n"
    }
    if Usr_sesion.Usr != "root" {
        return "Error: Solo el usuario root puede crear grupos\n"
    }

    // Verificar que el grupo no exista ya en el archivo users.txt
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
        // Inicializar el archivo users.txt si no existe
        inicializarUsersTxt(superBloque, archivo)
        numeroInodo = BuscarInodo(ruta, MountActual, superBloque, archivo)
        if numeroInodo == -1 {
            return "Error: No se pudo inicializar el archivo users.txt\n"
        }
    }

    contenido := LeerArchivo(numeroInodo, superBloque, archivo)
    if len(contenido) == 0 {
        return "Error: El archivo users.txt está vacío\n"
    }

    if ExisteGrupo(contenido, mkgrp.Name) {
        return fmt.Sprintf("Error: Ya existe un grupo con el nombre \"%s\"\n", mkgrp.Name)
    }

    // Agregar el nuevo grupo al archivo users.txt
    numero := ContarGrupos(contenido)
    grupo := AgregarGrupo(numero, mkgrp.Name)
    if AppendFile(numeroInodo, superBloque, archivo, grupo) {
        return fmt.Sprintf("Grupo \"%s\" creado con éxito\n", mkgrp.Name)
    }

    return "Error al crear el grupo\n"
}

func inicializarUsersTxt(superBloque SuperBlock, archivo *os.File) {
    contenidoInicial := "1,G,root\n1,U,root,root,123\n"
    numeroInodo := CrearArchivo("/users.txt", superBloque, archivo)
    AppendFile(numeroInodo, superBloque, archivo, contenidoInicial)
}

func CrearArchivo(ruta string, superBloque SuperBlock, archivo *os.File) int {
    // Implementar la lógica para crear un archivo en el sistema de archivos
    // y devolver el número de inodo del nuevo archivo
    

    return -1 // Placeholder
}

func AgregarGrupo(groupNumber int, groupName string) string {
    return fmt.Sprintf("%d,G,%s\n", groupNumber, groupName)
}

func ContarGrupos(contenido string) int {
    contador := 1
    lineas := strings.Split(contenido, "\n")
    for _, linea := range lineas {
        parametros := strings.Split(linea, ",")
        if len(parametros) > 1 && parametros[1] == "G" {
            contador++
        }
    }
    return contador
}

func ExisteGrupo(contenido string, groupName string) bool {
    lineas := strings.Split(contenido, "\n")
    for _, linea := range lineas {
        parametros := strings.Split(linea, ",")
        if len(parametros) > 2 && parametros[1] == "G" && parametros[2] == groupName {
            return true
        }
    }
    return false
}

func AppendFile(numeroInodo int, superBloque SuperBlock, archivo *os.File, contenido string) bool {
    inodo := NewInodes()
    archivo.Seek(int64(superBloque.S_inode_start+int32(numeroInodo)*int32(binary.Size(Inodes{}))), 0)
    binary.Read(archivo, binary.LittleEndian, &inodo)

    // Leer el contenido existente del archivo
    contenidoExistente := LeerArchivo(numeroInodo, superBloque, archivo)

    // Agregar el nuevo contenido al final del contenido existente
    nuevoContenido := contenidoExistente + contenido

    // Escribir el contenido combinado de vuelta al archivo
    for _, i := range inodo.I_block {
        if i != -1 {
            var fileBlock Fileblock
            archivo.Seek(int64(superBloque.S_block_start+int32(i)*int32(binary.Size(Fileblock{}))), 0)
            binary.Read(archivo, binary.LittleEndian, &fileBlock)
            // Agregar el contenido al bloque
            copy(fileBlock.B_content[:], nuevoContenido)
            archivo.Seek(int64(superBloque.S_block_start+int32(i)*int32(binary.Size(Fileblock{}))), 0)
            binary.Write(archivo, binary.LittleEndian, &fileBlock)
            return true
        }
    }
    return false
}