package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}
type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}
//first speed
func addCat(c echo.Context) error {
	cat := Cat{}
	defer c.Request().Body.Close()
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading request body for addcart %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	err = json.Unmarshal(b, &cat)
	if err != nil {
		log.Printf("Failed unmarshaling json data for addcart %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is your cat %#v", cat)
	return c.String(http.StatusOK,"we got your cat!")
}
//second speed
func addDog(c echo.Context)error{
	dog := Dog{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil{
		log.Printf("Failed reading request body for addDog %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is your dog %#v", dog)
	return c.String(http.StatusOK,"we got your dog!")
}
//third speed
func addFiker(c echo.Context)error{
	fiker := Fiker{}
	err := c.Bind(&fiker)
	if err != nil{
		log.Printf("Failed reading request body for addFiker %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is your fiker %#v", fiker)
	return c.String(http.StatusOK,"we got your fiker!")
}
func main() {
	r := mux.NewRouter()
	movies = append(movies, Movie{ID: "1", Isbn: "John Coltrane", Title: "Blue Train", Director: &Director{FirstName: "Abemelek", LastName: "Daniel"}})
	movies = append(movies, Movie{ID: "2", Isbn: "Gerry Mulligan", Title: "Jeru", Director: &Director{FirstName: "Serawit", LastName: "Sissay"}})
	movies = append(movies, Movie{ID: "3", Isbn: "Sarah Vaughan", Title: "Sarah Vaughan and Clifford Brown", Director: &Director{FirstName: "Love", LastName: "God"}})
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovies).Methods("DELETE")

	fmt.Printf("Starting server at post 4000\n")
	log.Fatal(http.ListenAndServe(":4000", r))
}
