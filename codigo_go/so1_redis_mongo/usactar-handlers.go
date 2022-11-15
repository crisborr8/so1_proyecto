package main

import (
	"fmt"
	"net/http"
	"so1_redis_mongo/Models"
)

func (app *application) getAllLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := app.models.DBMongo.GetAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if logs == nil {
		arr := []*Models.Data {}
		err = app.writeJSON(w, http.StatusOK,arr , "")
	}else {
		err = app.writeJSON(w, http.StatusOK, logs, "")
	}


	if err != nil {
		app.errorJSON(w, err)
		return
	}
	//w.Header().Set("Content-Type", "application/json")
	return
}


func (app *application) getLive(w http.ResponseWriter, r *http.Request) {
	live, err := app.redsModel.DBRedis.GetLive()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, live, "")

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	//w.Header().Set("Content-Type", "application/json")
	return

}

func (app *application)getLast10(writer http.ResponseWriter, request *http.Request) {
	last10, err := app.models.DBMongo.GetLast10()
	if err != nil {
		app.errorJSON(writer, err)
		return
	}

	err = app.writeJSON(writer, http.StatusOK, last10, "")

	if err != nil {
		app.errorJSON(writer, err)
		return
	}

	return
}

func (app *application)getCount(w http.ResponseWriter, r *http.Request) {
	count, err := app.models.DBMongo.GetCount()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, Models.Count{ Cantidad: count}, "")

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	return
}

func (app *application)purge(w http.ResponseWriter, r *http.Request) {
	msg, err := app.models.DBMongo.Purge()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, msg, "")

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	return

}



func (app *application) getTest(writer http.ResponseWriter, request *http.Request) {

	_, err := writer.Write([]byte("hola!!"))
	if err != nil {
		fmt.Println("prueba fallida no se ha podido devolver el valor")
		return
	}

	return
}

