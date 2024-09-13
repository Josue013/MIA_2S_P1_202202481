package filesystem

import (
	"fmt"
	"strconv"
	"strings"
)

// DividirComando recibe un comando y lo divide en un arreglo de strings
func DividirComando(comando string) string {
	var respuesta string
	//se divide el comando en un arreglo de strings por el enter
	comandos := strings.Split(comando, "\n")
	//se recorre el arreglo de strings
	for i := 0; i < len(comandos); i++ {
		//imprime el comando
		fmt.Println("Comando: " + comandos[i])
		//se analiza el comando
		respuesta += AnalizarComando(comandos[i])
	}
	return respuesta
}

// AnalizarComando recibe un comando y lo analiza
func AnalizarComando(comando string) string {
	var respuesta string
	//se divide el comando en un arreglo de strings por el espacio
	comandoSeparado := strings.Split(comando, " ")
	//Si encuentra el # en la posicion 0, es un comentario
	if strings.Contains(comandoSeparado[0], "#") {
		//imprime el comentario
		fmt.Println("Comentario: ")
		//Eliminiar el #
		comandoSeparado[0] = strings.Replace(comandoSeparado[0], "#", "", -1)
		respuesta += "Comentario: "
		//se recorre el arreglo de strings
		for i := 0; i < len(comandoSeparado); i++ {
			respuesta += comandoSeparado[i] + " "
		}
		respuesta += "\n"
		fmt.Println(respuesta)
	} else {
		//Si no es un comentario, entonces es un comando
		//Iterar sobre el comando
		for _, valor := range comandoSeparado {
			//el primer valor del comando lo pasamos a minusculas
			valor = strings.ToLower(valor)
			//Si el valor es mkdisk, entonces es un comando para crear un disco
			if valor == "mkdisk" {
				fmt.Println("====== Comando mkdisk ======")
				respuesta += "Ejecutando mkdisk\n"
				//Analizar el comando mkdisk
				respuesta += AnalizarMkdisk(&comandoSeparado) + "\n"
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString)
				return respuesta
			} else if valor == "rmdisk" {
				fmt.Println("====== Comando rmdisk ======")
				respuesta += "Ejecutando rmdisk\n"
				//Analizar el comando rmdisk
				respuesta += AnalizarRmdisk(&comandoSeparado)
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString) + "\n"
				return respuesta
			} else if valor == "fdisk" {
				fmt.Println("====== Comando fdisk ======")
				respuesta += "Ejecutando fdisk\n"
				//Analizar el comando fdisk
				respuesta += AnalizarFdisk(&comandoSeparado) + "\n"
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString) + "\n"
				return respuesta
			} else if valor == "mount" {
				fmt.Println("====== Comando mount ======")
				respuesta += "Ejecutando mount\n"
				//Analizar Comando Mount
				respuesta += analizarMount(&comandoSeparado) + "\n"
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString) + "\n"
				return respuesta
			} else if valor == "mkfs" {
				fmt.Println("====== Comando mkfs ======")
				respuesta += "Ejecutando mkfs\n"
				//Analizar Comando Mkfs
				respuesta += analizarMkfs(&comandoSeparado) + "\n"
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString)
				return respuesta
			} else if valor == "login" {
				fmt.Println("====== Comando login ======")
				respuesta += "Ejecutando login\n"
				//Analizar Comando Login
				respuesta += analizarLogin(&comandoSeparado) + "\n"
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString) + "\n"
				return respuesta
			} else if valor == "logout" {
				fmt.Println("====== Comando logout ======")
				respuesta += "Ejecutando comando logout\n"
				//Analizar el comando logout
				respuesta += analizarLogout(&comandoSeparado) + "\n"
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString) + "\n"
				return respuesta //Terminar la funcion
			} else if valor == "rep" {
				fmt.Println("====== Comando rep ======")
				respuesta += "Ejecutando comando rep\n"
				//Analizar el comando rep
				respuesta += AnalizarRep(&comandoSeparado) + "\n"
				//Pasar a string el comando separado
				comandoSeparadoString := strings.Join(comandoSeparado, " ")
				respuesta += AnalizarComando(comandoSeparadoString) + "\n"
				return respuesta

			} else if valor == "\n" {
				continue
			} else if valor == "\r" {
				continue
			} else if valor == "\t" {
				continue
			} else if valor == "" {
				continue
			} else {
				//Si no es ningun comando, entonces es un error
				fmt.Println("Error: Comando no reconocido")
				respuesta += "Error: Comando no reconocido\n"
			}
		}
	}

	return respuesta
}

// AnalizarMkdisk recibe un comando mkdisk y lo analiza
func AnalizarMkdisk(comando *[]string) string {
	//mkdisk -unit=M -path="/home 1/mis discos/Disco3.mia"
	//0 		1     2     3"/home/mis     4         5discos/Disco3.mia"
	var respuesta string
	//Eliminar el primer valor del comando
	*comando = (*comando)[1:]
	//-size=5 -unit=M -path="/home/mis discos/Disco3.mia"
	//Booleanos para saber si se encontro el size, unit, fit, path
	var size, unit, path, fit bool
	//Variables para guardar el valor del size, unit, fit, path
	var valorSize, valorUnit, valorFit, valorPath string
	//Iterar sobre el comando
	valorFit = "f"
	valorUnit = "m"
	for _, valor := range *comando {
		bandera := obtenerBandera(valor)
		banderaValor := obtenerValor(valor)
		if bandera == "-size" {
			size = true
			valorSize = banderaValor
			*comando = (*comando)[1:]
		} else if bandera == "-unit" {
			unit = true
			valorUnit = banderaValor
			valorUnit = strings.ToLower(valorUnit)
			*comando = (*comando)[1:]
		} else if bandera == "-fit" {
			fit = true
			valorFit = banderaValor
			valorFit = strings.ToLower(valorFit)
			*comando = (*comando)[1:]
		} else if bandera == "-path" {
			path = true
			// Verificar si el path tiene comillas
			// -path="/home 1/mis discos/Disco3.mia"
			if strings.Contains(banderaValor, "\"") {
					// Eliminar las comillas del inicio
					banderaValor = strings.Replace(banderaValor, "\"", "", -1)
					// Inicializar valorPath con el valor actual
					valorPath = banderaValor
					// Eliminar el primer valor del comandoSeparado
					*comando = (*comando)[1:]
					// Iterar sobre el comando
					Contador := 0
					for _, valor := range *comando {
							// Si el valor contiene comillas
							if strings.Contains(valor, "\"") {
									// Eliminar las comillas del final
									valor = strings.Replace(valor, "\"", "", -1)
									// Agregar el valor al path
									valorPath += " " + valor
									Contador++
									break
							} else {
									// Agregar el valor al path
									valorPath += " " + valor
									Contador++
							}
					}
					// Eliminar los valores del comando
					*comando = (*comando)[Contador:]
			} else {
					valorPath = banderaValor
					*comando = (*comando)[1:]
			}
	} else {
			fmt.Println("Error: Parametro no reconocida")
			respuesta += "Error: Parametro no reconocida\n"
	}

	}
	if !size {
		fmt.Println("Error: Falta el parametro size")
		respuesta += "Error: Falta el parametro size\n"
		return respuesta
	} else if !path {
		fmt.Println("Error: Falta el parametro path")
		respuesta += "Error: Falta el parametro path\n"
		return respuesta
	} else {
		if fit {
			if valorFit != "bf" && valorFit != "ff" && valorFit != "wf" {
				fmt.Println("Error: Fit no reconocido")
				respuesta += "Error: Fit no reconocido\n"
				return respuesta
			} else {
				if valorFit == "bf" {
					valorFit = "b"
				} else if valorFit == "ff" {
					valorFit = "f"
				} else if valorFit == "wf" {
					valorFit = "w"
				}
			}
		}
		if unit {
			if valorUnit != "k" && valorUnit != "m" {
				fmt.Println("Error: Unit no reconocido")
				respuesta += "Error: Unit no reconocido\n"
				return respuesta
			}
		}
		//Pasar a int el size
		sizeInt, err := strconv.Atoi(valorSize)
		if err != nil {
			fmt.Println("Error: Size no es un numero")
			respuesta += "Error: Size no es un numero\n"
			return respuesta
		}
		//Verificar que el size sea mayor a 0
		if sizeInt <= 0 {
			fmt.Println("Error: Size debe ser mayor a 0")
			respuesta += "Error: Size debe ser mayor a 0\n"
			return respuesta
		}
		//Imprimir los valores
		fmt.Println("Size: " + valorSize)
		fmt.Println("Unit: " + valorUnit)
		fmt.Println("Fit: " + valorFit)
		fmt.Println("Path: " + valorPath)
		//Llamar a la funcion para crear el disco
		respuesta += CrearDiscos(sizeInt, valorUnit, valorFit, valorPath)
		return respuesta
	}
}

// AnalizarRmdisk recibe un comando rmdisk y lo analiza
func AnalizarRmdisk(comando *[]string) string {
	//rmdisk -path=/home/Disco1.mia
	//respuesta
	var respuesta string
	//Booleanos para saber si se encontro el path
	var path bool
	//Variables para guardar el valor del path
	var valorPath string
	//Iterar sobre el comando
	for _, valor := range *comando {
		bandera := obtenerBandera(valor)
		banderaValor := obtenerValor(valor)
		if bandera == "-path" {
			path = true
			// Verificar si el path tiene comillas
			// -path="/home 1/mis discos/Disco3.mia"
			if strings.Contains(banderaValor, "\"") {
					// Eliminar las comillas del inicio
					banderaValor = strings.Replace(banderaValor, "\"", "", -1)
					// Inicializar valorPath con el valor actual
					valorPath = banderaValor
					// Eliminar el primer valor del comandoSeparado
					*comando = (*comando)[1:]
					// Iterar sobre el comando
					Contador := 0
					for _, valor := range *comando {
							// Si el valor contiene comillas
							if strings.Contains(valor, "\"") {
									// Eliminar las comillas del final
									valor = strings.Replace(valor, "\"", "", -1)
									// Agregar el valor al path
									valorPath += " " + valor
									Contador++
									break
							} else {
									// Agregar el valor al path
									valorPath += " " + valor
									Contador++
							}
					}
					// Eliminar los valores del comando
					*comando = (*comando)[Contador:]
			} else {
					valorPath = banderaValor
					*comando = (*comando)[1:]
			}
	} else {
			fmt.Println("Error: Parametro no reconocida")
			respuesta += "Error: Parametro no reconocida\n"
	}
	}
	//Obligatorios: path
	if !path {
		fmt.Println("Error: Falta el parametro path")
		respuesta += "Error: Falta el parametro path\n"
		return respuesta
	} else {
		//Imprimir los valores
		fmt.Println("Path: " + valorPath)
		//Llamar a la funcion para eliminar el disco
		respuesta += EliminarDiscos(valorPath)
		return respuesta
	}
}

// AnalizarFdisk recibe un comando fdisk y lo analiza
func AnalizarFdisk(comando *[]string) string {
	//fdisk -Size=300 -path=/home/Disco1.mia -name=Particion1
	*comando = (*comando)[1:]
	//respuesta
	var respuesta string
	//Booleanos para saber si se encontro el size, unit, fit, path
	var size, unit, path, name, typePart, fit bool
	//Variables para guardar el valor del size, unit, fit, path
	var valorSize, valorUnit, valorFit, valorPath, valorName, valorTypePart string
	valorFit = "f"
	valorUnit = "k"
	valorTypePart = "p"
	//Iterar sobre el comando
	for _, valor := range *comando {
		//Obtener la bandera
		bandera := obtenerBandera(valor)
		//Obtener el valor de la bandera
		banderaValor := obtenerValor(valor)
		//Si la bandera es -size
		if bandera == "-size" {
			size = true
			valorSize = banderaValor
			*comando = (*comando)[1:]
		} else if bandera == "-unit" {
			unit = true
			valorUnit = banderaValor
			valorUnit = strings.ToLower(valorUnit)
			*comando = (*comando)[1:]
		} else if bandera == "-fit" {
			fit = true
			valorFit = banderaValor
			valorFit = strings.ToLower(valorFit)
			*comando = (*comando)[1:]
		} else if bandera == "-name" {
			name = true
			valorName = banderaValor
			*comando = (*comando)[1:]
		} else if bandera == "-type" {
			typePart = true
			valorTypePart = banderaValor
			valorTypePart = strings.ToLower(valorTypePart)
			*comando = (*comando)[1:]
		} else if bandera == "-path" {
			path = true
			// Verificar si el path tiene comillas
			// -path="/home 1/mis discos/Disco3.mia"
			if strings.Contains(banderaValor, "\"") {
					// Eliminar las comillas del inicio
					banderaValor = strings.Replace(banderaValor, "\"", "", -1)
					// Inicializar valorPath con el valor actual
					valorPath = banderaValor
					// Eliminar el primer valor del comandoSeparado
					*comando = (*comando)[1:]
					// Iterar sobre el comando
					Contador := 0
					for _, valor := range *comando {
							// Si el valor contiene comillas
							if strings.Contains(valor, "\"") {
									// Eliminar las comillas del final
									valor = strings.Replace(valor, "\"", "", -1)
									// Agregar el valor al path
									valorPath += " " + valor
									Contador++
									break
							} else {
									// Agregar el valor al path
									valorPath += " " + valor
									Contador++
							}
					}
					// Eliminar los valores del comando
					*comando = (*comando)[Contador:]
			} else {
					valorPath = banderaValor
					*comando = (*comando)[1:]
			}
	} else {
			fmt.Println("Error: Parametro no reconocida")
			respuesta += "Error: Parametro no reconocida\n"
	}
	}
	//Obligatorios: name, path y size
	if !name {
		fmt.Println("Error: Falta el parametro name")
		respuesta += "Error: Falta el parametro name\n"
		return respuesta
	} else if !path {
		fmt.Println("Error: Falta el parametro path")
		respuesta += "Error: Falta el parametro path\n"
		return respuesta
	} else if !size {
		fmt.Println("Error: Falta el parametro size")
		respuesta += "Error: Falta el parametro size\n"
		return respuesta
	} else {
		//Opcionales: unit, fit, type
		if fit {
			if valorFit != "bf" && valorFit != "ff" && valorFit != "wf" {
				fmt.Println("Error: Fit no reconocido")
				respuesta += "Error: Fit no reconocido\n"
				return respuesta
			} else {
				if valorFit == "bf" {
					valorFit = "b"
				} else if valorFit == "ff" {
					valorFit = "f"
				} else if valorFit == "wf" {
					valorFit = "w"
				}
			}
		}
		if unit {
			if valorUnit != "k" && valorUnit != "m" && valorUnit != "b" {
				fmt.Println("Error: Unit no reconocido")
				respuesta += "Error: Unit no reconocido\n"
				return respuesta
			}
		}
		if typePart {
			if valorTypePart != "p" && valorTypePart != "e" && valorTypePart != "l" {
				fmt.Println("Error: Type no reconocido")
				respuesta += "Error: Type no reconocido\n"
				return respuesta
			}
		}
		//Pasar a int el size
		sizeInt, err := strconv.Atoi(valorSize)
		if err != nil {
			fmt.Println("Error: Size no es un numero")
			respuesta += "Error: Size no es un numero\n"
			return respuesta
		}
		//Verificar que el size sea mayor a 0
		if sizeInt <= 0 {
			fmt.Println("Error: Size debe ser mayor a 0")
			respuesta += "Error: Size debe ser mayor a 0\n"
			return respuesta
		}
		//Imprimir los valores
		fmt.Println("Size: " + valorSize)
		fmt.Println("Unit: " + valorUnit)
		fmt.Println("Fit: " + valorFit)
		fmt.Println("Path: " + valorPath)
		fmt.Println("Name: " + valorName)
		fmt.Println("Type: " + valorTypePart)
		//Llamar a la funcion para crear la particion
		respuesta += Fdisk(sizeInt, valorUnit, valorFit, valorPath, valorName, valorTypePart)
		return respuesta
	}
}

func analizarMount(comandoSeparado *[]string) string {
	//respuesta
	var respuesta string
	//mount -driveletter=A -name=Part1 #id=A118
	*comandoSeparado = (*comandoSeparado)[1:]
	//Booleanos para verificar si se ingresaron los parametros
	var banderaPath, banderaName bool
	//Variables para almacenar los valores de los parametros
	var pathValor, nameValor string
	//Iterar sobre el comando separado
	for _, valor := range *comandoSeparado {
			bandera := obtenerBandera(valor)
			banderaValor := obtenerValor(valor)
			if bandera == "-path" {
					banderaPath = true
					// Verificar si el path tiene comillas
					// -path="/home 1/mis discos/Disco3.mia"
					if strings.Contains(banderaValor, "\"") {
							// Eliminar las comillas del inicio
							banderaValor = strings.Replace(banderaValor, "\"", "", -1)
							// Inicializar pathValor con el valor actual
							pathValor = banderaValor
							// Eliminar el primer valor del comandoSeparado
							*comandoSeparado = (*comandoSeparado)[1:]
							// Iterar sobre el comando
							Contador := 0
							for _, valor := range *comandoSeparado {
									// Si el valor contiene comillas
									if strings.Contains(valor, "\"") {
											// Eliminar las comillas del final
											valor = strings.Replace(valor, "\"", "", -1)
											// Agregar el valor al path
											pathValor += " " + valor
											Contador++
											break
									} else {
											// Agregar el valor al path
											pathValor += " " + valor
											Contador++
									}
							}
							// Eliminar los valores del comando
							*comandoSeparado = (*comandoSeparado)[Contador:]
					} else {
							pathValor = banderaValor
							*comandoSeparado = (*comandoSeparado)[1:]
					}
			} else if bandera == "-name" {
					banderaName = true
					nameValor = banderaValor
					*comandoSeparado = (*comandoSeparado)[1:]
			} else {
					fmt.Println("Parametro no reconocido: ", bandera)
					respuesta += "Parametro no reconocido: " + bandera + "\n"
			}
	}
	// Obligatorios: -path, -name
	// Verificar si se ingresaron los parametros obligatorios
	if !banderaPath {
			fmt.Println("El parametro -path es obligatorio")
			respuesta += "El parametro -path es obligatorio\n"
			return respuesta
	} else if !banderaName {
			fmt.Println("El parametro -name es obligatorio")
			respuesta += "El parametro -name es obligatorio\n"
			return respuesta
	} else {
			// Imprimir los valores de los parametros
			fmt.Println("Path: ", pathValor)
			fmt.Println("Name: ", nameValor)
			// Llamar a la funcion para montar la particion
			respuesta += MountPartition(pathValor, nameValor)
			return respuesta
	}
}

func analizarMkfs(comandoSeparado *[]string) string {
	// mkfs -type=full -id=B145 -fs=3fs
	//respuesta
	var respuesta string
	*comandoSeparado = (*comandoSeparado)[1:]
	//Booleanos para verificar si se ingresaron los parametros
	var banderaType, banderaId, banderaFs bool
	//Variables para almacenar los valores de los parametros
	var typeValor, idValor, fsValor string
	typeValor = "full"
	fsValor = "2fs"
	//Iterar sobre el comando separado
	for _, valor := range *comandoSeparado {
		bandera := obtenerBandera(valor)
		banderaValor := obtenerValor(valor)
		if bandera == "-type" {
			banderaType = true
			typeValor = banderaValor
			typeValor = strings.ToLower(typeValor)
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-id" {
			banderaId = true
			idValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-fs" {
			banderaFs = true
			fsValor = banderaValor
			fsValor = strings.ToLower(fsValor)
			*comandoSeparado = (*comandoSeparado)[1:]
		} else {
			fmt.Println("Parametro no reconocido: ", bandera)
			respuesta += "Parametro no reconocido: " + bandera + "\n"
		}
	}
	//Obligatorios: -id
	//Verificar si se ingresaron los parametros obligatorios
	if !banderaId {
		fmt.Println("El parametro -id es obligatorio")
		respuesta += "El parametro -id es obligatorio\n"
		return respuesta
	} else {
		//Verificar si se ingresaron los parametros aceptados
		if banderaType {
			if typeValor != "full" {
				fmt.Println("El valor del parametro -type no es valido")
				respuesta += "El valor del parametro -type no es valido\n"
				return respuesta
			}
		}
		if banderaFs {
			if fsValor != "2fs" && fsValor != "3fs" {
				fmt.Println("El valor del parametro -fs no es valido")
				respuesta += "El valor del parametro -fs no es valido\n"
				return respuesta
			}
		}
		//Imprimir los valores de los parametros
		fmt.Println("Type: ", typeValor)
		fmt.Println("Id: ", idValor)
		fmt.Println("Fs: ", fsValor)
		//Llamar a la funcion para formatear la particion
		respuesta += Mkfs(typeValor, idValor, fsValor)
		return respuesta
	}
}

func analizarLogin(comandoSeparado *[]string) string {
	var respuesta string
	//Eliminar el primer valor del comandoSeparado
	*comandoSeparado = (*comandoSeparado)[1:]
	//Booleanos para verificar si se ingreso el parametro obligatorio
	var user, pwd, id bool
	//Variables para guardar los valores de los parametros
	var userValor, pwdValor, idValor string
	//Recorrer el comandoSeparado
	for _, valor := range *comandoSeparado {
		bandera := obtenerBandera(valor)
		banderaValor := obtenerValor(valor)
		if bandera == "-user" {
			if banderaValor[0] == '"' {
				banderaValor = banderaValor[1:]
				//Eliminar el primer valor del comandoSeparado
				Contador := 0
				*comandoSeparado = (*comandoSeparado)[1:]
				for _, valor := range *comandoSeparado {
					banderaValor += " " + valor
					Contador++
					//Eliminar \r y \n
					if strings.Contains(valor, "\r") {
						valor = strings.Replace(valor, "\r", "", -1)
					}
					if strings.Contains(valor, "\n") {
						valor = strings.Replace(valor, "\n", "", -1)
					}
					if strings.Contains(valor, "\"") {
						break
					}
				}
				//Eliminar los valores que ya se analizaron
				*comandoSeparado = (*comandoSeparado)[Contador-1:]
				//Eliminar las comillas del final
				banderaValor = strings.Replace(banderaValor, "\"", "", -1)
			}
			user = true
			userValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]

		} else if bandera == "-pass" {
			pwd = true
			pwdValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else if bandera == "-id" {
			id = true
			idValor = banderaValor
			*comandoSeparado = (*comandoSeparado)[1:]
		} else {
			fmt.Println("Parametro no reconocido")
			respuesta += "Parametro no reconocido\n"
		}
	}
	//Verificar si se ingreso el parametro obligatorio
	if !user {
		fmt.Println("No se ingreso el parametro obligatorio user")
		respuesta += "No se ingreso el parametro obligatorio user\n"
		return respuesta
	} else if !pwd {
		fmt.Println("No se ingreso el parametro obligatorio pass")
		respuesta += "No se ingreso el parametro obligatorio \n"
		return respuesta
	} else if !id {
		fmt.Println("No se ingreso el parametro obligatorio id")
		respuesta += "No se ingreso el parametro obligatorio id\n"
		return respuesta
	} else {
		//Imprimir los valores y ejecutar el comando
		fmt.Println("user: ", userValor)
		fmt.Println("pass: ", pwdValor)
		fmt.Println("id: ", idValor)
		//Ejecutar el comando
		respuesta += Login(userValor, pwdValor, idValor)
		return respuesta
	}
}

func AnalizarRep(comandoSeparado *[]string) string {
	var respuesta string
	//Eliminar el primer valor del comandoSeparado
	*comandoSeparado = (*comandoSeparado)[1:]
	//Booleanos para verificar si se ingreso el parametro obligatorio
	var id, name, path_file_ls, path bool
	//Variables para guardar los valores de los parametros
	var idValor, nameValor, path_file_fsValor, pathValor string
	//Recorrer el comandoSeparado
	for _, valor := range *comandoSeparado {
			bandera := obtenerBandera(valor)
			banderaValor := obtenerValor(valor)
			if bandera == "-id" {
					id = true
					idValor = banderaValor
					*comandoSeparado = (*comandoSeparado)[1:]
			} else if bandera == "-name" {
					name = true
					nameValor = banderaValor
					//Pasar el valor a minusculas
					nameValor = strings.ToLower(nameValor)
					*comandoSeparado = (*comandoSeparado)[1:]
			} else if bandera == "-path_file_ls" {
					// Verificar si el path tiene comillas
					// -path_file_ls="/home 1/mis discos/Disco3.mia"
					if strings.Contains(banderaValor, "\"") {
							// Eliminar las comillas del inicio
							banderaValor = strings.Replace(banderaValor, "\"", "", -1)
							// Inicializar path_file_fsValor con el valor actual
							path_file_fsValor = banderaValor
							// Eliminar el primer valor del comandoSeparado
							*comandoSeparado = (*comandoSeparado)[1:]
							// Iterar sobre el comando
							Contador := 0
							for _, valor := range *comandoSeparado {
									// Si el valor contiene comillas
									if strings.Contains(valor, "\"") {
											// Eliminar las comillas del final
											valor = strings.Replace(valor, "\"", "", -1)
											// Agregar el valor al path
											path_file_fsValor += " " + valor
											Contador++
											break
									} else {
											// Agregar el valor al path
											path_file_fsValor += " " + valor
											Contador++
									}
							}
							// Eliminar los valores del comando
							*comandoSeparado = (*comandoSeparado)[Contador:]
					} else {
							path_file_fsValor = banderaValor
							*comandoSeparado = (*comandoSeparado)[1:]
					}
					path_file_ls = true
			} else if bandera == "-path" {
					// Verificar si el path tiene comillas
					// -path="/home 1/mis discos/Disco3.mia"
					if strings.Contains(banderaValor, "\"") {
							// Eliminar las comillas del inicio
							banderaValor = strings.Replace(banderaValor, "\"", "", -1)
							// Inicializar pathValor con el valor actual
							pathValor = banderaValor
							// Eliminar el primer valor del comandoSeparado
							*comandoSeparado = (*comandoSeparado)[1:]
							// Iterar sobre el comando
							Contador := 0
							for _, valor := range *comandoSeparado {
									// Si el valor contiene comillas
									if strings.Contains(valor, "\"") {
											// Eliminar las comillas del final
											valor = strings.Replace(valor, "\"", "", -1)
											// Agregar el valor al path
											pathValor += " " + valor
											Contador++
											break
									} else {
											// Agregar el valor al path
											pathValor += " " + valor
											Contador++
									}
							}
							// Eliminar los valores del comando
							*comandoSeparado = (*comandoSeparado)[Contador:]
					} else {
							pathValor = banderaValor
							*comandoSeparado = (*comandoSeparado)[1:]
					}
					path = true
			} else {
					fmt.Println("Parametro no reconocido")
					respuesta += "Parametro no reconocido\n"
			}
	}
	//Verificar si se ingreso el parametro obligatorio
	if !id {
			fmt.Println("No se ingreso el parametro obligatorio id")
			respuesta += "No se ingreso el parametro obligatorio id\n"
			return respuesta
	} else if !name {
			fmt.Println("No se ingreso el parametro obligatorio name")
			respuesta += "No se ingreso el parametro obligatorio name\n"
			return respuesta
	} else if !path {
			fmt.Println("No se ingreso el parametro obligatorio path")
			respuesta += "No se ingreso el parametro obligatorio path\n"
			return respuesta
	} else {
			if nameValor == "disk" {
					//Imprimir los valores y ejecutar el comando
					fmt.Println("id: ", idValor)
					fmt.Println("name: ", nameValor)
					fmt.Println("path: ", pathValor)
					//Ejecutar el comando
					respuesta += ReporteDisk(idValor, pathValor)
					return respuesta
			} else if nameValor == "file" {
					if !path_file_ls {
							fmt.Println("No se ingreso el parametro obligatorio path_file_ls")
							respuesta += "No se ingreso el parametro obligatorio path_file_ls\n"
							return respuesta
					} else {
							//Imprimir los valores y ejecutar el comando
							fmt.Println("id: ", idValor)
							fmt.Println("name: ", nameValor)
							fmt.Println("path: ", pathValor)
							fmt.Println("path_file_ls: ", path_file_fsValor)
							//Ejecutar el comando
							respuesta += ReporteFile(idValor, pathValor, path_file_fsValor)
							return respuesta
					}
			} else if nameValor == "sb" {
					//Imprimir los valores y ejecutar el comando
					fmt.Println("id: ", idValor)
					fmt.Println("name: ", nameValor)
					fmt.Println("path: ", pathValor)
					//Ejecutar el comando
					respuesta += ReporteSB(idValor, pathValor)
					return respuesta
			} else if nameValor == "mbr" {
					//Imprimir los valores y ejecutar el comando
					fmt.Println("id: ", idValor)
					fmt.Println("name: ", nameValor)
					fmt.Println("path: ", pathValor)
					//Ejecutar el comando
					respuesta += ReporteMBR(idValor, pathValor)
					return respuesta
			} else if nameValor == "inode" {
					//Imprimir los valores y ejecutar el comando
					fmt.Println("id: ", idValor)
					fmt.Println("name: ", nameValor)
					fmt.Println("path: ", pathValor)
					//Ejecutar el comando
					respuesta += ReporteInode(idValor, pathValor)
					return respuesta
			} else if nameValor == "bm_inode" {
					//Imprimir los valores y ejecutar el comando
					fmt.Println("id: ", idValor)
					fmt.Println("name: ", nameValor)
					fmt.Println("path: ", pathValor)
					//Ejecutar el comando
					respuesta += ReporteBmInode(idValor, pathValor)
					return respuesta
			} else {
					fmt.Println("Los valores de name deben ser: disk, sb, tree o file")
					respuesta += "Los valores de name deben ser: disk, sb, tree o file\n"
					return respuesta
			}
			/*else if nameValor == "tree" {
					//Imprimir los valores y ejecutar el comando
					fmt.Println("id: ", idValor)
					fmt.Println("name: ", nameValor)
					fmt.Println("path: ", pathValor)
					//Ejecutar el comando
					respuesta += ReporteTree(idValor, pathValor)
					return respuesta
			}*/
	}
}

func analizarLogout(comandoSeparado *[]string) string {
	var respuesta string
	//Eliminar el primer valor del comandoSeparado
	*comandoSeparado = (*comandoSeparado)[1:]
	//Llamar a la funcion logout
	respuesta += Logout()
	return respuesta
}

func obtenerBandera(bandera string) string {
	//-size
	var banderaValor string
	for _, valor := range bandera {
		if valor == '=' {
			break
		}
		banderaValor += string(valor)
	}
	banderaValor = strings.ToLower(banderaValor)
	return banderaValor
}

func obtenerValor(bandera string) string {
	var banderaValor string
	var boolBandera bool
	for _, valor := range bandera {
		if boolBandera {
			banderaValor += string(valor)
		}
		if valor == '=' {
			boolBandera = true
		}
	}
	return banderaValor
}
