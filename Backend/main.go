package main

import (
	"MIA_2S_P1_202202481/Backend/filesystem"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Comando struct {
	//json:"peticion"
	Comando string `json:"peticion"`
}

type User struct {
	Carnet int
	Nombre string
}

type Respuesta struct {
	ResponseBack string `json:"respuesta"`
	Error        bool   `json:"error"`
}

type allTasks []User

var tasks = allTasks{
	{
		Carnet: 202202481,
		Nombre: "Josué Nabí Hurtarte Pinto",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func leerComando(w http.ResponseWriter, r *http.Request) {
	var newComando Comando
	var newRespuesta Respuesta
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Inserte un comando valido")
		newRespuesta.ResponseBack = "Inserte un comando valido"
	}
	json.Unmarshal(reqBody, &newComando)
	newRespuesta.ResponseBack = filesystem.DividirComando(newComando.Comando)
	fmt.Println(newRespuesta.ResponseBack)
	//Agregar la respuesta a la peticion
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRespuesta)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenido a mi servidor")
}

func main() {
	//Rutas
	fmt.Println("Josué Nabí Hurtarte Pinto - 202202481 - MIA PROYECTO 1")
	router := mux.NewRouter().StrictSlash(true)
	//Endpoints
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/Josue", getTasks).Methods("GET")
	router.HandleFunc("/command", leerComando).Methods("POST")

	headres := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	//Servidor
	fmt.Println("Servidor corriendo en el puerto http://localhost:3000")
	http.ListenAndServe(":3000", handlers.CORS(headres, methods, origins)(router))
}
