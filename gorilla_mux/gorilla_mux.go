package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/v1/users/{id:[0-9]+}", usersHandler).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:5000", router))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Your id is %v", vars["id"])
}
