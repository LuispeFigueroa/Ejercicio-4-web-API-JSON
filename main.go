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

	log.Println("Servidor iniciado en http://localhost:24087")
	log.Fatal(http.ListenAndServe(":24087", nil))

}

// funcion para cargar el archivo JSON con los lucadores usando os
func loadFighters() {
	file, err := os.ReadFile("data/fighters.json")

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
	switch r.Method {
	case http.MethodGet:
		handleGetFighter(w, r)
	case http.MethodPost:
		handlePostFighter(w, r)
	case http.MethodPatch:
		handlePatchFighter(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// handler para manejtar todos las request de tipo GET
func handleGetFighter(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("id")

	if idParam == "" {
		writeJSON(w, http.StatusOK, fighters)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
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

func handlePostFighter(w http.ResponseWriter, r *http.Request) {
	var newFighter Fighter
	err := json.NewDecoder(r.Body).Decode(&newFighter)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//si alguno de los campos no esta completo, se devuelve un error y se solicita que llene todos los campos
	if newFighter.Name == "" || newFighter.Country == "" || newFighter.Record == "" || newFighter.Specialty == "" || newFighter.Height == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	//se genera un nuevo ID para el nuevo luchador y se agrega a la lista de luchadores
	newFighter.ID = generateNextID()
	fighters = append(fighters, newFighter)
	saveFighters()

	writeJSON(w, http.StatusCreated, newFighter)
}

// funcion para generar el id del nuevo luchador
// se recorre la lista  de luchadores hasta encontrar el ID mas alto y se le suma 1 para generar el nuevo ID
func generateNextID() int {
	maxID := 0
	for _, fighter := range fighters {
		if fighter.ID > maxID {
			maxID = fighter.ID
		}
	}

	return maxID + 1
}

// funcion para guardar los luchadores en el archivo JSON
func saveFighters() {
	data, err := json.MarshalIndent(fighters, "", " ")
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	err = os.WriteFile("data/fighters.json", data, 0644)
	if err != nil {
		log.Println("Error writing JSON file:", err)
	}
}

func handlePatchFighter(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("id")
	//si no se proporciona un ID, se devuelve un error indicando que el ID es requerido
	if idParam == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	//se comprueba que la ID sea un numero valido y que exista en la lista de luchadores, si no se encuentra se devuelve un error indicando que el ID es invalido
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	//se decodifica el cuerpo de la solicitud para obtener los campos que se desean actualizar, si el cuerpo de la solicitud no es un JSON valido se devuelve un error indicando que el cuerpo de la solicitud es invalido
	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// se recorren todos los luchadores hasta encontrar el luchador con la ID proporcionada,
	// si se encuentra se actualizan los campos que se desean actualizar y se guarda la lista de luchadores en el archivo JSON,
	// si no se encuentra se devuelve un error indicando que el luchador no fue encontrado
	for i, fighter := range fighters {
		if fighter.ID == id {

			if name, ok := updates["name"].(string); ok {
				fighters[i].Name = name
			}

			if country, ok := updates["country"].(string); ok {
				fighters[i].Country = country
			}

			if record, ok := updates["record"].(string); ok {
				fighters[i].Record = record
			}

			if specialty, ok := updates["specialty"].(string); ok {
				fighters[i].Specialty = specialty

			}

			if height, ok := updates["height"].(string); ok {
				fighters[i].Height = height
			}

			saveFighters()
			writeJSON(w, http.StatusOK, fighters[i])
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
