package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Fighter struct {
	ID        int    `json:"id"`
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

	err = json.Unmarshal(file, &fighters)
	if err != nil {
		log.Fatal("Error al parsear el archivo JSON: ", err)
	}
}

// ping para ver si el server esta funcionando
func pingHandler(w http.ResponseWriter, r *http.Request) {
	response := Message{Message: "pong"}
	writeJSON(w, http.StatusOK, response)
}

func fightersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	idParam := query.Get("id")
	if idParam == "" {
		writeJSON(w, http.StatusOK, fighters)
		return
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	for _, fighter := range fighters {
		if fighter.ID == id {
			writeJSON(w, http.StatusOK, fighter)
			return
		}
	}
	http.Error(w, "Fighter not found", http.StatusNotFound)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
