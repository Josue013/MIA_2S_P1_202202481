package filesystem

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// respuesta += ReporteDisk(idValor, pathValor)
func ReporteDisk(id string, pathValor string) string {
	var respuesta string

	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Nombre del archivo: " + fileName)
	fmt.Println("Ruta del archivo: " + dirPath)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		respuesta += "Error al crear el directorio\n"
		fmt.Println("Error al crear el directorio")
		return respuesta
	}

	//Buscar la particion montada
	indice := VerificarParticionMontada(id)
	if indice == -1 {
		respuesta += "La particion no esta montada"
		return respuesta
	}

	MountActual := particionesMontadas[indice]

	//Abrir el archivo
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
		respuesta += "Error al abrir el archivo\n"
		fmt.Println("Error al abrir el archivo")
		return respuesta
	}
	defer archivo.Close()
	Dot := "digraph grid {bgcolor=\"darkslategray\" label=\" Reporte Disk \"layout=dot compound=true "
	Dot += "labelloc = \"t\"edge [weigth=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"]"
	Dot += "node[shape=record, color=darkslategray1, fontname=\"Helvetica\"]a0[label=\"MBR"

	//Leer el MBR
	disk := MBR{}
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		respuesta += "Error al leer el MBR\n"
		fmt.Println("Error al leer el MBR")
		return respuesta
	}
	sizeMBR := int(disk.Mbr_tamano)
	libreMBR := int(disk.Mbr_tamano)

	//Crear el MBR
	if disk.Mbr_partition_1.Part_size != 0 {
		libreMBR -= int(disk.Mbr_partition_1.Part_size)
		Dot += "|"
		if disk.Mbr_partition_1.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition_1.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition_1.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			ebr := NewEBR()
			Desplazamiento := int(disk.Mbr_partition_1.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &ebr)
			if err != nil {
				fmt.Println("Error al leer el EBR")
				respuesta += "Error al leer el EBR\n"
				return respuesta
			}
			if ebr.Part_size != 0 {
				Dot += "|{"
				PrimerEBR := true
				for {
					if !PrimerEBR {
						Dot += "|EBR"
					} else {
						PrimerEBR = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					fmt.Println("Nombre de la particion: " + string(ebr.Part_name[:]))
					porcentaje := (float64(ebr.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libreExtendida -= int(ebr.Part_size)

					Desplazamiento += int(ebr.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					err = binary.Read(archivo, binary.LittleEndian, &ebr)
					if err != nil {
						fmt.Println("Error al leer el EBR")
						respuesta += "Error al leer el EBR\n"
						return respuesta
					}
					if ebr.Part_size == 0 {
						break
					}
				}

				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}

	}
	if disk.Mbr_partition_2.Part_size != 0 {
		libreMBR -= int(disk.Mbr_partition_2.Part_size)
		Dot += "|"
		if disk.Mbr_partition_2.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition_2.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition_2.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			ebr := NewEBR()
			Desplazamiento := int(disk.Mbr_partition_2.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &ebr)
			if err != nil {
				fmt.Println("Error al leer el EBR")
				respuesta += "Error al leer el EBR\n"
				return respuesta
			}
			if ebr.Part_size != 0 {
				Dot += "|{"
				PrimerEBR := true
				for {
					if !PrimerEBR {
						Dot += "|EBR"
					} else {
						PrimerEBR = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					fmt.Println("Nombre de la particion: " + string(ebr.Part_name[:]))
					porcentaje := (float64(ebr.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libreExtendida -= int(ebr.Part_size)

					Desplazamiento += int(ebr.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					err = binary.Read(archivo, binary.LittleEndian, &ebr)
					if err != nil {
						fmt.Println("Error al leer el EBR")
						respuesta += "Error al leer el EBR\n"
						return respuesta
					}
					if ebr.Part_size == 0 {
						break
					}
				}

				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}
	}
	if disk.Mbr_partition_3.Part_size != 0 {
		libreMBR -= int(disk.Mbr_partition_3.Part_size)
		Dot += "|"
		if disk.Mbr_partition_3.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition_3.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition_3.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			ebr := NewEBR()
			Desplazamiento := int(disk.Mbr_partition_3.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &ebr)
			if err != nil {
				fmt.Println("Error al leer el EBR")
				respuesta += "Error al leer el EBR\n"
				return respuesta
			}
			if ebr.Part_size != 0 {
				Dot += "|{"
				PrimerEBR := true
				for {
					if !PrimerEBR {
						Dot += "|EBR"
					} else {
						PrimerEBR = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					fmt.Println("Nombre de la particion: " + string(ebr.Part_name[:]))
					porcentaje := (float64(ebr.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libreExtendida -= int(ebr.Part_size)

					Desplazamiento += int(ebr.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					err = binary.Read(archivo, binary.LittleEndian, &ebr)
					if err != nil {
						fmt.Println("Error al leer el EBR")
						respuesta += "Error al leer el EBR\n"
						return respuesta
					}
					if ebr.Part_size == 0 {
						break
					}
				}

				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}

	}
	if disk.Mbr_partition_4.Part_size != 0 {
		libreMBR -= int(disk.Mbr_partition_4.Part_size)
		Dot += "|"
		if disk.Mbr_partition_4.Part_type == [1]byte{'p'} {
			Dot += "Primaria"
			porcentaje := (float64(disk.Mbr_partition_4.Part_size) * float64(100)) / float64(sizeMBR)
			Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
		} else {
			libreExtendida := int(disk.Mbr_partition_4.Part_size)
			Dot += "{Extendida"
			//Leer el EBR
			ebr := NewEBR()
			Desplazamiento := int(disk.Mbr_partition_4.Part_start)
			archivo.Seek(int64(Desplazamiento), 0)
			err = binary.Read(archivo, binary.LittleEndian, &ebr)
			if err != nil {
				fmt.Println("Error al leer el EBR")
				respuesta += "Error al leer el EBR\n"
				return respuesta
			}
			if ebr.Part_size != 0 {
				Dot += "|{"
				PrimerEBR := true
				for {
					if !PrimerEBR {
						Dot += "|EBR"
					} else {
						PrimerEBR = false
						Dot += "EBR"
					}
					Dot += "|Logica"
					fmt.Println("Nombre de la particion: " + string(ebr.Part_name[:]))
					porcentaje := (float64(ebr.Part_size) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					libreExtendida -= int(ebr.Part_size)

					Desplazamiento += int(ebr.Part_size) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(Desplazamiento), 0)
					err = binary.Read(archivo, binary.LittleEndian, &ebr)
					if err != nil {
						fmt.Println("Error al leer el EBR")
						respuesta += "Error al leer el EBR\n"
						return respuesta
					}
					if ebr.Part_size == 0 {
						break
					}
				}

				if libreExtendida > 0 {
					Dot += "|Libre"
					porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				}
				Dot += "}}"
			} else {
				Dot += "|Libre"
				porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
				Dot += "}"
			}
		}
	}
	if libreMBR > 0 {
		Dot += "|Libre"
		porcentaje := (float64(libreMBR) * float64(100)) / float64(sizeMBR)
		Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
	}

	Dot += "\"];\n}"
	//Crear el archivo dot
	//-path=/home/user/reports/report2.pdf
	extension := path.Ext(pathValor) //Obtener la extension
	//Archivo sin extension
	fileName = strings.TrimSuffix(fileName, extension) //Quitar la extension
	DotName := dirPath + fileName + ".dot"
	//Crear el archivo .dot
	file, err := os.Create(DotName)
	if err != nil {
		fmt.Println("Error al crear el archivo .dot")
		respuesta += "Error al crear el archivo .dot\n"
		return respuesta
	}
	defer file.Close()
	//Escribir el archivo .dot
	_, err = file.WriteString(Dot)
	if err != nil {
		fmt.Println("Error al escribir el archivo .dot")
		respuesta += "Error al escribir el archivo .dot\n"
		return respuesta
	}
	fmt.Println("Archivo .dot creado")

	//Quitar el punto a la extension
	extension = extension[1:]

	//Crear el reporte
	cmd := exec.Command("dot", "-T", extension, DotName, "-o", pathValor)
	fmt.Println("dot -T" + extension + " " + DotName + " -o " + pathValor)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al crear el reporte con extension")
		respuesta += "Error al crear el reporte con extension\n"
		return respuesta
	}

	return "Reporte Disk creado con exito\n"
}

// Reporte SuperBlock
func ReporteSB(id string, pathValor string) string {
	var respuesta string

	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Nombre del archivo: " + fileName)
	fmt.Println("Ruta del archivo: " + dirPath)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		respuesta += "Error al crear el directorio\n"
		fmt.Println("Error al crear el directorio")
		return respuesta
	}

	//Buscar la particion montada
	indice := VerificarParticionMontada(id)
	if indice == -1 {
		respuesta += "La particion no esta montada"
		return respuesta
	}

	MountActual := particionesMontadas[indice]
	//Abrir el archivo
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
		respuesta += "Error al abrir el archivo\n"
		fmt.Println("Error al abrir el archivo")
		return respuesta
	}
	defer archivo.Close()
	//Leer el superbloque
	superBloque := NewSuperBlock()
	archivo.Seek(int64(MountActual.Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &superBloque)
	if err != nil {
		respuesta += "Error al leer el superbloque\n"
		fmt.Println("Error al leer el superbloque")
		return respuesta
	}
	Dot := "digraph grid {bgcolor=\"#00441b\" fontcolor=\"white\" label=\" Reporte SuperBlock \"layout=dot "
	Dot += "labelloc = \"t\"edge [weigth=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"]"
	Dot += "a0[shape=none, color=cyan, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"white\">SuperBlock</TD><TD></TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"springgreen3\">s_filesystem_type</TD><TD bgcolor=\"springgreen3\">" + strconv.Itoa(int(superBloque.S_filesystem_type)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"white\">s_inodes_count</TD><TD bgcolor=\"white\">" + strconv.Itoa(int(superBloque.S_inodes_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"springgreen3\">s_blocks_count</TD><TD bgcolor=\"springgreen3\">" + strconv.Itoa(int(superBloque.S_blocks_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"white\">s_free_blocks_count</TD><TD bgcolor=\"white\">" + strconv.Itoa(int(superBloque.S_free_blocks_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"springgreen3\">s_free_inodes_count</TD><TD bgcolor=\"springgreen3\">" + strconv.Itoa(int(superBloque.S_free_inodes_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"white\">s_mtime</TD><TD bgcolor=\"white\">" + string(superBloque.S_mtime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"springgreen3\">s_umtime</TD><TD bgcolor=\"springgreen3\">" + string(superBloque.S_umtime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"white\">s_mnt_count</TD><TD bgcolor=\"white\">" + strconv.Itoa(int(superBloque.S_mnt_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"springgreen3\">s_magic</TD><TD bgcolor=\"springgreen3\">" + strconv.Itoa(int(superBloque.S_magic)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"white\">s_inode_size</TD><TD bgcolor=\"white\">" + strconv.Itoa(int(superBloque.S_inode_size)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"springgreen3\">s_block_size</TD><TD bgcolor=\"springgreen3\">" + strconv.Itoa(int(superBloque.S_block_size)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"white\">s_first_ino</TD><TD bgcolor=\"white\">" + strconv.Itoa(int(superBloque.S_first_ino)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"springgreen3\">s_first_blo</TD><TD bgcolor=\"springgreen3\">" + strconv.Itoa(int(superBloque.S_first_blo)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"white\">s_bm_inode_start</TD><TD bgcolor=\"white\">" + strconv.Itoa(int(superBloque.S_bm_inode_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"springgreen3\">s_bm_block_start</TD><TD bgcolor=\"springgreen3\">" + strconv.Itoa(int(superBloque.S_bm_block_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"white\">s_inode_start</TD><TD bgcolor=\"white\">" + strconv.Itoa(int(superBloque.S_inode_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"springgreen3\">s_block_start</TD><TD bgcolor=\"springgreen3\">" + strconv.Itoa(int(superBloque.S_block_start)) + "</TD></TR>\n"
	Dot += "</TABLE>>];\n"
	Dot += "}"

	//Crear el archivo dot
	extension := path.Ext(pathValor)
	//Archivo sin extension
	fileName = strings.TrimSuffix(fileName, extension)
	DotName := dirPath + fileName + ".dot"
	//Crear el archivo .dot
	file, err := os.Create(DotName)
	if err != nil {
		fmt.Println("Error al crear el archivo .dot")
		respuesta += "Error al crear el archivo .dot\n"
		return respuesta
	}
	defer file.Close()
	//Escribir el archivo .dot
	_, err = file.WriteString(Dot)
	if err != nil {
		fmt.Println("Error al escribir el archivo .dot")
		respuesta += "Error al escribir el archivo .dot\n"
		return respuesta
	}
	fmt.Println("Archivo .dot creado")

	//Quitar el punto a la extension
	extension = extension[1:]

	//Crear el reporte
	cmd := exec.Command("dot", "-T", extension, DotName, "-o", pathValor)
	fmt.Println("dot -T " + extension + " " + DotName + " -o " + pathValor)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al crear el reporte")
		respuesta += "Error al crear el reporte\n"
		return respuesta
	}

	return "Reporte SuperBlock creado con exito\n"

}

func ReporteFile(idValor string, pathValor string, rutaValor string) string {
	var respuesta string
	//Buscar la particion montada
	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Nombre del archivo: " + fileName)
	fmt.Println("Ruta del archivo: " + dirPath)
	//Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		respuesta += "Error al crear el directorio\n"
		fmt.Println("Error al crear el directorio")
		return respuesta
	}

	indice := VerificarParticionMontada(idValor)
	if indice == -1 {
		respuesta += "La particion no esta montada"
		return respuesta
	}
	MountActual := particionesMontadas[indice]
	//Abrir el archivo
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
		respuesta += "Error al abrir el archivo\n"
		fmt.Println("Error al abrir el archivo")
		return respuesta
	}
	defer archivo.Close()
	//Leer el superbloque
	superBloque := NewSuperBlock()
	archivo.Seek(int64(MountActual.Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &superBloque)
	if err != nil {
		respuesta += "Error al leer el superbloque\n"
		fmt.Println("Error al leer el superbloque")
		return respuesta
	}
	//Buscar el inodo de la ruta
	numeroInodo := BuscarInodo(rutaValor, MountActual, superBloque, archivo)
	if numeroInodo == -1 {
		respuesta += "La ruta no existe\n"
		return respuesta
	}
	//Leer el inodo
	cadena := LeerArchivo(numeroInodo, superBloque, archivo)
	if len(cadena) == 0 {
		respuesta += "El archivo esta vacio\n"
		return respuesta
	}
	Dot := "digraph G{\n"
	Dot += "a[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD colspan=\"2\" bgcolor=\"lightgrey\" >" + rutaValor + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">Contenido</TD></TR>\n"
	Dot += "<TR><TD>" + cadena + "</TD></TR>\n"
	Dot += "</TABLE>>];\n}"
	//Crear el archivo dot
	extension := path.Ext(pathValor)

	if extension == ".txt" {
		//Crear el archivo .txt
		file, err := os.Create(pathValor)
		if err != nil {
			fmt.Println("Error al crear el archivo .txt")
			respuesta += "Error al crear el archivo .txt\n"
			return respuesta
		}
		defer file.Close()

		//Escribir el archivo .txt
		_, err = file.WriteString(cadena)
		if err != nil {
			fmt.Println("Error al escribir el archivo .txt")
			respuesta += "Error al escribir el archivo .txt\n"
			return respuesta
		}
		fmt.Println("Archivo .txt creado")
		return "Reporte File creado con exito\n"

	} else {
		//Archivo sin extension
		fileName = strings.TrimSuffix(fileName, extension)
		DotName := dirPath + fileName + ".dot"
		//Crear el archivo .dot
		file, err := os.Create(DotName)
		if err != nil {
			fmt.Println("Error al crear el archivo .dot")
			respuesta += "Error al crear el archivo .dot\n"
			return respuesta
		}
		defer file.Close()

		//Escribir el archivo .dot
		_, err = file.WriteString(Dot)
		if err != nil {
			fmt.Println("Error al escribir el archivo .dot")
			respuesta += "Error al escribir el archivo .dot\n"
			return respuesta
		}
		fmt.Println("Archivo .dot creado")

		//Quitar el punto a la extension
		extension = extension[1:]

		//Crear el reporte
		cmd := exec.Command("dot", "-T", extension, DotName, "-o", pathValor)
		fmt.Println("dot -T " + extension + " " + DotName + " -o " + pathValor)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error al crear el reporte")
			respuesta += "Error al crear el reporte\n"
			return respuesta
		}

		return "Reporte File creado con exito\n"
	}

}

// REPORTE MBR
func ReporteMBR(id string, pathValor string) string {
	var respuesta string

	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Nombre del archivo: " + fileName)
	fmt.Println("Ruta del archivo: " + dirPath)
	// Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
			respuesta += "Error al crear el directorio\n"
			fmt.Println("Error al crear el directorio")
			return respuesta
	}

	// Buscar la particion montada
	indice := VerificarParticionMontada(id)
	if indice == -1 {
			respuesta += "La particion no esta montada"
			return respuesta
	}

	MountActual := particionesMontadas[indice]

	// Abrir el archivo
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
			respuesta += "Error al abrir el archivo\n"
			fmt.Println("Error al abrir el archivo")
			return respuesta
	}
	defer archivo.Close()

	// Leer el MBR
	disk := MBR{}
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
			respuesta += "Error al leer el MBR\n"
			fmt.Println("Error al leer el MBR")
			return respuesta
	}

	// Crear el contenido DOT
	dotContent := fmt.Sprintf(`digraph G {
			node [shape=plaintext]
			tabla [label=<
					<table border="0" cellborder="1" cellspacing="0">
							<tr><td colspan="2" bgcolor="indigo"><font color="white">REPORTE MBR</font></td></tr>
							<tr><td>mbr_tamano</td><td>%d</td></tr>
							<tr><td bgcolor="plum2">mbr_fecha_creacion</td><td bgcolor="plum2">%s</td></tr>
							<tr><td>mbr_disk_signature</td><td>%d</td></tr>
	`, disk.Mbr_tamano, string(disk.Mbr_fecha_creacion[:]), disk.Mbr_disk_signature)

	// Agregar las particiones a la tabla
	particiones := []Partition{disk.Mbr_partition_1, disk.Mbr_partition_2, disk.Mbr_partition_3, disk.Mbr_partition_4}
	for i, part := range particiones {
			if part.Part_size == 0 {
					continue
			}

			partName := strings.TrimRight(string(part.Part_name[:]), "\x00")
			partStatus := rune(part.Part_status[0])
			partType := rune(part.Part_type[0])
			partFit := rune(part.Part_fit[0])

			dotContent += fmt.Sprintf(`
					<tr><td colspan="2" bgcolor="#f07d7d"> PARTICIÓN %d </td></tr>
					<tr><td>part_status</td><td>%c</td></tr>
					<tr><td bgcolor="#f5b4af">part_type</td><td bgcolor="#f5b4af">%c</td></tr>
					<tr><td>part_fit</td><td>%c</td></tr>
					<tr><td bgcolor="#f5b4af">part_start</td><td bgcolor="#f5b4af">%d</td></tr>
					<tr><td>part_size</td><td>%d</td></tr>
					<tr><td bgcolor="#f5b4af">part_name</td><td bgcolor="#f5b4af">%s</td></tr>
			`, i+1, partStatus, partType, partFit, part.Part_start, part.Part_size, partName)

			if (partType == 'E' || partType == 'e') {
					dotContent += recorrerEBR(MountActual.Path, part.Part_start)
			}
	}

	// Cerrar la tabla y el contenido DOT
	dotContent += "</table>>] }"

	// Guardar el contenido DOT en un archivo
	dotFileName := dirPath + fileName + ".dot"
	file, err := os.Create(dotFileName)
	if err != nil {
			fmt.Println("Error al crear el archivo .dot")
			respuesta += "Error al crear el archivo .dot\n"
			return respuesta
	}
	defer file.Close()

	_, err = file.WriteString(dotContent)
	if err != nil {
			fmt.Println("Error al escribir el archivo .dot")
			respuesta += "Error al escribir el archivo .dot\n"
			return respuesta
	}
	fmt.Println("Archivo .dot creado")

	// Quitar el punto a la extension
	extension := path.Ext(pathValor)[1:]

	// Crear el reporte
	cmd := exec.Command("dot", "-T", extension, dotFileName, "-o", pathValor)
	fmt.Println("dot -T" + extension + " " + dotFileName + " -o " + pathValor)
	err = cmd.Run()
	if err != nil {
			fmt.Println("Error al crear el reporte con extension")
			respuesta += "Error al crear el reporte con extension\n"
			return respuesta
	}

	return "Reporte MBR creado con exito\n"
}

func recorrerEBR(ruta string, whereToStart int32) string {
	contenido := ""
	var temp EBR
	archivo, err := os.OpenFile(ruta, os.O_RDWR, 0664)
	if err != nil {
			fmt.Println("Error al abrir el archivo para leer EBR")
			return ""
	}
	defer archivo.Close()

	archivo.Seek(int64(whereToStart), 0)
	err = binary.Read(archivo, binary.LittleEndian, &temp)
	if err != nil {
			fmt.Println("Error al leer el EBR")
			return ""
	}

	flag := true
	for flag {
			if temp.Part_size == 0 {
					flag = false
			} else if temp.Part_next != -1 && temp.Part_mount[0] != '5' {
					contenido += "\t\t\t<TR><TD bgcolor=\"pink\" COLSPAN=\"2\">Particion Logica</TD></TR>\n"
					contenido += "\t\t\t<TR><TD> part_status </TD><TD>"
					contenido += string(temp.Part_mount[:])
					contenido += "</TD></TR>\n"
					contenido += "\t\t\t<TR><TD bgcolor=\"#D3D3D3\"> part_next </TD><TD bgcolor=\"#D3D3D3\">"
					contenido += strconv.Itoa(int(temp.Part_next))
					contenido += "</TD></TR>\n"
					contenido += "\t\t\t<TR><TD> part_fit </TD><TD>"
					contenido += string(temp.Part_fit[:])
					contenido += "</TD></TR>\n"
					contenido += "\t\t\t<TR><TD bgcolor=\"#D3D3D3\"> part_start </TD><TD bgcolor=\"#D3D3D3\">" + strconv.Itoa(int(temp.Part_start)) + "</TD></TR>\n"
					contenido += "\t\t\t<TR><TD> part_size </TD><TD>" + strconv.Itoa(int(temp.Part_size)) + "</TD></TR>\n"
					contenido += "\t\t\t<TR><TD bgcolor=\"#D3D3D3\"> part_name </TD><TD bgcolor=\"#D3D3D3\">" + strings.TrimRight(string(temp.Part_name[:]), "\x00") + "</TD></TR>\n"
			} else if temp.Part_next == -1 {
					contenido += "\t\t\t<TR><TD bgcolor=\"pink\" COLSPAN=\"2\">Particion Logica</TD></TR>\n"
					contenido += "\t\t\t<TR><TD> part_status </TD><TD>"
					contenido += string(temp.Part_mount[:])
					contenido += "</TD></TR>\n"
					contenido += "\t\t\t<TR><TD> part_fit </TD><TD>"
					contenido += string(temp.Part_fit[:])
					contenido += "</TD></TR>\n"
					contenido += "\t\t\t<TR><TD bgcolor=\"#D3D3D3\"> part_start </TD><TD bgcolor=\"#D3D3D3\">" + strconv.Itoa(int(temp.Part_start)) + "</TD></TR>\n"
					contenido += "\t\t\t<TR><TD> part_size </TD><TD>" + strconv.Itoa(int(temp.Part_size)) + "</TD></TR>\n"
					contenido += "\t\t\t<TR><TD bgcolor=\"#D3D3D3\"> part_name </TD><TD bgcolor=\"#D3D3D3\">" + strings.TrimRight(string(temp.Part_name[:]), "\x00") + "</TD></TR>\n"
					flag = false
			}
			if temp.Part_next != -1 {
					archivo.Seek(int64(temp.Part_next), 0)
					err = binary.Read(archivo, binary.LittleEndian, &temp)
					if err != nil {
							fmt.Println("Error al leer el siguiente EBR")
							return ""
					}
			}
	}
	return contenido
}

// REPORTE INODE
func ReporteInode(id string, pathValor string) string {
	var respuesta string

	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Nombre del archivo: " + fileName)
	fmt.Println("Ruta del archivo: " + dirPath)
	// Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		respuesta += "Error al crear el directorio\n"
		fmt.Println("Error al crear el directorio")
		return respuesta
	}

	// Buscar la particion montada
	indice := VerificarParticionMontada(id)
	if indice == -1 {
		respuesta += "La particion no esta montada"
		return respuesta
	}

	MountActual := particionesMontadas[indice]

	// Abrir el archivo
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
		respuesta += "Error al abrir el archivo\n"
		fmt.Println("Error al abrir el archivo")
		return respuesta
	}
	defer archivo.Close()

	// Leer el superbloque
	superBloque := NewSuperBlock()
	archivo.Seek(int64(MountActual.Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &superBloque)
	if err != nil {
		respuesta += "Error al leer el superbloque\n"
		fmt.Println("Error al leer el superbloque")
		return respuesta
	}

	// Leer el bitmap de inodos
	bitmap := make([]byte, superBloque.S_inodes_count)
	archivo.Seek(int64(superBloque.S_bm_inode_start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &bitmap)
	if err != nil {
		respuesta += "Error al leer el bitmap de inodos\n"
		fmt.Println("Error al leer el bitmap de inodos")
		return respuesta
	}

	// Iniciar el contenido DOT
	dotContent := `digraph G {
			rankdir=LR;
			node [shape=plaintext]
	`

	// Iterar sobre cada inodo utilizado
	for i := int32(0); i < superBloque.S_inodes_count; i++ {
		if bitmap[i] == 0 {
			continue
		}

		inode := NewInodes()
		// Leer el inodo
		archivo.Seek(int64(superBloque.S_inode_start+(i*superBloque.S_inode_size)), 0)
		err := binary.Read(archivo, binary.LittleEndian, &inode)
		if err != nil {
			respuesta += "Error al leer el inodo\n"
			fmt.Println("Error al leer el inodo")
			return respuesta
		}

		// Convertir tiempos a string
		atime := string(inode.I_atime[:])
		ctime := string(inode.I_ctime[:])
		mtime := string(inode.I_mtime[:])

		// Definir el contenido DOT para el inodo actual
		dotContent += fmt.Sprintf(`inode%d [label=<
					<table border="0" cellborder="1" cellspacing="0">
							<tr><td colspan="2"> REPORTE INODO %d </td></tr>
							<tr><td>i_uid</td><td>%d</td></tr>
							<tr><td>i_gid</td><td>%d</td></tr>
							<tr><td>i_size</td><td>%d</td></tr>
							<tr><td>i_atime</td><td>%s</td></tr>
							<tr><td>i_ctime</td><td>%s</td></tr>
							<tr><td>i_mtime</td><td>%s</td></tr>
							<tr><td>i_type</td><td>%c</td></tr>
							<tr><td>i_perm</td><td>%d</td></tr>
							<tr><td colspan="2">BLOQUES DIRECTOS</td></tr>
					`, i, i, inode.I_uid, inode.I_gid, inode.I_size, atime, ctime, mtime, rune(inode.I_type[0]), inode.I_perm)

		// Agregar los bloques directos a la tabla hasta el índice 11
		for j, block := range inode.I_block {
			if j > 11 {
				break
			}
			dotContent += fmt.Sprintf("<tr><td>%d</td><td>%d</td></tr>", j+1, block)
		}

		// Agregar los bloques indirectos a la tabla
		dotContent += fmt.Sprintf(`
							<tr><td colspan="2">BLOQUE INDIRECTO</td></tr>
							<tr><td>%d</td><td>%d</td></tr>
							<tr><td colspan="2">BLOQUE INDIRECTO DOBLE</td></tr>
							<tr><td>%d</td><td>%d</td></tr>
							<tr><td colspan="2">BLOQUE INDIRECTO TRIPLE</td></tr>
							<tr><td>%d</td><td>%d</td></tr>
					</table>>];
			`, 13, inode.I_block[12], 14, inode.I_block[13], 15, inode.I_block[14])

		// Agregar enlace al siguiente inodo si no es el último
		if i < superBloque.S_inodes_count-1 {
			dotContent += fmt.Sprintf("inode%d -> inode%d;\n", i, i+1) 
		}
	}

	// Cerrar el contenido DOT
	dotContent += "}"

	// Crear el archivo DOT
	dotFileName := dirPath + fileName + ".dot"
	file, err := os.Create(dotFileName)
	if err != nil {
		fmt.Println("Error al crear el archivo .dot")
		respuesta += "Error al crear el archivo .dot\n"
		return respuesta
	}
	defer file.Close()

	// Escribir el contenido DOT en el archivo
	_, err = file.WriteString(dotContent)
	if err != nil {
		fmt.Println("Error al escribir el archivo .dot")
		respuesta += "Error al escribir el archivo .dot\n"
		return respuesta
	}
	fmt.Println("Archivo .dot creado")

	// Quitar el punto a la extension
	extension := path.Ext(pathValor)[1:]

	// Crear el reporte
	cmd := exec.Command("dot", "-T", extension, dotFileName, "-o", pathValor)
	fmt.Println("dot -T" + extension + " " + dotFileName + " -o " + pathValor)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al crear el reporte con extension")
		respuesta += "Error al crear el reporte con extension\n"
		return respuesta
	}

	return "Reporte Inode creado con exito\n"
}

// REPORTE BmInode
func ReporteBmInode(id string, pathValor string) string {
	var respuesta string

	fileName := path.Base(pathValor)
	dirPath := strings.TrimSuffix(pathValor, fileName)
	fmt.Println("Nombre del archivo: " + fileName)
	fmt.Println("Ruta del archivo: " + dirPath)
	// Crear el directorio si no existe
	err := os.MkdirAll(dirPath, 0664)
	if err != nil {
		respuesta += "Error al crear el directorio\n"
		fmt.Println("Error al crear el directorio")
		return respuesta
	}

	// Buscar la particion montada
	indice := VerificarParticionMontada(id)
	if indice == -1 {
		respuesta += "La particion no esta montada"
		return respuesta
	}

	MountActual := particionesMontadas[indice]

	// Abrir el archivo
	archivo, err := os.OpenFile(MountActual.Path, os.O_RDWR, 0664)
	if err != nil {
		respuesta += "Error al abrir el archivo\n"
		fmt.Println("Error al abrir el archivo")
		return respuesta
	}
	defer archivo.Close()

	// Leer el superbloque
	superBloque := NewSuperBlock()
	archivo.Seek(int64(MountActual.Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &superBloque)
	if err != nil {
		respuesta += "Error al leer el superbloque\n"
		fmt.Println("Error al leer el superbloque")
		return respuesta
	}

	// Calcular el número total de inodos
	totalInodes := superBloque.S_inodes_count

	// Obtener el contenido del bitmap de inodos
	var bitmapContent strings.Builder

	for i := int32(0); i < totalInodes; i++ {
		// Establecer el puntero
		_, err := archivo.Seek(int64(superBloque.S_bm_inode_start+i), 0)
		if err != nil {
			return fmt.Errorf("error al establecer el puntero en el archivo: %v", err).Error()
		}

		// Leer un byte (carácter '0' o '1')
		char := make([]byte, 1)
		_, err = archivo.Read(char)
		if err != nil {
			return fmt.Errorf("error al leer el byte del archivo: %v", err).Error()
		}

		// Convertir el byte leído a '0' o '1'
		if char[0] == 0 {
			bitmapContent.WriteByte('0')
		} else {
			bitmapContent.WriteByte('1')
		}

		// Agregar un carácter de nueva línea cada 20 caracteres (20 inodos)
		if (i+1)%20 == 0 {
			bitmapContent.WriteString("\n")
		}
	}

	// Crear el archivo TXT
	txtFile, err := os.Create(pathValor)
	if err != nil {
		return fmt.Errorf("error al crear el archivo TXT: %v", err).Error()
	}
	defer txtFile.Close()

	// Escribir el contenido del bitmap en el archivo TXT
	_, err = txtFile.WriteString(bitmapContent.String())
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo TXT: %v", err).Error()
	}

	fmt.Println("Archivo del bitmap de inodos generado:", pathValor)
	return "Reporte BmInode creado con exito\n"
}
