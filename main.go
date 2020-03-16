package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/soajs/soajs.golang"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Message string `json:"message"`
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	resp := Response{}
	resp.Message = fmt.Sprintf("heartbeat")
	respJson, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func hello(w http.ResponseWriter, r *http.Request) {
	soajs := r.Context().Value(soajsgo.SoajsKey).(soajsgo.ContextData)
	respJsonSOA, err := json.Marshal(soajs)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJsonSOA)
}

func interconnect(w http.ResponseWriter, r *http.Request) {
	soajs := r.Context().Value(soajsgo.SoajsKey).(soajsgo.ContextData)
	respJsonSOA, err := json.Marshal(soajs)
	if err != nil {
		panic(err)
	}

	log.Println("micro1")
	uri := soajs.Awareness.Path("micro2");
	uri = uri + "hello?key=" + soajs.Tenant.Key.EKey
	log.Println(uri)

	req, err := http.NewRequest("GET", "http://"+uri, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.Header.Set("Cache-Control", "no-cache")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
	fmt.Printf("%s\n", body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJsonSOA)
}

//main function
func main() {
	router := mux.NewRouter()

	jsonFile, err := os.Open("soa.json")
	if err != nil {
		log.Println(err)
	}
	log.Println("Successfully Opened soa.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result soajsgo.Config
	json.Unmarshal([]byte(byteValue), &result)

	soajs, err := soajsgo.NewFromConfig(context.Background(), result)
	if err != nil {
		log.Fatal(err)
	}

	router.Use(soajs.Middleware)

	router.HandleFunc("/hello", hello).Methods("GET")
	router.HandleFunc("/interconnect", interconnect).Methods("GET")

	router.HandleFunc("/heartbeat", Heartbeat)

	log.Println("starting")

	port := fmt.Sprintf(":%d", result.ServicePort)
	log.Fatal(http.ListenAndServe(port, router))
}
