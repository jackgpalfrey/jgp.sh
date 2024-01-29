package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var words = Words{
	wordList: make([]string, 0),
}
var db *sql.DB

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
		Name:    "JGP.SH",
		Version: "indev",
		Status:  "OK",
	}

	err := json.NewEncoder(w).Encode(&report)
	if err != nil {
		log.Fatal(err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	short := params["short"]

	rows, err := db.Query(`SELECT long FROM urls WHERE short=$1`, short)
	defer rows.Close()
	if err != nil {
		msg := err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "`+msg+`"}`)
		return
	}

	rows.Next()
	var long string
	rows.Scan(&long)
	log.Printf("%s", long)
	if long == "" {
		msg := "URL '" + short + "' not found"
		fmt.Fprintf(w, `{"error": "`+msg+`"}`)
		return
	}

	http.Redirect(w, r, long, http.StatusPermanentRedirect)
}

func enumerate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query(`SELECT short, long FROM urls`)
	defer rows.Close()
	if err != nil {
		msg := err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "`+msg+`"}`)
		return
	}

	enumeration := map[string]string{}

	var short, long string
	for rows.Next() {
		err := rows.Scan(&short, &long)
		if err != nil {
			log.Fatal(err)
		}
		enumeration[short] = long
	}

	json.NewEncoder(w).Encode(&enumeration)
}

type CreateShortURLResponse struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

type CreateShortURLRequest struct {
	Long string `json:"long"`
}

func create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reqBody := CreateShortURLRequest{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "`+err.Error()+`"}`)
		return
	}

	long := reqBody.Long

	short, err := words.getRandom()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Server Error"}`)
		return
	}

	_, err = db.Exec("INSERT INTO urls (short,long) VALUES ($1,$2)", short, long)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "`+err.Error()+`"}`)
		return
	}

	response := CreateShortURLResponse{
		Short: short,
		Long:  long,
	}

	json.NewEncoder(w).Encode(&response)
}

func main() {
	log.Printf("Starting Server...")

	log.Printf("Loading Words...")
	loadErr := words.loadFrom("./words")
	if loadErr != nil {
		log.Fatal(loadErr)
	}
	log.Println("Loaded")

	log.Printf("Connecting to DB...")
	db = connectToDatabase()
	log.Println("Connected")

	router := mux.NewRouter()
	router.HandleFunc("/", home)
	router.HandleFunc("/vitals", healthcheck)
	router.HandleFunc("/create", create).Methods("POST")
	router.HandleFunc("/enumerate", enumerate)
	router.HandleFunc("/{short}", redirect)

	log.Printf("Listening on port 8080\n")
	serveErr := http.ListenAndServe(":8080", router)
	if serveErr != nil {
		log.Fatal(serveErr)
	}
}
