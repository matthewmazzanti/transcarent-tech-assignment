package main

import (
	"net/http"
	"log"
	"encoding/json"
)

func main() {
	log.Println("Hello world!")
}

func getJson(url string) (interface{}, error) {
	httpRes, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer httpRes.Body.Close()

	var jsonRes interface{} = nil
	err = json.NewDecoder(httpRes.Body).Decode(&jsonRes)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return jsonRes, nil
}
