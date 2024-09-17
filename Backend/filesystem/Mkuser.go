package filesystem

import (
    "encoding/binary"
    "fmt"
    "os"
    "strconv"
    "strings"
)

// Estructura para el comando MKUSR
type MKUSR struct {
    User string
    Pass string
    Grp  string
}

// Función para crear un usuario
func commandMkusr(mkusr *MKUSR) string {
    // Verificar que haya una sesión activa y que el usuario sea root
    if Usr_sesion.Uid == -1 {
        return "Error: Necesita iniciar sesión\n"
    }
    if Usr_sesion.Usr != "root" {
        return "Error: Solo el usuario root puede crear usuarios\n"
    }

    // Verificar que la partición esté montada
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

    // Leer el superbloque
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

    contenido := LeerArchivo2(numeroInodo, superBloque, archivo)
    fmt.Println("Contenido actual de users.txt:", contenido) // Log de depuración

    if len(contenido) == 0 {
        return "Error: El archivo users.txt está vacío\n"
    }

    // Verificar si el grupo existe
    if !ExisteGrupo(contenido, mkusr.Grp) {
        return fmt.Sprintf("Error: No existe un grupo con el nombre \"%s\"\n", mkusr.Grp)
    }

    // Verificar si el usuario ya existe
    if ExisteUsuario(contenido, mkusr.User) {
        return fmt.Sprintf("Error: Ya existe un usuario con el nombre \"%s\"\n", mkusr.User)
    }

    // Agregar el nuevo usuario al archivo users.txt
    numero := ObtenerNumeroGrupo(contenido, mkusr.Grp)
    usuario := AgregarUsuario(numero, mkusr.Grp, mkusr.User, mkusr.Pass)
    nuevoContenido := contenido + usuario
    fmt.Println("Nuevo contenido a escribir:", nuevoContenido) // Log de depuración

    if EscribirArchivo2(numeroInodo, superBloque, archivo, nuevoContenido) {
        return fmt.Sprintf("Usuario \"%s\" creado con éxito\n", mkusr.User)
    }

    return "Error al crear el usuario\n"
}

// Función para agregar un usuario al formato de users.txt
func AgregarUsuario(groupNumber int, groupName string, userName string, password string) string {
    return fmt.Sprintf("%d,U,%s,%s,%s\n", groupNumber, groupName, userName, password)
}

// Función para obtener el número del grupo
func ObtenerNumeroGrupo(contenido string, groupName string) int {
    lineas := strings.Split(contenido, "\n")
    for _, linea := range lineas {
        parametros := strings.Split(linea, ",")
        if len(parametros) > 2 && parametros[1] == "G" && parametros[2] == groupName {
            numero, _ := strconv.Atoi(parametros[0])
            return numero
        }
    }
    return -1
}

// Función para verificar si el usuario existe
func ExisteUsuario(contenido string, userName string) bool {
    lineas := strings.Split(contenido, "\n")
    for _, linea := range lineas {
        parametros := strings.Split(linea, ",")
        if len(parametros) > 2 && parametros[1] == "U" && parametros[3] == userName {
            return true
        }
    }
    return false
}

// Función para escribir el contenido actualizado en el archivo
func EscribirArchivo2(numeroInodo int, superBloque SuperBlock, archivo *os.File, contenido string) bool {
    contenidoBytes := []byte(contenido)
    inodo := Inodes{}

    // Leer el inodo del archivo
    archivo.Seek(int64(superBloque.S_inode_start+int32(numeroInodo)*int32(binary.Size(Inodes{}))), 0)
    err := binary.Read(archivo, binary.LittleEndian, &inodo)
    if err != nil {
        fmt.Println("Error al leer el inodo")
        return false
    }

    // Escribir el contenido en los bloques
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

    fmt.Println("Escritura en archivo completada correctamente") // Log de depuración
    return true
}

// Función para leer el contenido del archivo
func LeerArchivo2(numeroInodo int, superBloque SuperBlock, archivo *os.File) string {
    var respuesta string
    inodo := NewInodes()

    // Leer el inodo
    archivo.Seek(int64(superBloque.S_inode_start+int32(numeroInodo)*int32(binary.Size(Inodes{}))), 0)
    binary.Read(archivo, binary.LittleEndian, &inodo)

    if inodo.I_size == 0 {
        return respuesta
    }

    // Leer el contenido del archivo desde los bloques asignados
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
