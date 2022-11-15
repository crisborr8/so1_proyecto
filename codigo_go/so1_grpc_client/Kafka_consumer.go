package main

import (
	services_matchpb "Kafka_Redis_Client/GoOutput/SendMatches"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Match struct {
	Team1 string `json:"team1" bson:"team1"`
	Team2 string `json:"team2" bson:"team2"`
	Score string `json:"score" bson:"score"`
	Phase int    `json:"phase" bson:"phase"`
}

/*
{
"team1": "Guatemala",
"team2": "Argentina",
"score":"2-5",
"phase":<1,2,3,4>
}
Esto
*/

func StartKafka(host string, port string, hostGrpc string, mongoConect string) {
	conf := kafka.ReaderConfig{
		Brokers:        []string{host + ":" + port},
		Topic:          "topic_matches",
		MaxBytes:       10,
		SessionTimeout: 10 * time.Second,
	}

	conn, err := kafka.Dial("tcp", host + ":" + port)
	defer conn.Close()


	controller, _ := conn.Controller()
	conncontroller, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Fatal(err)
	}
	err = conncontroller.CreateTopics(kafka.TopicConfig{
		Topic:              "topic_matches",
	})
	if err != nil {
		log.Fatal(err)
	}


	cc, err := grpc.Dial(hostGrpc + ":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	matchConnect := services_matchpb.NewSendMatchesClient(cc)
	mongoClient, err := openDB(mongoConect)

	if err != nil {
		log.Fatal(err)
	}

	reader := kafka.NewReader(conf)

	var MatchRead Match

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			//fmt.Println("Some error occured", err)
			continue
		} else {
			wg := new(sync.WaitGroup)
			wg.Add(2)
			err = json.Unmarshal(m.Value, &MatchRead)

			if err != nil {
				fmt.Println("Error al parsear", err)
			}

			go func() {
				// send to redis

				req := &services_matchpb.Match{
					Team1: MatchRead.Team1,
					Team2: MatchRead.Team2,
					Score: strings.TrimSpace(MatchRead.Score),
					Phase: int32(MatchRead.Phase),
				}

				res, err := matchConnect.SendMatch(context.Background(), req)
				if err != nil || res.GetError() {
					fmt.Println("Error al enviar", err)
				}
				wg.Done()
			}()

			go func() {
				// send to mongo
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				_, err := mongoClient.Database("so1_proyecto").
					Collection("logs").InsertOne(ctx,MatchRead )
				if err != nil {
					fmt.Println("Error al enviar", err)
				}
				wg.Done()
			}()

			wg.Wait()

		}
		time.Sleep( 100 * time.Millisecond)

	}
}

func openDB(mongoConect string) (*mongo.Client, error) {
	//credential := options.Credential{
	//	Username: "dbuser",
	//	Password: "sopes1",
	//}
	//client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + mongoConect).SetAuth(
	//	credential ).SetDirect(true).SetAppName(""))

	mongoDBConnectionString := "#cadena conexion con mongo"


	clientOptions := options.Client().ApplyURI(mongoDBConnectionString).SetDirect(true)

	client, err := mongo.NewClient(clientOptions)


	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Se ha realizado la conexión a la base de datos correctamente")
	return client, nil
}
