package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {

	var js []byte
	var err error
	if wrap == "" {
		js, err = json.MarshalIndent(data, "", "\t")
		if err != nil {
			return err
		}
	} else {
		wrapper := make(map[string]interface{})
		wrapper[wrap] = data

		js, err = json.Marshal(wrapper)

		if err != nil {
			return err
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err2 := w.Write(js)
	if err2 != nil {
		return err2
	}
	return nil

}

func (app *application) errorJSON(w http.ResponseWriter, err error) {
	log.Println(err)
	type jsonError struct {
		Message string `json:"message"`
	}
	theError := jsonError{
		Message: err.Error(),
	}

	err = app.writeJSON(w, http.StatusBadRequest, theError, "error")
	if err != nil {
		return
	}

}
