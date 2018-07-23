package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.Handle("/admin", requireAdminCookie(http.HandlerFunc(handleAdmin)))
	mux.HandleFunc("/admin", requireAdminCookie2(handleAdmin))
}

func requireAdminCookie(h http.Handler) http.Handler {
	return nil
}

func requireAdminCookie2(h http.HandlerFunc) http.HandlerFunc {
	return nil
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {

}
