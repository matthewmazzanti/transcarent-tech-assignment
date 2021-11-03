package main

import (
	"net/http"
	"fmt"
	"html"
	"log"
	"encoding/json"
)

type UserPosts struct {
	Id int `json:"id"`
	UserInfo User `json:"userInfo"`
	Posts []Post `json:"posts"`
}

type User struct {
	Name string `json:"name"`
	Username string `json:"username"`
	Email string `json:"email"`
}

type Post struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
}

var baseUrl = "https://jsonplaceholder.typicode.com"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	userPosts, err := getUserPosts(10)
	if err != nil {
		log.Println(err)
		return
	}

	data, err := json.MarshalIndent(userPosts, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(data))
}

func runServer() {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUserPosts(id int) (*UserPosts, error) {
	user, err := getUser(id)
	if err != nil {
		return nil, err
	}

	posts, err := getPosts(id)
	if err != nil {
		return nil, err
	}

	return &UserPosts{
		Id: id,
		UserInfo: user,
		Posts: posts,
	}, nil
}

func getUser(id int) (User, error) {
	url := fmt.Sprintf("%s/users/%d", baseUrl, id)
	res, err := getJson(url)
	if err != nil { return User{}, err }

	return parseUser(res)
}

func parseUser(res interface{}) (User, error) {
	data, ok := res.(map[string]interface{})
	if !ok {
		return User{}, fmt.Errorf("returned non-object json")
	}

	name, err := indexStr(data, "name")
	if err != nil { return User{}, err }

	username, err := indexStr(data, "username")
	if err != nil { return User{}, err }

	email, err := indexStr(data, "email")
	if err != nil { return User{}, err }

	return User{
		Name: name,
		Username: username,
		Email: email,
	}, nil
}

func getPosts(id int) ([]Post, error) {
	url := fmt.Sprintf("%s/posts?userId=%d", baseUrl, id)
	res, err := getJson(url)
	if err != nil { return nil, err }

	return parsePosts(res)
}

func parsePosts(res interface{}) ([]Post, error) {
	data, ok := res.([]interface{})
	if !ok {
		return nil, fmt.Errorf("non-list json")
	}

	posts := make([]Post, len(data))
	for i, postIface := range data {
		post, err := parsePost(postIface)
		if err != nil {
			return nil, err
		}

		posts[i] = post
	}

	return posts, nil
}

func parsePost(res interface{}) (Post, error) {
	postData, ok := res.(map[string]interface{})
	if !ok {
		return Post{}, fmt.Errorf("non-object json")
	}

	postId, err := indexInt(postData, "id")
	if err != nil { return Post{}, err }

	title, err := indexStr(postData, "title")
	if err != nil { return Post{}, err }

	body, err := indexStr(postData, "body")
	if err != nil { return Post{}, err }

	return Post{
		Id: postId,
		Title: title,
		Body: body,
	}, nil

}

func indexInt(data map[string]interface{}, key string) (int, error) {
	valIface, ok := data[key]
	if !ok {
		return 0, fmt.Errorf("Data does not contain key: \"%s\"", key)
	}

	valFloat, ok := valIface.(float64)
	if !ok {
		return 0, fmt.Errorf(
			"Value at key \"%s\" was not a number",
			key,
		)
	}

	val := int(valFloat)
	if valFloat != float64(val) {
		return 0, fmt.Errorf(
			"Number at %f key \"%s\" was not an integer",
			valFloat,
			key,
		)
	}

	return val, nil
}

func indexStr(data map[string]interface{}, key string) (string, error) {
	valIface, ok := data[key]
	if !ok {
		return "", fmt.Errorf("Data does not contain key: %s", key)
	}

	val, ok := valIface.(string)
	if !ok {
		return "", fmt.Errorf(
			"Value at key \"%s\" was not a string",
			key,
		)
	}

	return val, nil
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
