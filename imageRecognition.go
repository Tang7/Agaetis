package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// ImagePrefix for copied image
var ImagePrefix = "image-"

// ImageRecognition handler to show ImageRecognition Page
func ImageRecognition(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		renderStaticPage(w, "ImageRecognition.html")
		return
	}
	// Make a POST request when submit
	// Form a image file from file type input and return a multipart file
	f, _, err := r.FormFile("image")
	if err != nil {
		log.Print("Error when form file ", err)
	}
	defer f.Close()
	// create a temp file
	t, err := ioutil.TempFile(TempImgPath, ImagePrefix)
	if err != nil {
		log.Print("Error when create temp file ", err)
	}
	defer t.Close()
	// copy file f to temp file t
	_, err = io.Copy(t, f)
	if err != nil {
		log.Print("Error when copy file ", err)
	}
	idOffset := len(ImagePrefix) + len(TempImgPath) - 2 // len("./")
	http.Redirect(w, r, "/imageRecognition/viewImage?id="+t.Name()[idOffset:], 302)
}

// viewImage serve copied temp file and display it
func viewImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, TempImgPath+ImagePrefix+r.FormValue("id"))
}
