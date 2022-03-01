package server

import (
	"forum/pkg/models"
	"log"
	"net/http"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// ErrorHandler handles
// func (s *AppContext) ErrorHandler(w http.ResponseWriter, code int, msg string) {
// 	srvErr := struct {
// 		ErrCode int
// 		ErrMsg  string
// 	}{ErrCode: code, ErrMsg: msg}

// 	w.WriteHeader(code)

// 	err := s.Template.ExecuteTemplate(w, "error.html", srvErr)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

func (s *AppContext) serverErr(w http.ResponseWriter) {
	err :=  models.Err{ErrCode: 500,  ErrMsg:  "Internal Server Error"}
	w.WriteHeader(http.StatusInternalServerError)
	errTemp := s.Template.ExecuteTemplate(w, "error.html", err)
	if errTemp != nil {
		http.Error(w, errTemp.Error(), http.StatusInternalServerError)
	}
}

func (s *AppContext) clientErr(w http.ResponseWriter, err models.Err){
	w.WriteHeader(err.ErrCode)
	errTemp := s.Template.ExecuteTemplate(w, "error.html", err)
	if errTemp != nil {
		http.Error(w, errTemp.Error(), http.StatusInternalServerError)
	}
}

func (s *AppContext) notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	err := models.Err{ErrCode: http.StatusNotFound, ErrMsg: "Not Found"}
	errTemp := s.Template.ExecuteTemplate(w, "error.html", err)
	if errTemp != nil {
		http.Error(w, errTemp.Error(), http.StatusInternalServerError)
	}
}

func (s *AppContext) badReq(w http.ResponseWriter) {
	err := models.Err{ErrCode: http.StatusNotFound, ErrMsg:  "Bad Request"}
	w.WriteHeader(http.StatusNotFound)
	errTemp := s.Template.ExecuteTemplate(w, "error.html",err)
	if errTemp != nil {
		http.Error(w, errTemp.Error(), http.StatusInternalServerError)
	}
}


func (s *AppContext) methodNotAllowed(w http.ResponseWriter) {
	err := models.Err{ErrCode: http.StatusMethodNotAllowed, ErrMsg: "Method Not Allowed"}
	w.WriteHeader(http.StatusMethodNotAllowed)
	errTemp := s.Template.ExecuteTemplate(w, "error.html", err)
	if errTemp != nil {
		http.Error(w, errTemp.Error(), http.StatusInternalServerError)
	}
}

 
 

