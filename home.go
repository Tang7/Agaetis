package main

import (
	"net/http"
	"time"
)

var countryMap = map[string]string {
	"Finland": "Europe/Helsinki",
}

func Home(w http.ResponseWriter, r *http.Request) {
	curLocation, _ := time.LoadLocation(countryMap["Finland"])
	t := time.Now().In(curLocation)
	homePageVars := PageVariables {
		Data: t.Format("Mon Jan _2 2006"),
		Time: t.Format("3:04PM"),
	}
	render(w, "HomePage.html", homePageVars)
}