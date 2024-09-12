package filesystem

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"time"
)

// Mkfs crea un sistema de archivos en la particion indicada
/* respuesta += Mkfs(typeValor, idValor, fsValor) */

func Mkfs(typeValor string, idValor string, fsValor string) string {
	var respuesta string

	indice := VerificarParticionMontada(idValor)
	if indice == -1 {
		respuesta += "La particion no esta montada"
		return respuesta
	}

	MountActual := particionesMontadas[indice]

	if MountActual.Part_type != [1]byte{'p'} {
		respuesta += "La particion no es primaria"
		return respuesta
	}

	//Cantidad de estructuras que caben en la particion
	var n int
	if fsValor == "2fs" {
		n = int(math.Floor(float64(int(MountActual.Size)-int(binary.Size(SuperBlock{}))) / float64(4+int(binary.Size(Inodes{}))+3*int(binary.Size(Fileblock{})))))
	} /* else {
		n = int(math.Floor(float64(int(MountActual.Size)-int(binary.Size(SuperBlock{}))) / float64(4+int(binary.Size(Inodes{}))+3*int(binary.Size(FolderBlock{}))+binary.Size(Journal{}))))
	} */

	//Crear el superbloque
	sb := NewSuperBlock()
	sb.S_inodes_count = int32(n)
	sb.S_blocks_count = int32(n * 3)
	sb.S_free_blocks_count = int32(n * 3)
	sb.S_free_inodes_count = int32(n)
	fechaActual := time.Now()
	fecha := fechaActual.Format("2006-01-02 15:04:05")
	copy(sb.S_mtime[:], fecha)
	copy(sb.S_umtime[:], fecha)
	sb.S_mnt_count = 1
	if fsValor == "2fs" {
		respuesta += Crear2FS(sb, MountActual, n)
		return respuesta
	} /* else {
		Crear3FS(sb, MountActual, n)
	} */
	return respuesta
}

// Crear2FS crea un sistema de archivos 2fs
func Crear2FS(sb SuperBlock, MountActual Mount, n int) string {

	sb.S_filesystem_type = 2
	sb.S_bm_inode_start = int32(MountActual.Start) + int32(binary.Size(SuperBlock{}))
	sb.S_bm_block_start = sb.S_bm_inode_start + int32(n)
	sb.S_inode_start = sb.S_bm_block_start + int32(3*n)
	sb.S_block_start = sb.S_inode_start + int32(n*int(binary.Size(Inodes{})))
	//Crear el bloque 0, inodo 0 y el usuario root
	sb.S_free_blocks_count--
	sb.S_free_inodes_count--
	sb.S_free_blocks_count--
	sb.S_free_inodes_count--

	//Creación del super bloque
	//Abrir el archivo

	file, err := os.OpenFile(MountActual.Path, os.O_WRONLY, 0664)
	if err != nil {
		println("Error al abrir el archivo")
		return "Error al abrir el archivo"
	}
	defer file.Close()

	file.Seek(int64(MountActual.Start), 0)
	binary.Write(file, binary.LittleEndian, &sb)

	//Crear el bitmap de inodos
	var llenar byte = 0
	file.Seek(int64(sb.S_bm_inode_start), 0)
	for i := 0; i < n; i++ {
		binary.Write(file, binary.LittleEndian, &llenar)
	}

	//Crear el bitmap de bloques
	file.Seek(int64(sb.S_bm_block_start), 0)
	for i := 0; i < n*3; i++ {
		binary.Write(file, binary.LittleEndian, &llenar)
	}

	//Crear el inodo 0
	inodo0 := NewInodes()

	//Crear el bloque 0
	var bloque0 Fileblock

	//Formatear inodos
	file.Seek(int64(sb.S_inode_start), 0)
	for i := 0; i < n; i++ {
		binary.Write(file, binary.LittleEndian, &inodo0)
	}

	//Formatear bloques
	file.Seek(int64(sb.S_block_start), 0)
	for i := 0; i < n*3; i++ {
		binary.Write(file, binary.LittleEndian, &bloque0)
	}

	//Crear el directorio raíz
	//Crear el inodo
	inodo0.I_uid = 1
	inodo0.I_gid = 1
	fechaActual := time.Now()
	fecha := fechaActual.Format("2006-01-02 15:04:05")
	copy(inodo0.I_atime[:], fecha)
	copy(inodo0.I_ctime[:], fecha)
	copy(inodo0.I_mtime[:], fecha)
	inodo0.I_type = [1]byte{'0'}
	inodo0.I_perm = 664
	inodo0.I_block[0] = 0

	//Crear el bloque carpeta

	var bloqueCarpeta FolderBlock
	bloqueCarpeta.B_content[0].B_inodo = 0
	copy(bloqueCarpeta.B_content[0].B_name[:], ".")
	bloqueCarpeta.B_content[1].B_inodo = 0
	copy(bloqueCarpeta.B_content[1].B_name[:], "..")
	bloqueCarpeta.B_content[2].B_inodo = 1
	copy(bloqueCarpeta.B_content[2].B_name[:], "users.txt")
	bloqueCarpeta.B_content[3].B_inodo = -1

	data := "1,G,root\n1,U,root,root,123\n"

	//Escribir el inodo y el bloque en el archivo

	inodo1 := NewInodes()
	inodo1.I_uid = 1
	inodo1.I_gid = 1
	fechaActual = time.Now()
	fecha = fechaActual.Format("2006-01-02 15:04:05")
	copy(inodo1.I_atime[:], fecha)
	copy(inodo1.I_ctime[:], fecha)
	copy(inodo1.I_mtime[:], fecha)
	inodo1.I_type = [1]byte{'1'}
	inodo1.I_perm = 664
	inodo1.I_block[0] = 1
	inodo1.I_size = int32(len(data)) + int32(binary.Size(Fileblock{}))

	inodo0.I_size = inodo1.I_size + int32(binary.Size(FolderBlock{})) + int32(binary.Size(FolderBlock{}))

	var bloqueArchivo Fileblock
	copy(bloqueArchivo.B_content[:], data)

	//Escribir el inodo en el archivo
	file.Seek(int64(sb.S_bm_inode_start), 0)
	var bit byte = 1
	binary.Write(file, binary.LittleEndian, &bit)
	binary.Write(file, binary.LittleEndian, &bit)

	file.Seek(int64(sb.S_bm_block_start), 0)
	binary.Write(file, binary.LittleEndian, &bit)
	binary.Write(file, binary.LittleEndian, &bit)

	file.Seek(int64(sb.S_inode_start), 0)
	binary.Write(file, binary.LittleEndian, &inodo0)
	binary.Write(file, binary.LittleEndian, &inodo1)

	file.Seek(int64(sb.S_block_start), 0)
	binary.Write(file, binary.LittleEndian, &bloqueCarpeta)
	binary.Write(file, binary.LittleEndian, &bloqueArchivo)

	fmt.Println("Sistema de archivos 2FS creado con éxito en el disco: " + MountActual.Id)
	return "Sistema de archivos 2FS creado con éxito en el disco: " + MountActual.Id

}