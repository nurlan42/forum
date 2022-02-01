package server

import (
	"log"
	"net/http"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ErrorHandler(w http.ResponseWriter, code int, msg string) {
	srvErr := struct {
		ErrCode int
		ErrMsg  string
	}{ErrCode: code, ErrMsg: msg}

	w.WriteHeader(code)

	err := tmpl.ExecuteTemplate(w, "error.html", srvErr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
