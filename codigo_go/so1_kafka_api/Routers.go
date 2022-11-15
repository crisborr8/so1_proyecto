package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler{
	router := httprouter.New()

	router.HandlerFunc( http.MethodPost, "/v1/insert", app.insertRegister )
	router.HandlerFunc( http.MethodGet, "/", app.getTest )

	return app.enableCORS(router)


}
