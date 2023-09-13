package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/sclevine/agouti"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/cookie", handleCookie)
	server := httptest.NewServer(mux)
	defer server.Close()

	driver := agouti.ChromeDriver()
	//driver := agouti.PhantomJS()
	if err := driver.Start(); err != nil {
		log.Fatal(err)
	}
	defer driver.Stop() //nolint:errcheck

	page, err := driver.NewPage()
	if err != nil {
		log.Fatal(err)
	}
	if err := page.Navigate(server.URL + "/cookie"); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 3)

	cookies, err := page.GetCookies()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cookies = %+v\n", cookies)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("test")
	if err != nil {
		http.Error(w, fmt.Sprintf("No cookie\n\n%v", err), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "cookie = %+v", cookie)
}

func handleCookie(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "test",
		Value:    "test",
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}
