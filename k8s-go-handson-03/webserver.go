package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

type BluemixServiceConfig struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

func getBluemixServiceConfig(fileName string) BluemixServiceConfig {
	fmt.Printf("get bluemix service config from file:%s \n", fileName)
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error when get bluemix service config", err.Error())
		//os.Exit(1)
	}
	//fmt.Printf("%s\n", string(raw))
	var bluemixServiceConfig BluemixServiceConfig
	json.Unmarshal(raw, &bluemixServiceConfig)
	fmt.Printf("json:%v\n", bluemixServiceConfig)
	return bluemixServiceConfig
}

func indexHandler(redisClient *redis.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count, err := redisClient.Incr("counter").Result()
		if err != nil {
			fmt.Println("error at incr", err)
		}
		strCount := strconv.FormatInt(count, 10)
		fmt.Println("Accessing:", r.URL.Path, "Count:"+strCount)
		var version = "Version 1.0"
		var hostName = os.Getenv("HOSTNAME")
		var output = version + "\nHello Bluemix Kubernetes Cluster! \n" + "HostName:" + hostName + "\n" + "Review Count:" + strCount + "\n"
		fmt.Fprintf(w, output)
	})
}

func killHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Kill server!")
	os.Exit(1)
}

func getRedisClient() *redis.Client {
	var redisClient *redis.Client
	var bxServiceFile = os.Getenv("BLUEMIX_SERVICE_FILE")
	if bxServiceFile == "" {
		bxServiceFile = "/bluemix-service/binding"
	}
	bluemixService := getBluemixServiceConfig(bxServiceFile)
	redisHost := bluemixService.Hostname
	redisPort := bluemixService.Port
	password := bluemixService.Password

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: password,
		DB:       0,
	})
	fmt.Println("Connected to redis host->" + redisHost + ":" + redisPort)
	return redisClient
}

func main() {

	redisDB := getRedisClient()
	http.HandleFunc("/kill", killHandler)
	//http.HandleFunc("/", indexHandler)
	http.Handle("/", indexHandler(redisDB))
	port := os.Getenv("WEBAPP_PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("WebServer listening port:" + port)
	http.ListenAndServe(":"+port, nil)
}
