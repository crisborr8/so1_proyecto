/* package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "192.168.1.16:9092", "topic_test2", 0)



	if err != nil {
		fmt.Println("no connection")
		return
	}
	err = conn.SetWriteDeadline(time.Now().Add(time.Second * 3))

	if err != nil {
		return
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	done := make(chan bool, 1)

	//kafkaReader := kafka.NewReader( kafka.ReaderConfig{
	//	Brokers: []string{"192.168.1.16:9092"},
	//	Topic: "topic_test2",
	//	Partition: 0,
	//	MinBytes:  10e3, // 10KB
	//	MaxBytes:  10e6, // 10
	//	MaxWait: 10 * time.Second,
	//} )
	//



	go func() {

		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		_ = conn.Close()
		done <- true
	}()


	batch := conn.ReadBatch(1e3, 1e6) // 1 - 1000
	bytes := make([]byte, 1e3)
	i := 1;
	for {
		_, err := batch.Read(bytes)
		if err != nil {
			time.Sleep(1000 * time.Millisecond)
			break;

		}else {
			fmt.Println( i, " " , string(bytes))
			i++;

			bytes[0] = 0;
			for j := 1; j < len(bytes); j *= 2 {
				copy(bytes[j:],bytes[:j])
			}
		}

	}
	_ = conn.Close()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

}


*/

package main

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

type config struct {
	port int
	env  string
	kafkaConn  string

}

type application struct {
	config config
	logger *log.Logger
}


func main()  {

	/*fmt.Println("Okay...")
	//appKafka.StartKafka()
	//
	//fmt.Println("Kafka has been started...")
	//
	//time.Sleep(30 * time.Second) */



	var kafkaConn string
	if os.Getenv("KAFKA") == "" {
		kafkaConn = "localhost:9092"
	}else {
		kafkaConn = os.Getenv("KAFKA") + ":9092"
	}

	fmt.Println(kafkaConn)

	conn, err := kafka.Dial("tcp", kafkaConn)
	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()


	controller, err  := conn.Controller()
	if err != nil {
		log.Fatal(err)
	}

	conncontroller, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Fatal(err)
	}
	err = conncontroller.CreateTopics(kafka.TopicConfig{
		Topic:              "topic_matches",
		NumPartitions: 1,
		ReplicationFactor: 1,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("creado correctamente")




	cfg := config{port: 80, env: "development", kafkaConn: kafkaConn}

	app := &application{
		config: cfg,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  2 * time.Minute,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 60 * time.Second,

	}

	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}