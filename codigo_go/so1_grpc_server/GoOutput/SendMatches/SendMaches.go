package services_matchpb

import (
	context "context"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
	"strings"
)

type Server struct {
	RedisPool *redis.Pool
}

func (s Server) SendMatch(ctx context.Context, match *Match) (*Confirm, error) {
	conn := s.RedisPool.Get()
	defer conn.Close()

	match.Team2 = strings.ReplaceAll(match.GetTeam2(), " ", "%")

	match.Team1 = strings.ReplaceAll(match.GetTeam1(), " ", "%")

	_, err := conn.Do("HSET", strconv.Itoa(int(match.GetPhase())),
		strings.ToLower(match.GetTeam1())+"_"+
			strings.ToLower(match.GetTeam2()), "true")

	if err != nil {
		log.Fatalln(err)
	}

	countB, _ := conn.Do("HGET", strconv.Itoa(int(match.GetPhase()))+"_"+
		strings.ToLower(match.GetTeam1())+"_"+
		strings.ToLower(match.GetTeam2()),
		match.GetScore())
	count, _ := redis.Int(countB, err)

	_, err = conn.Do("HSET", strconv.Itoa(int(match.GetPhase()))+"_"+
		strings.ToLower(match.GetTeam1())+"_"+
		strings.ToLower(match.GetTeam2()),
		match.GetScore(), count+1)

	return &Confirm{
		Error: err != nil,
	}, nil
}

func (s Server) mustEmbedUnimplementedSendMatchesServer() {
	panic("implement me")
}
