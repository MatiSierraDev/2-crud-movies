package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// estructura de las peliculas
type Movies struct {
	ID       string `json:"id"`
	Nombre   string `json:"nombre"`
	Release  string `json:"release"`
	Director string `json:"director"`
}

const serverPort = 3030

var movies []Movies

// funciones para manejar las peticiones
func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if r.URL.Path != "/movies" {
		http.Error(w, "Error bad request", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Error Method not supported", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(movies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.URL.Path != "/movies" {
		http.Error(w, "Bad request", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Error bad Request", http.StatusNotFound)
		return
	}

	var movie Movies

	json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100000000))

	for _, item := range movies {
		if movie.ID == item.ID {
			http.Error(w, "Bad Request", http.StatusConflict)
			return
		}
	}

	movies = append(movies, movie)

	fmt.Fprintf(w, "peliculas: %v\n", movies)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movies
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dataId := mux.Vars(r)["id"]

	for _, movie := range movies {
		if movie.ID == dataId {
			json.NewEncoder(w).Encode(movie)
			fmt.Println(movie)
			return
		}
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	if path != "/" {
		http.Error(w, "Error bad request index", http.StatusNotFound)
	}

	if method != "GET" {
		http.Error(w, "Error Method not supported.", http.StatusNotFound)
	}

	w.Write([]byte(path))
}

func main() {
	//ENRUTADOR GORILLA/MUX
	mux := mux.NewRouter()

	movies = append(movies, Movies{ID: "1", Nombre: "Dragon Ball: The Legend of Shen Long", Release: "December 20, 1986", Director: "Daisuke Nishio"})
	movies = append(movies, Movies{ID: "2", Nombre: "Dragon Ball: Sleeping Princess in Devil's Castle", Release: "July 18, 1987", Director: "Daisuke Nishio"})
	mux.HandleFunc("/", Index).Methods("GET")
	mux.HandleFunc("/movies", GetAllMovies).Methods("GET")
	mux.HandleFunc("/movies/{id}", GetMovie).Methods("GET")
	mux.HandleFunc("/movies", CreateMovie).Methods("POST")
	mux.HandleFunc("/movies/{id}", DeleteMovie).Methods("DELETE")
	mux.HandleFunc("/movies/{id}", UpdateMovie).Methods("PUT")

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: mux,
	}

	fmt.Printf("Listening on server %d ", serverPort)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
