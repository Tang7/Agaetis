package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var htmlPage = template.Must(template.ParseFiles("upload.html"))

func check(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		htmlPage.Execute(w, nil)
		return
	}
	// Make a POST request when submit
	// Form a image file from file type input and return a multipart file
	f, _, err := r.FormFile("image")
	check(err)
	defer f.Close()
	// create a temp file
	t, err := ioutil.TempFile(".", "image-")
	check(err)
	defer t.Close()
	// copy file f to temp file t
	_, err2 := io.Copy(t, f)
	check(err2)
	http.Redirect(w, r, "/view?id="+t.Name()[6:], 302)
}

func viewImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, "image-"+r.FormValue("id"))
}

func errorHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				w.WriteHeader(500)
				htmlPage.Execute(w, e)
			}
		}()
		fn(w, r)
	}
}

func main() {
	http.HandleFunc("/", errorHandler(uploadImage))
	http.HandleFunc("/view", errorHandler(viewImage))
	http.ListenAndServe(":8080", nil)
}
