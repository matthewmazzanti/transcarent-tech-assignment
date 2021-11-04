package main

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"strings"
	"strconv"
	"sync"
	"context"
)

// Define the data structure. Since the expected result has a very rigid
// structure, and since we are ingesting data, unpacking data to a concrete
// type allows us to run validation and sanity check against the input,
// and ensures that the rendered json is always valid.
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

	userPosts, status, err := getUserPosts(1)
	fmt.Println(userPosts)
	fmt.Println(status)
	fmt.Println(err)

	/*
	serverExit := &sync.WaitGroup{}
	runServer(serverExit)
	serverExit.Wait()
	*/
}

func errorStatus(status int) bool {
	return status < 200 || status >= 300
}

func writeError(w http.ResponseWriter, status int, msg string) {
	errRes := fmt.Sprintf(`{
  "code": %d,
  "error": "%s"
}`, status, msg)

	http.Error(w, errRes, status)
}

func runServer(wg *sync.WaitGroup) *http.Server {
	path := "/v1/user-posts/"

	handler := http.NewServeMux()
	handler.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			writeError(w, 404, "Not found")
		}

		subpath := strings.TrimPrefix(r.URL.Path, path)
		id, err := strconv.Atoi(subpath)
		if err != nil || id < 0 {
			writeError(w, 404, "Not Found")
			return
		}

		userPosts, status, err := getUserPosts(id)
		if status == 404 {
			writeError(w, 404, "Not found")
			return
		}

		if err != nil || errorStatus(status) {
			writeError(w, 500, "Something went wrong")
			return
		}

		userPostsJson, err := json.MarshalIndent(userPosts, "", "  ")
		if err != nil {
			writeError(w, 500, "Something went wrong")
			return
		}

		fmt.Fprintf(w, "%v", string(userPostsJson))
	})

	srv := &http.Server{
		Addr: ":8080",
		Handler: handler,
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("Server stopped due to error: %v", err)
		}
	}()

	return srv
}

// Request both the user and their posts, and stitch together into a UserPosts
// struct
func getUserPosts(id int) (*UserPosts, int, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resChan := make(chan interface{})
	var user *User = nil
	var posts []Post = nil

	// Run both gets in parallel
	go func() {
		resChan <- getUser(ctx, id)
	}()

	go func() {
		resChan <- getPosts(ctx, id)
	}()

	// Iterate until we have both the user and the post data. The requests
	// are defined such that a 200 status and no errors will always have
	// non-null values for these
	for user == nil || posts == nil {
		resIface := <-resChan
		switch res := resIface.(type) {
		case UserRes:
			// If we error out, cancel the in-flight request
			if errorStatus(res.status) {
				cancel()
				return nil, res.status, nil
			}
			if res.err != nil {
				cancel()
				return nil, res.status, res.err
			}
			user = res.user

		case PostsRes:
			// If we error out, cancel the in-flight request
			if errorStatus(res.status) {
				cancel()
				return nil, res.status, nil
			}
			if res.err != nil {
				cancel()
				return nil, res.status, res.err
			}
			posts = res.posts

		// Shouldnt happen, since this covers all cases of return types
		default:
			log.Fatalf("Unexpected type returned from resChan")
		}
	}

	close(resChan)

	return &UserPosts{
		Id: id,
		UserInfo: *user,
		Posts: posts,
	}, 200, nil
}

// Data structure to contain the state of request to user endpoint
type UserRes struct {
	user *User
	status int
	err error
}

// Make a get request to the user's endpoint, and validate the response
func getUser(ctx context.Context, id int) UserRes {
	url := fmt.Sprintf("%s/users/%d", baseUrl, id)

	res, status, err := getJson(ctx, url)
	if errorStatus(status) {
		return UserRes{ status: status }
	}
	if err != nil {
		return UserRes{ status: status, err: err }
	}

	user, err := parseUser(res)
	return UserRes{ user: &user, status: status, err: err }
}

// Unpack JSON data into the "User" data structure. If fields are missing or
// of the wrong type, return an error.
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

// Data structure to contain the state of request to posts endpoint
type PostsRes struct {
	posts []Post
	status int
	err error
}

// Make a get request to the posts endpoint with a userId filter, and validate
// the response
func getPosts(ctx context.Context, id int) PostsRes {
	url := fmt.Sprintf("%s/posts?userId=%d", baseUrl, id)
	res, status, err := getJson(ctx, url)
	if errorStatus(status) {
		return PostsRes{ posts: nil, status: status, err: nil }
	}
	if err != nil {
		return PostsRes{ posts: nil, status: status, err: err }
	}

	posts, err := parsePosts(res)
	return PostsRes{ posts: posts, status: status, err: err }
}

// Unpack multiple posts in a list
func parsePosts(res interface{}) ([]Post, error) {
	data, ok := res.([]interface{})
	if !ok {
		return nil, fmt.Errorf("non-list json")
	}

	// We know the length of the list, so we dont need to `append`
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

// Unpack JSON data into the "Post" data structure. If fields are missing or
// of the wrong type, return an error.
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

// Access a map[string]interface{} key, checking that the accessed value is an
// integer. Since JSON only defines float64s, this requires checking that the
// returned interface is a float, and that the value of that float is integral
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

	// If the value is an int, val == float64(int(val)). Save off the int
	// cast to return later
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

// Access a map[string]interface{} key, checking that the accessed value is a
// string
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

// Make a GET request to provided URL, parsing the response as json. Returns
// three values:
// An interface{} of the parsed json, if applicable
// An HTTP status code
// An error. This may be an error in the request or in the parsing of the json
func getJson(ctx context.Context, url string) (interface{}, int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	httpRes, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode < 200 || httpRes.StatusCode >= 300 {
		return nil, httpRes.StatusCode, nil
	}

	var jsonRes interface{} = nil
	err = json.NewDecoder(httpRes.Body).Decode(&jsonRes)
	if err != nil {
		return nil, httpRes.StatusCode, err
	}

	return jsonRes, httpRes.StatusCode, nil
}
