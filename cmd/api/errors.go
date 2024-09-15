package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err)
	http.Error(w, "сервер не смог обработать запрос", http.StatusInternalServerError)
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "запрошенный ресурс не найден", http.StatusNotFound)
}
