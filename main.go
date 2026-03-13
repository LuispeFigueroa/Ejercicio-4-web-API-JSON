package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	case http.MethodDelete:
		handleDeleteFighter(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handler para manejtar todos las request de tipo GET
func handleGetFighter(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("id")
	nameParam := query.Get("name")
	countryParam := query.Get("country")
	recordParam := query.Get("record")
	specialtyParam := query.Get("specialty")
	heightParam := query.Get("height")

	var results []Fighter

	var id int
	var err error

	if idParam != "" {
		id, err = strconv.Atoi(idParam)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid ID")
			return
		}
	}
	// recorremos todo el slice de luchadores
	for _, fighter := range fighters {
		// si se envio un ID, pero no es el mismo que este luchador, continuamos al siguiente. Esto se hace con cada paramtro que podria enviarse.
		if idParam != "" && fighter.ID != id {
			continue
		}
		//uso el EqualFold para comparar los strings sin que sea case sensitive
		if nameParam != "" && !strings.EqualFold(fighter.Name, nameParam) {
			continue
		}
		if countryParam != "" && !strings.EqualFold(fighter.Country, countryParam) {
			continue
		}
		if specialtyParam != "" && !strings.EqualFold(fighter.Specialty, specialtyParam) {
			continue
		}
		if recordParam != "" && !strings.EqualFold(fighter.Record, recordParam) {
			continue
		}
		if heightParam != "" && !strings.EqualFold(fighter.Height, heightParam) {
			continue
		}
		//cada vez que se encuentra un luchador que cumple con los parametros enviados, se agrega a la lista de resultados
		results = append(results, fighter)
	}
	// si no se encuentra ningun luchador que cumpla con los parametros enviados, se devuelve un error indicando que no se encontraron luchadores
	if len(results) == 0 {
		writeError(w, http.StatusNotFound, "Fighter not found")
		return
	}
	writeJSON(w, http.StatusOK, results)
}

func handlePostFighter(w http.ResponseWriter, r *http.Request) {
	var newFighter Fighter
	err := json.NewDecoder(r.Body).Decode(&newFighter)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	//si alguno de los campos no esta completo, se devuelve un error y se solicita que llene todos los campos
	if newFighter.Name == "" || newFighter.Country == "" || newFighter.Record == "" || newFighter.Specialty == "" || newFighter.Height == "" {
		writeError(w, http.StatusBadRequest, "Missing required fields")
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
		writeError(w, http.StatusBadRequest, "ID is required")
		return
	}
	//se comprueba que la ID sea un numero valido y que exista en la lista de luchadores, si no se encuentra se devuelve un error indicando que el ID es invalido
	id, err := strconv.Atoi(idParam)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	//se decodifica el cuerpo de la solicitud para obtener los campos que se desean actualizar, si el cuerpo de la solicitud no es un JSON valido se devuelve un error indicando que el cuerpo de la solicitud es invalido
	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid Request Body")
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
	writeError(w, http.StatusNotFound, "Fighter not found")
}

func handleDeleteFighter(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idParam := query.Get("id")
	if idParam == "" {
		writeError(w, http.StatusBadRequest, "ID is required")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	for i, fighter := range fighters {
		if fighter.ID == id {
			fighters = append(fighters[:i], fighters[i+1:]...)
			saveFighters()
			writeJSON(w, http.StatusOK, Message{Message: "Fighter deleted successfully"})
			return
		}
	}
	writeError(w, http.StatusNotFound, "Fighter not found")
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		//aqui hago el log del error en caso de que ocurra un error al codificar la respuesta JSON, pero no se devuelve un error al cliente porque ya se ha establecido el encabezado y el código de estado
		log.Println("Error encoding JSON response: ", err)
	}
}

// funcion para regresar todos los errores como JSON
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
