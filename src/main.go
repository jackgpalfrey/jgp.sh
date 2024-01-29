package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "JGP.SH")
}

type HealthcheckReport struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	log.Print("HEALTHCHECK")
	w.Header().Set("Content-Type", "application/json")
	report := HealthcheckReport{
		Name: "JGP.SH",
		Version: "indev",
		Status: "OK",
	}

	err := json.NewEncoder(w).Encode(&report)
	if err != nil {
		log.Fatal(err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request){
	http.Redirect(w,r,"http://google.com", http.StatusPermanentRedirect)
}

type CreateShortURLResponse struct {
	Short string `json:"short"`
	Long string `json:"long"`
}

func create(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	response := CreateShortURLResponse{
		Short: "unknown",
		Long: "unknown",
	}

	json.NewEncoder(w).Encode(&response)
	
}

func main() {
	fmt.Printf("Hello, world")

	router := mux.NewRouter()
	router.HandleFunc("/", home)
	router.HandleFunc("/vitals", healthcheck)
	router.HandleFunc("/create", create).Methods("POST")
	router.HandleFunc("/{code}", redirect)

	log.Printf("Listening on port 8080\n")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
