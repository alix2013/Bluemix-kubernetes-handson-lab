package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	redis "github.com/go-redis/redis"
)

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
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "",
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
