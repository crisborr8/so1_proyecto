package Models

import (
	"github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	DBMongo DBMongoModel
}

type ModelRedis struct {
	DBRedis DBRedisModel
}

func NewModelRedis(db *redis.Pool) ModelRedis {
	return ModelRedis{
		DBRedis: DBRedisModel{
			DB: db,
		},
	}
}

func NewModels(db *mongo.Client) Models {
	return Models{
		DBMongo: DBMongoModel{
			DB: db,
		},
	}
}

type Data struct {
	Team1 string `json:"team1" bson:"team1" `
	Team2 string `json:"team2" bson:"team2" `
	Score string ` json:"score" bson:"score"`
	Phase int `json:"phase" bson:"phase"`
}


type Match struct {
	Team1 string `json:"team1" bson:"team1" `
	Team2 string `json:"team2" bson:"team2" `
	Score string ` json:"score" bson:"score"`
	Phase string `json:"phase" bson:"phase"`
}

type Partido struct {
	Team1       string            `json:"team_1" bson:"team_1"`
	Team2        string            `json:"team_2" bson:"team_2"`
	Predicciones map[string]string `json:"predicciones" bson:"predicciones"`
}

type Live struct {
	Octavos []*Partido `json:"octavos" bson:"octavos"`
	Cuartos []*Partido `json:"cuartos" bson:"cuartos"`
	Semis   []*Partido `json:"semis" bson:"semis"`
	Final   []*Partido `json:"final" bson:"final"`
}

type Count struct {
	Cantidad int64 `bson:"cantidad" json:"cantidad"`
}

type Message struct {
	Message  string `json:"message" bson:"message"`
}


