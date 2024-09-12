package filesystem

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
respuesta += Login(userValor, pwdValor, idValor) */
//Login es la funcion que permite el acceso al sistema
func Login(user string, pwd string, id string) string {
	var respuesta string
	//Verificar que el id exista en la lista de particiones montadas
	indice := VerificarParticionMontada(id)
	if indice == -1 {
		respuesta += "La particion no esta montada"
		return respuesta
	}
	//Verificar que no haya una sesion activa
	if Usr_sesion.Uid != -1 {
		respuesta += "Ya hay una sesion activa\n"
		return respuesta
	}
	MountActual := particionesMontadas[indice]
	superBloque := NewSuperBlock()
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
		respuesta += "Error al abrir el archivo\n"
		return respuesta
	}
	defer archivo.Close()
	archivo.Seek(int64(MountActual.Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &superBloque)
	if err != nil {
		respuesta += "Error al leer el superBloque\n"
		return respuesta
	}
	if !(superBloque.S_inode_size != 0 && superBloque.S_block_size != 0) {
		respuesta += "El sistema de archivos no esta formateado\n"
		return respuesta
	}
	ruta := "/users.txt"
	//Verificar que el archivo users.txt exista
	numeroInodo := BuscarInodo(ruta, MountActual, superBloque, archivo)
	if numeroInodo == -1 {
		respuesta += "El archivo users.txt no existe\n"
		return respuesta
	}

	LeerUsuarios := LeerArchivo(numeroInodo, superBloque, archivo)
	if len(LeerUsuarios) == 0 {
		respuesta += "El archivo users.txt esta vacio\n"
		return respuesta
	}

	fmt.Println("Usuarios: " + LeerUsuarios)

	//Dividir el archivo en lineas
	lineas := strings.Split(LeerUsuarios, "\n")
	//Recorrer las lineas
	for _, linea := range lineas {
		if len(linea) == 0 {
			break
		}
		if linea[2] == 'U' || linea[2] == 'u' {
			in := strings.Split(linea, ",")
			if in[3] == user && in[4] == pwd {
				uid, _ := strconv.Atoi(in[0])
				Usr_sesion.Uid = int32(uid)
				Usr_sesion.Usr = user
				Usr_sesion.Pass = pwd
				Usr_sesion.Pid = id
				Usr_sesion.Grp = in[2]
				break
			}
		}
		if len(linea) == 0 {
			break
		}
		if linea[2] == 'G' || linea[2] == 'g' {
			in := strings.Split(linea, ",")
			if in[2] == Usr_sesion.Grp {
				gid, _ := strconv.Atoi(in[0])
				Usr_sesion.Gid = int32(gid)
				break
			}
		}
	}
	respuesta += "Sesion iniciada correctamente con el usuario: " + user + "\n"
	fmt.Println("Sesion iniciada correctamente con el usuario: " + user)
	return respuesta
}

// Logout es la funcion que permite cerrar la sesion
func Logout() string {
	var respuesta string
	if Usr_sesion.Uid == -1 {
		respuesta += "No hay una sesion activa\n"
		return respuesta
	}
	Usr_sesion = NuevoUsuarioActual()
	respuesta += "Sesion cerrada correctamente\n"
	return respuesta
}

// Funcion que busca un inodo en el sistema de archivos
func BuscarInodo(ruta string, MountActual Mount, superBloque SuperBlock, archivo *os.File) int {
	pathSplit := strings.Split(ruta, "/")
	var newPath []string
	for _, s := range pathSplit {
		if s != "" {
			newPath = append(newPath, s)
		}
	}

	pathSplit = newPath
	//Leer el inodo raíz
	inodoRaiz := NewInodes()
	archivo.Seek(int64(superBloque.S_inode_start), 0)
	err := binary.Read(archivo, binary.LittleEndian, &inodoRaiz)
	if err != nil {
		fmt.Println("Error al leer el inodo raíz")
		return -1
	}

	//Buscar el numero de inodo del archivo
	numeroInodo := BuscarInodoRec(inodoRaiz, pathSplit, superBloque, archivo)
	return numeroInodo
}

// Funcion recursiva que busca un inodo en el sistema de archivos
func BuscarInodoRec(inodo Inodes, pathSplit []string, superBloque SuperBlock, archivo *os.File) int {
	contador := 0
	if len(pathSplit) == 0 {
		return contador
	}
	actual := pathSplit[0]
	path := pathSplit[1:]
	for _, i := range inodo.I_block {
		if i != -1 {
			Desplazamiento := (superBloque.S_block_start) + (int32(i) * int32(binary.Size(Fileblock{})))
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
					//Bucar el siguiente inodo
					inodoSiguiente := NewInodes()
					archivo.Seek(int64(superBloque.S_inode_start)+int64(j.B_inodo*int32(binary.Size(Inodes{}))), 0)
					err := binary.Read(archivo, binary.LittleEndian, &inodoSiguiente)
					if err != nil {
						fmt.Println("Error al leer el inodo")
						return -1
					}
					return BuscarInodoRec(inodoSiguiente, path, superBloque, archivo)
				}
			}
		}
	}
	return -1
}

func LeerArchivo(numeroInodo int, superBloque SuperBlock, archivo *os.File) string {
	var respuesta string
	inodo := NewInodes()
	archivo.Seek(int64(superBloque.S_inode_start+int32(numeroInodo)*int32(binary.Size(Inodes{}))), 0)
	binary.Read(archivo, binary.LittleEndian, &inodo)

	if inodo.I_size == 0 {
		respuesta += "El archivo esta vacio\n"
		return respuesta
	}

	//Buscar el inodo del archivo
	for _, i := range inodo.I_block {
		if i != -1 {
			Desplazamiento := superBloque.S_block_start + (int32(i) * int32(binary.Size(Fileblock{})))
			var fileBlock Fileblock
			archivo.Seek(int64(Desplazamiento), 0)
			binary.Read(archivo, binary.LittleEndian, &fileBlock)
			lectura := strings.TrimRight(string(fileBlock.B_content[:]), string(rune(0)))
			lectura = ObtenerCadena(lectura, 64)
			respuesta += lectura
		}
	}
	return respuesta
}

// Funcion que obtiene una cadena de un arreglo de bytes
func ObtenerCadena(cadena string, size int) string {
	contenidoFinal := ""
	CantidadDeCaracteres := len(cadena)
	if CantidadDeCaracteres < size {
		contenidoFinal = cadena
	} else {
		for i := 0; i < size; i++ {
			contenidoFinal += string(cadena[0])
			cadena = cadena[1:]
		}
	}
	return contenidoFinal
}
