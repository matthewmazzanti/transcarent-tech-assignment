package main

import (
	"net/http"
	"log"
	"encoding/json"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	user, err := getJson("https://jsonplaceholder.typicode.com/users/1")
	if err != nil {
		log.Println("Unable to fetch user")
		return
	}

	posts, err := getJson("https://jsonplaceholder.typicode.com/posts?userId=1")
	if err != nil {
		log.Println("Unable to fetch user")
		return
	}

	log.Println(user)
	log.Println(posts)
}

func getJson(url string) (interface{}, error) {
	httpRes, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer httpRes.Body.Close()

	var jsonRes interface{} = nil
	err = json.NewDecoder(httpRes.Body).Decode(&jsonRes)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return jsonRes, nil
}
