package main

import "net/http"

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", app.createUserHandler)
	mux.HandleFunc("GET /tokens/get/{user_id}", app.createTokensHandler)
	mux.HandleFunc("POST /tokens/refresh", app.renewAccessTokenHandler)

	return mux
}
