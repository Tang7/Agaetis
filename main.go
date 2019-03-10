package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

// PageVariables variables to display in page
type PageVariables struct {
	Data string
	Time string
}

// ImgPath image folder path
var ImgPath = "./img"

// TempImgPath temp image folder store the image during image upload
var TempImgPath = "./tmp/"

func main() {
	// disable low level warning
	os.Setenv("TF_CPP_MIN_LOG_LEVEL", "2")
	// serve everything in below folders as a file
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	http.HandleFunc("/", Home)
	http.HandleFunc("/imageRecognition", UploadImage)
	http.HandleFunc("/changeLog", ChangeLog)
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
