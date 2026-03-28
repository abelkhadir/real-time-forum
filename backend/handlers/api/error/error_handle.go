package error_handler

import (
	"fmt"
	"html/template"
	"net/http"
)

type Error struct {
	ErrorMessage string
	ErrorStatus  int
}

func ErrorPage(w http.ResponseWriter, errMessage string, errStatus int) {
	template, err := template.ParseGlob("./frontend/error.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internel server Error", http.StatusInternalServerError)
		return
	}

	Data := Error{
		ErrorMessage: errMessage,
		ErrorStatus:  errStatus,
	}

	w.WriteHeader(errStatus)
	if err := template.Execute(w, Data); err != nil {
		fmt.Println(err)
		http.Error(w, "internel server Error", http.StatusInternalServerError)
		return
	}
}
