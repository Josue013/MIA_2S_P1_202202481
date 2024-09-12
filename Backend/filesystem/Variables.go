package filesystem

import "fmt"

type Mount struct {
	Id        string
	Path      string
	Name      string
	Part_type [1]byte
	Start     int32
	Size      int32
}

type usuarioActual struct {
	Uid  int32
	Gid  int32
	Grp  string
	Usr  string
	Pass string
	Pid  string
}

func NuevoUsuarioActual() usuarioActual {
	var usr usuarioActual
	usr.Uid = -1
	usr.Gid = -1
	usr.Grp = ""
	usr.Usr = ""
	usr.Pass = ""
	usr.Pid = ""
	return usr
}

var particionesMontadas []Mount

var pathsParticiones []string

var Usr_sesion usuarioActual = NuevoUsuarioActual()

func VerificarParticionMontada(id string) int {
	for i := 0; i < len(particionesMontadas); i++ {
		fmt.Println(particionesMontadas[i].Id)
		if particionesMontadas[i].Id == id {
			return i
		}
	}
	return -1
}
