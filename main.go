package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// estructura de las peliculas
type Movies struct {
	Nombre   string `json:"nombre"`
	Isbn     string `json:"isbn"`
	Director string `json:"director"`
}

var movies []Movies

// funciones para manejar las peticiones
func GetAllMovies(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/movies" {
		http.Error(w, "Error bad request", http.StatusNotFound)
	}

	if r.Method != "GET" {
		http.Error(w, "Error Method not supported", http.StatusNotFound)
	}

	fmt.Fprintf(w, "%s\n", r.URL.Path)
	fmt.Fprintf(w, "%s\n", r.Method)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Error bad Request", http.StatusNotFound)
	}

	var movie Movies
	json.NewDecoder(r.Body).Decode(&movie)

	movies = append(movies, movie)

	fmt.Fprintf(w, "%s\n", r.URL.Path)
	fmt.Fprintf(w, "%s\n", r.Method)
	fmt.Fprintf(w, "%v\n", movie)
	fmt.Fprintf(w, "peliculas: %v\n", movies)
}

func Index(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	method := r.Method

	if path != "/home" {
		http.Error(w, "Error bad request", http.StatusNotFound)
	}

	if method != "GET" {
		http.Error(w, "Error Method not supported.", http.StatusNotFound)
	}

	fmt.Fprintf(w, "PATH: %s\n", path)
	fmt.Fprintf(w, "METHOD: %s\n", method)

	w.Write([]byte(path))
}

func main() {

	// create server
	fmt.Println("Listening on port 3000")

	//endPoints
	http.HandleFunc("/home", Index)
	// http.HandleFunc("/movies", GetAllMovies)
	http.HandleFunc("/movies", CreateMovie)

	err := http.ListenAndServe("127.0.0.1:3000", nil)

	if err != nil {
		log.Fatal(err)
	}

}
