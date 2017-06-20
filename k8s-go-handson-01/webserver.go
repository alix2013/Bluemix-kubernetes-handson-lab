package main

import (
	"fmt"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Accessing:", r.URL.Path)
	var version = "Version 2.0"
	var hostName = os.Getenv("HOSTNAME")
	var output = version + "\nHello Bluemix Kubernetes Cluster! \n" + "HostName:" + hostName + "\n"
	fmt.Fprintf(w, output)
}
func killHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Kill server!")
	os.Exit(1)
}

func main() {
	http.HandleFunc("/kill", killHandler)
	http.HandleFunc("/", indexHandler)
	port := os.Getenv("WEBAPP_PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("WebServer listening port:" + port)
	http.ListenAndServe(":"+port, nil)
}

