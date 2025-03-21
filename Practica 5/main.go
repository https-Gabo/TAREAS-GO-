package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Usuario struct {
	ID     int    `json:"id"`
	Nombre string `json:"name"`
	Email  string `json:"email"`
}

var usuarios []Usuario

// Usuarios

func Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("test-header", "header")
		json.NewEncoder(w).Encode(usuarios)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error al leer el body", http.StatusBadRequest)
		}
		var user Usuario
		err = json.Unmarshal(body, &user)
		if err != nil {
			http.Error(w, "Error parseando el JSON", http.StatusBadRequest)
		}
		user.ID = len(usuarios) + 1
		usuarios = append(usuarios, user)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("test-header", "header")
		json.NewEncoder(w).Encode(user)
	case http.MethodDelete:
	case http.MethodPut:
		fmt.Println("No implementado")
	default:
		http.Error(w, "Método no permitido", 405)
	}
}

func Ping(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "pong")
	default:
		http.Error(w, "Método no permitido", 405)
	}
}
func Index(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("./public/index.html")
	if err != nil {
		fmt.Fprintln(w, "error leyendo el html")
		return
	}
	fmt.Fprintln(w, string(content))
}
func main() {
	usuarios = append(usuarios, Usuario{
		ID:     1,
		Nombre: "Alfredo",
		Email:  "Alfredo@mail.com",
	})
	http.HandleFunc("/ping", Ping)
	http.HandleFunc("/v1/users", Users)
	http.HandleFunc("/", Index)

	fmt.Println("Servidor escuchando en el puerto 3000")
	http.ListenAndServe(":3000", nil)
}
