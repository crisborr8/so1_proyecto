package main

	import (
		"github.com/julienschmidt/httprouter"
		"net/http"
	)

func (app *application) routes() http.Handler{
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/logs", app.getAllLogs)
	router.HandlerFunc(http.MethodGet, "/live", app.getLive)
	router.HandlerFunc(http.MethodGet, "/last", app.getLast10)
	router.HandlerFunc(http.MethodGet, "/count" , app.getCount)
	router.HandlerFunc(http.MethodGet, "/purge", app.purge)
	router.HandlerFunc( http.MethodGet, "/", app.getTest )
	return app.enableCORS(router)
}

