package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", showForm)
	mux.HandleFunc("/upload", upload)

	log.Println("serving on :5678")
	if err := http.ListenAndServe(":5678", mux); err != nil {
		log.Fatalf("failed to ListenAndServe: %v", err)
	}
}

func showForm(w http.ResponseWriter, r *http.Request) {
	html := `
<html>
<head><title>Multiple file upload example in Go</title></head>
<body>
<h1>Select file to upload</h1>
<form action="/upload" method="post" enctype="multipart/form-data">
Select file to upload:
<p><input type="file" name="file1" id="file1"></p>
<p><input type="file" name="file2" id="file2"></p>
<input type="submit" value="Upload" name="submit">
</form>
</body>
</html>
`
	fmt.Fprintln(w, strings.TrimSpace(html))
}

func upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(5 * 1024 * 1024 * 1024 /* 5MB */); err != nil {
		internalServerError(w, fmt.Errorf("failed to parse form"))
		return
	}
	for _, name := range []string{"file1", "file2"} {
		file, handler, err := r.FormFile(name)
		if err != nil {
			internalServerError(w, err)
			return
		}
		f, err := os.Create("/tmp/" + handler.Filename)
		if err != nil {
			internalServerError(w, err)
			return
		}
		// Copy file content
		if _, err := io.Copy(f, file); err != nil {
			internalServerError(w, err)
			return
		}
		// Show a name of uploaded file
		fmt.Fprintln(w, "Uploaded files")
		fmt.Fprintln(w, handler.Filename)
		f.Close()
	}
}

func internalServerError(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Sprintf("Internal Server Error\n\n%v", err), http.StatusInternalServerError)
}
