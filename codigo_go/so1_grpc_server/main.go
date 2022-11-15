package main

import (
	services_matchpb "Redis_Server/GoOutput/SendMatches"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

func main() {

	dbcR := ""
	if os.Getenv("REDIS") == "" {
		dbcR = "192.168.1.18"
	}else {
		dbcR = os.Getenv("REDIS")
	}
	dbcR += ":6379"



	pool := &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 60 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", dbcR) },
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	listener, err := net.Listen("tcp","0.0.0.0:50051")
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()

	services_matchpb.RegisterSendMatchesServer( s, &services_matchpb.Server{ RedisPool: pool  } )
	if err = s.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
