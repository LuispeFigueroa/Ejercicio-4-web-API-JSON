package main

import (
	"log"
	"net/http"
	"os"
)

type Fighter struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Record    string `json:"record"`
	Specialty string `json:"specialty"`
	Height    string `json:"height"`
}

type Message struct {
	Message string `json:"text"`
}

var fighters []Fighter

func main() {
	//cargar a los luchadores en el JSON
	loadFighters()

	http.HandleFunc("/api/ping", pingHandler)
	http.HandleFunc("/api/Fighters", fightersHandler)

}

// funcion para cargar el archivo JSON con los lucadores usando os
func loadFighters() {
	file, err := os.ReadFile("/data/fighters.json")

	if err != nil {
		log.Fatal("Error al leer el archivo JSON: ", err)
	}
}

// ping para ver si el server esta funcionando
func pingHandler(w http.ResponseWriter, r *http.Request) {
	response = Message{Message: "pong"}
	writeJSON(w, http.StatusOK, response)
}

func fightersHandler(w http.ResponseWriter, r *http.Request) {

}
