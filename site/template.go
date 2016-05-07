package site

import (
	"html/template"
	"log"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	t := template.Must(template.ParseFiles(root + "/template/" + name + ".tpl"))
	err := t.Execute(w, data)
	if err != nil {
		log.Print(err)
	}
}

func renderError(w http.ResponseWriter, status int) {
	data := struct {
		Err int
	}{
		status,
	}

	renderTemplate(w, "error", data)
}
