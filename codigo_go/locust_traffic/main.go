package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/myzhan/boomer"
)

// This is a tool like Apache Benchmark a.k.a "ab".
// It doesn't implement all the features supported by ab.

var client *http.Client
var postBody []byte


var count uint64
var verbose bool

var method string
var url string
var timeout int
var postFile string
var contentType string
var concurrency int

var disableCompression bool
var disableKeepalive bool
var iterations uint64

var countrys = []string{
	"Catar",
	"Ecuador",
	"Senegal",
	"Países Bajos",
	"Inglaterra",
	"Irán",
	"Estados Unidos",
	"Gales",
	"Argentina",
	"Arabia Saudí",
	"México",
	"Polonia",
	"Francia",
	"Australia",
	"Dinamarca",
	"Túnez",
	"España",
	"Costa Rica",
	"Alemania",
	"Japón",
	"Bélgica",
	"Canadá",
	"Marruecos",
	"Croacia",
	"Brasil",
	"Serbia",
	"Suiza",
	"Camerún",
	"Portugal",
	"Ghana",
	"Uruguay",
	"Corea del Sur",
}

func randomBody() []byte {
	//var value  = ""



	rand.Seed(time.Now().UnixNano())
	team1 := 0 + rand.Intn(len(countrys))
	team2 := 0 + rand.Intn(len(countrys))

	score1 := rand.Intn(10)
	score2 := rand.Intn(10)
	phase := 1 + rand.Intn(3)

	json := fmt.Sprintf(` {
    "team1": "%s",
    "team2": "%s",
    "score": "%d-%d",
    "phase": %d
} `, countrys[team1], countrys[team2],
		score1, score2, phase)
	fmt.Println(json)

	return []byte(json)
}

func workerR() {
	//fmt.Println(url)
	request, err := http.NewRequest(method, url, bytes.NewBuffer(randomBody()))
	if err != nil {

		log.Fatalf("%v\n", err)
	}

	request.Header.Set("Content-Type", contentType)

	startTime := time.Now()
	response, err := client.Do(request)
	elapsed := time.Since(startTime)

	if err != nil {
		if verbose {
			log.Printf("%v\n", err)
		}
		boomer.RecordFailure("http", "error", 0.0, err.Error())
	} else {
		boomer.RecordSuccess("http", strconv.Itoa(response.StatusCode),
			elapsed.Nanoseconds()/int64(time.Millisecond), response.ContentLength)

		if verbose {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Printf("%v\n", err)
			} else {
				//log.Printf("Status Code: %d\n", response.StatusCode)
				log.Println(string(body))
			}

		} else {
			io.Copy(io.Discard, response.Body)
		}

		response.Body.Close()
	}
	iterations++

	if iterations > count {
		os.Exit(0)
	}

}

func worker() {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(postBody))
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	request.Header.Set("Content-Type", contentType)

	startTime := time.Now()
	response, err := client.Do(request)
	elapsed := time.Since(startTime)

	if err != nil {
		if verbose {
			log.Printf("%v\n", err)
		}
		boomer.RecordFailure("http", "error", 0.0, err.Error())
	} else {
		boomer.RecordSuccess("http", strconv.Itoa(response.StatusCode),
			elapsed.Nanoseconds()/int64(time.Millisecond), response.ContentLength)

		if verbose {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Printf("%v\n", err)
			} else {
				log.Printf("Status Code: %d\n", response.StatusCode)
				log.Println(string(body))
			}

		} else {
			io.Copy(ioutil.Discard, response.Body)
		}

		response.Body.Close()
	}
	iterations++

	if iterations > count {
		os.Exit(0)
	}
}

func main() {
	iterations = 0;

	argsWithoutProg := os.Args[1:]
	params := make(map[string]string)

	for i := 0; i < len(argsWithoutProg); i += 2 {
		params[argsWithoutProg[i]] = argsWithoutProg[i+1]
	}

	//fmt.Println(argsWithoutProg)


	/*switch (argsWithoutProg [ ]) {
	}
	*/

	url = params["--url"]
	if url == "" {
		url = "http://35.193.180.239.nip.io/input/v1/insert"
	}

	if params ["-n"] == ""{
		count = math.MaxUint64 - 1
	}else {
		count, _ = strconv.ParseUint(params["-n"],10,64)
	}


	if params["--concurrency"] == ""{
		concurrency = 1
	}else {
		concurrency, _ = strconv.Atoi(params["--concurrency"])
	}

	timeoutS := params["--timeout"]
	if timeoutS == "" {
		timeoutS = "100"
	}



	timeoutInt, _ := strconv.Atoi(timeoutS)

	if params["-f"] == "" {
		log.Println("[INFO] Archivo .json no específicado, se va a utilizar tráfico aleatorio")
		//log.Fatalln("Error grave: No se ha especificado el archivo de entrada ")
		//return
		os.Args = []string{"--run-task", "wokerR"}
	} else {
		os.Args = []string{"--run-task", "woker"}
	}

	flag.StringVar(&method, "method", "POST", "HTTP method, one of GET, POST")
	flag.StringVar(&url, "url", url, "URL")
	flag.IntVar(&timeout, "timeout", timeoutInt, "Milisegundos para esperar la solicitud")

	flag.StringVar(&postFile, "post-file", params["-f"], "File containing data to POST. Remember also to set --content-type")
	flag.StringVar(&contentType, "content-type", "application/json", "Content-type header")

	flag.BoolVar(&disableCompression, "disable-compression", false, "Disable compression")
	flag.BoolVar(&disableKeepalive, "disable-keepalive", false, "Disable keepalive")

	flag.BoolVar(&verbose, "verbose", true, "Print debug log")

	//fmt.Print(f)
	flag.Parse()

	log.Printf(`HTTP Locust: benchmark is running with these args:
method: %s
url: %s
timeout: %d
post-file: %s
content-type: %s
disable-compression: %t
disable-keepalive: %t
verbose: %t`, method, url, timeout, postFile,
		contentType, disableCompression, disableKeepalive, verbose)

	if url == "" {
		log.Fatalln("--url can't be empty string, please specify a URL that you want to test.")
	}

	if method != "GET" && method != "POST" {
		log.Fatalln("HTTP method must be one of GET, POST.")
	}

	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 2000
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConnsPerHost: 2000,
		DisableCompression:  disableCompression,
		DisableKeepAlives:   disableKeepalive,
	}
	client = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeout) * time.Millisecond,
	}

	task := &boomer.Task{
		Name:   "worker",
		Weight: concurrency,
		Fn:     worker,
	}
	taskR := &boomer.Task{
		Name:   "workerR",
		Weight: concurrency,
		Fn:     workerR,
	}

	if method == "POST" {
		if postFile != "" {
			//log.Fatalln("--post-file can't be empty string when method is POST")
			tmp, err := os.ReadFile(postFile)
			if err != nil {
				log.Fatalf("%v\n", err)
			}
			postBody = tmp
			boomer.Run(task)
		} else {
			// se ejecuta de forma aleatoria
			boomer.Run(taskR)
		}
	} else {
		boomer.Run(task)
	}
}
