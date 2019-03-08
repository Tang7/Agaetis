package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type PageVariables struct {
	Data string
	Time string
}

var ImgPath = "./img"
var TempImgPath = "./tmp/"

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/imageRecognition", ImageRecognition)
	http.HandleFunc("/imageRecognition/viewImage", viewImage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func render(w http.ResponseWriter, htmlPage string, pageVars PageVariables) {
	htmlPage = fmt.Sprintf("html/%s", htmlPage)
	curPage := template.Must(template.ParseFiles(htmlPage))
	if err := curPage.Execute(w, pageVars); err != nil {
		log.Print("Page execute error: ", err)
	}
}

func renderStaticPage(w http.ResponseWriter, htmlPage string) {
	htmlPage = fmt.Sprintf("html/%s", htmlPage)
	curPage := template.Must(template.ParseFiles(htmlPage))
	if err := curPage.Execute(w, nil); err != nil {
		log.Print("Page execute error: ", err)
	}
}
