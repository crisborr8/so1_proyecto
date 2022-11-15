package Models

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	"sync"
	"time"
)

type DBMongoModel struct {
	DB *mongo.Client
}

type DBRedisModel struct {
	DB *redis.Pool
}

func (dbm *DBMongoModel) GetAll() ([]*Data, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cursor, err := dbm.DB.Database("so1_proyecto").Collection("logs").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var AllLogs []*Data
	err = cursor.All(ctx, &AllLogs)

	if err != nil {
		return nil, err
	}

	return AllLogs, nil
}

func (dbm *DBMongoModel) GetLast10() ([]*Data, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	opt := options.Find().SetSort(bson.D{{"_id", -1}}).SetLimit(10)

	cursor, err := dbm.DB.Database("so1_proyecto").Collection("logs").
		Find(ctx, bson.D{}, opt)
	if err != nil {
		return nil, err
	}
	var AllLogs []*Data
	err = cursor.All(ctx, &AllLogs)

	if err != nil {
		return nil, err
	}
	return AllLogs, nil
}

func (dbm *DBMongoModel) GetCount() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	count, err := dbm.DB.Database("so1_proyecto").Collection("logs").CountDocuments(ctx, bson.M{})
	if err != nil {
		return -1, err
	}

	return count, err
}


func (dbm *DBMongoModel) Purge() (*Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := dbm.DB.Database("so1_proyecto").Collection("logs").Drop(ctx)
	if err != nil {
		return nil, err
	}
	return &Message{ Message: "Contenido eliminado correctamente"}, err

}

func (dbmR *DBRedisModel) GetLive() (*Live, error) {

	conn := dbmR.DB.Get()
	conn1 := dbmR.DB.Get()
	conn2 := dbmR.DB.Get()
	conn3 := dbmR.DB.Get()
	defer conn.Close()
	defer conn1.Close()
	defer conn2.Close()
	defer conn3.Close()

	live_16sB, err := conn.Do("HKEYS", "1")
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}

	live_8sB, err := conn.Do("HKEYS", "2")
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	live_4sB, err := conn.Do("HKEYS", "3")
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	live_finalB, err := conn.Do("HKEYS", "4")
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	liveResonse := Live{
		Octavos: make([]*Partido, 0, 64),
		Cuartos: make([]*Partido, 0, 64),
		Semis:   make([]*Partido, 0, 64),
		Final:   make([]*Partido, 0, 64),
	}

	live_16s, err := redis.Strings(live_16sB, err)
	live_8s, err := redis.Strings(live_8sB, err)
	live_4s, err := redis.Strings(live_4sB, err)
	live_final, err := redis.Strings(live_finalB, err)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	// round of 16

	var readerMatchB interface{}

	wg := new(sync.WaitGroup)
	wg.Add(4)

	go func() {
		for _, t := range live_16s {
			readerMatchB, err = conn.Do("HGETALL", "1_"+t)

			if err != nil {
				log.Fatalln(err)
			}

			matches, err := redis.StringMap(readerMatchB, err)
			if err != nil {
				log.Fatalln(err)
			}

			teams := strings.Split(t, "_")

			match := &Partido{
				Team1:        strings.ReplaceAll(teams[0],"%"," "),
				Team2:        strings.ReplaceAll(teams[1], "%", " "),
				Predicciones: matches,
			}

			liveResonse.Octavos = append(liveResonse.Octavos, match)
		}

		wg.Done()
	}()

	go func() {
		for _, t := range live_8s {
			readerMatchB, err = conn1.Do("HGETALL", "2_"+t)

			if err != nil {
				log.Fatalln(err)
			}

			matches, err := redis.StringMap(readerMatchB, err)
			if err != nil {
				log.Fatalln(err)
			}

			teams := strings.Split(t, "_")

			match := &Partido{
				Team1:        strings.ReplaceAll(teams[0],"%"," "),
				Team2:        strings.ReplaceAll(teams[1], "%", " "),
				Predicciones: matches,
			}

			liveResonse.Cuartos = append(liveResonse.Cuartos, match)
		}

		wg.Done()
	}()

	go func() {
		for _, t := range live_4s {
			readerMatchB, err = conn2.Do("HGETALL", "3_"+t)

			if err != nil {
				log.Fatalln(err)
			}

			matches, err := redis.StringMap(readerMatchB, err)
			if err != nil {

				log.Fatalln(err)
			}

			teams := strings.Split(t, "_")

			match := &Partido{
				Team1:        strings.ReplaceAll(teams[0],"%"," "),
				Team2:        strings.ReplaceAll(teams[1], "%", " "),
				Predicciones: matches,
			}

			liveResonse.Semis = append(liveResonse.Semis, match)
		}

		wg.Done()
	}()

	go func() {
		for _, t := range live_final {
			readerMatchB, err = conn3.Do("HGETALL", "4_"+t)

			if err != nil {
				log.Fatalln(err)
			}

			matches, err := redis.StringMap(readerMatchB, err)
			if err != nil {
				log.Fatalln(err)
			}

			teams := strings.Split(t, "_")

			match := &Partido{
				Team1:        strings.ReplaceAll(teams[0],"%"," "),
				Team2:        strings.ReplaceAll(teams[1], "%", " "),
				Predicciones: matches,
			}

			liveResonse.Final = append(liveResonse.Final, match)
		}
		wg.Done()
	}()

	wg.Wait()
	return &liveResonse, nil
}
