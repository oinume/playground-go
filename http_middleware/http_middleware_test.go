package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAdmin(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/admin", requireAdminCookie(http.HandlerFunc(admin)))
	s := httptest.NewServer(mux)
	defer s.Close()

	resp, err := http.Get(s.URL + "/admin")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := resp.StatusCode, http.StatusForbidden; got != want {
		t.Errorf("status code mismatch: got=%v, want=%v", got, want)
	}
}

func admin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
