package main

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"so1_redis_mongo/Models"
	"time"
)

type config struct {
port int
env  string
dbc  string
dbcR string
}

type application struct {
	config config
	logger *log.Logger
	models Models.Models
	redsModel Models.ModelRedis
}



func main ( ) {

	dbc := "mongodb://"
	if  os.Getenv("MONGO") == "" {
		dbc +=  "192.168.1.18:27017"
	}else {
		dbc += os.Getenv("MONGO") + ":27017"
	}

	dbcR := ""
	if os.Getenv("REDIS") == "" {
		dbcR = "192.168.1.18:6379"
	}else {
		dbcR = os.Getenv("REDIS") + ":6379"
	}

	cfg := config{port: 8080, env: "development", dbc: dbc, dbcR: dbcR}

	logger := log.New(os.Stdout, "Log", log.Ldate|log.Ltime)

	dbMongo, err := openDB(cfg)

	if err != nil {
		log.Fatal(err)
	}



	app := &application{
		config: cfg,
		logger: logger,
		models: Models.NewModels(dbMongo),
		redsModel: Models.NewModelRedis( openDBRedis(cfg) ),
	}


	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	logger.Println("Iniciando el servidor en el puerto:", cfg.port)

	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}



}

func openDB(cfg config) (*mongo.Client, error) {
	//client, err := mongo.NewClient(options.Client().ApplyURI(cfg.dbc))

	mongoDBConnectionString := "#cadena de conexión a mongo"


	clientOptions := options.Client().ApplyURI(mongoDBConnectionString).SetDirect(true)

	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		fmt.Println("Error base conexión no realizada con éxito")
		return nil, err
	}

	fmt.Println("Se ha realizado la conexión a la base de datos correctamente")
	return client, nil
}

func openDBRedis(cfg config)(*redis.Pool) {
	pool := &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 60 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", cfg.dbcR) },
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return pool
}

