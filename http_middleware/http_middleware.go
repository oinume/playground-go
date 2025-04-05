package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/admin", requireAdminCookie(http.HandlerFunc(handleAdmin)))
	mux.HandleFunc("/admin2", requireAdminCookie2(handleAdmin))
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func requireAdminCookie(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("admin")
		if err != nil {
			http.Error(w, "No admin cookie", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func requireAdminCookie2(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("admin")
		if err != nil {
			http.Error(w, "No admin cookie", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "This is admin page")
}
