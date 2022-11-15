package main

import (
	"fmt"
	"os"
)

func main()  {
	fmt.Println("test")
	//fmt.Println("Okay...")
	var portKafka, hostKafka, hostGrpc, hostMongo string
	fmt.Println("=== KAFKA  ====")
	if os.Getenv("KAFKA") == "" {
		hostKafka = "localhost"
	}else {
		fmt.Println(os.Getenv("KAFKA"))
		hostKafka = os.Getenv("KAFKA")
	}

	portKafka = "9092"
	fmt.Println("=== GRPC  ====")

	if os.Getenv("GRPC") == "" {
		hostGrpc = "localhost"
	}else {
		fmt.Println(os.Getenv("GRPC"))
		hostGrpc = os.Getenv("GRPC")
	}

	fmt.Println("=== MONGASO  ====")

	if os.Getenv("MONGO") == "" {
		hostMongo = "192.168.1.18:27017"
	}else {
		fmt.Println(os.Getenv("MONGO"))
		//10.32.8.133
		hostMongo = os.Getenv("MONGO") + ":27017"
	}

	StartKafka(hostKafka, portKafka, hostGrpc, hostMongo)

	//fmt.Println("Kafka has been started...")

	//time.Sleep(2 * time.Minute)

}
