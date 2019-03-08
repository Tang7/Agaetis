package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type PageVariables struct {
	Data string
	Time string
}

var countryMap = map[string]string {
	"Finland": "Europe/Helsinki",
}

var homePage = template.Must(template.ParseFiles("HomePage.html"))

func check(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		curLocation, _ := time.LoadLocation(countryMap["Finland"])
		t := time.Now().In(curLocation)
		HomePageVars := PageVariables {
			Data: t.Format("Mon Jan _2 2006"),
			Time: t.Format("3:04PM"),
		}
		fmt.Println(HomePageVars)

		homePage.Execute(w, HomePageVars)
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
				homePage.Execute(w, e)
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
