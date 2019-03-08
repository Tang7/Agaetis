package main

import (
	"net/http"
	"time"
)

// ChangeLog handler to show change log page
func ChangeLog(w http.ResponseWriter, r *http.Request) {
	curLocation, _ := time.LoadLocation(countryMap["Finland"])
	t := time.Now().In(curLocation)
	curPageVars := PageVariables{
		Data: t.Format("Mon Jan _2 2006"),
		Time: t.Format("3:04PM"),
	}
	render(w, "ChangeLog.html", curPageVars)
}
