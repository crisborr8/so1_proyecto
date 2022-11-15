package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"kafka_API/Models"
	"net/http"
	"time"
)

func (app *application) insertRegister(w http.ResponseWriter, r *http.Request) {

	var newMatch Models.Match

	err := json.NewDecoder(r.Body).Decode(&newMatch)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	fmt.Println( newMatch )

//	newMatch.Team2 = strings.ReplaceAll(newMatch.Team2, " ", "_")
//	newMatch.Team1 = strings.ReplaceAll(newMatch.Team1, " ", "_")
//	newMatch.Score = strings.ReplaceAll(newMatch.Score, " ", "")



	blob, err := json.Marshal(newMatch)

	if err != nil {
		app.errorJSON(w, err)
		return
	}


/*
	jsonString := string(blob)*/

	conn , err := kafka.DialLeader(context.Background(), "tcp", app.config.kafkaConn, "topic_matches", 0)
	if err != nil {
		fmt.Println(err)
		app.errorJSON(w, err)
		return
	}
	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		fmt.Println(err)
		app.errorJSON(w, err)
		return
	}

	_, err = conn.WriteMessages(  kafka.Message{Value: blob})
	if err != nil {
		fmt.Println(err)
		app.errorJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, "ok!", "")
	if err != nil {
		fmt.Println(err)
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

