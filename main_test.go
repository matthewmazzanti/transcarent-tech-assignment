package main

import (
    "testing"
    "encoding/json"
    "reflect"
    "sync"
    "context"
    "log"
    "fmt"
)

// User test data
var expUserStr = `{
  "id": 1,
  "name": "Leanne Graham",
  "username": "Bret",
  "email": "Sincere@april.biz",
  "address": {
    "street": "Kulas Light",
    "suite": "Apt. 556",
    "city": "Gwenborough",
    "zipcode": "92998-3874",
    "geo": {
      "lat": "-37.3159",
      "lng": "81.1496"
    }
  },
  "phone": "1-770-736-8031 x56442",
  "website": "hildegard.org",
  "company": {
    "name": "Romaguera-Crona",
    "catchPhrase": "Multi-layered client-server neural-net",
    "bs": "harness real-time e-markets"
  }
}`

var expUser = &User{
	Name: "Leanne Graham",
	Username: "Bret",
	Email: "Sincere@april.biz",
}

func expUserJson() interface{} {
	var res interface{}
	json.Unmarshal([]byte(expUserStr), &res)
	return res
}

// Posts test data
var expPostsStr = `[
  {
    "userId": 1,
    "id": 1,
    "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
    "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
  },
  {
    "userId": 1,
    "id": 2,
    "title": "qui est esse",
    "body": "est rerum tempore vitae\nsequi sint nihil reprehenderit dolor beatae ea dolores neque\nfugiat blanditiis voluptate porro vel nihil molestiae ut reiciendis\nqui aperiam non debitis possimus qui neque nisi nulla"
  },
  {
    "userId": 1,
    "id": 3,
    "title": "ea molestias quasi exercitationem repellat qui ipsa sit aut",
    "body": "et iusto sed quo iure\nvoluptatem occaecati omnis eligendi aut ad\nvoluptatem doloribus vel accusantium quis pariatur\nmolestiae porro eius odio et labore et velit aut"
  },
  {
    "userId": 1,
    "id": 4,
    "title": "eum et est occaecati",
    "body": "ullam et saepe reiciendis voluptatem adipisci\nsit amet autem assumenda provident rerum culpa\nquis hic commodi nesciunt rem tenetur doloremque ipsam iure\nquis sunt voluptatem rerum illo velit"
  },
  {
    "userId": 1,
    "id": 5,
    "title": "nesciunt quas odio",
    "body": "repudiandae veniam quaerat sunt sed\nalias aut fugiat sit autem sed est\nvoluptatem omnis possimus esse voluptatibus quis\nest aut tenetur dolor neque"
  },
  {
    "userId": 1,
    "id": 6,
    "title": "dolorem eum magni eos aperiam quia",
    "body": "ut aspernatur corporis harum nihil quis provident sequi\nmollitia nobis aliquid molestiae\nperspiciatis et ea nemo ab reprehenderit accusantium quas\nvoluptate dolores velit et doloremque molestiae"
  },
  {
    "userId": 1,
    "id": 7,
    "title": "magnam facilis autem",
    "body": "dolore placeat quibusdam ea quo vitae\nmagni quis enim qui quis quo nemo aut saepe\nquidem repellat excepturi ut quia\nsunt ut sequi eos ea sed quas"
  },
  {
    "userId": 1,
    "id": 8,
    "title": "dolorem dolore est ipsam",
    "body": "dignissimos aperiam dolorem qui eum\nfacilis quibusdam animi sint suscipit qui sint possimus cum\nquaerat magni maiores excepturi\nipsam ut commodi dolor voluptatum modi aut vitae"
  },
  {
    "userId": 1,
    "id": 9,
    "title": "nesciunt iure omnis dolorem tempora et accusantium",
    "body": "consectetur animi nesciunt iure dolore\nenim quia ad\nveniam autem ut quam aut nobis\net est aut quod aut provident voluptas autem voluptas"
  },
  {
    "userId": 1,
    "id": 10,
    "title": "optio molestias id quia eum",
    "body": "quo et expedita modi cum officia vel magni\ndoloribus qui repudiandae\nvero nisi sit\nquos veniam quod sed accusamus veritatis error"
  }
]`

var expPosts = []Post{
	{
		Id: 1,
		Title: "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
		Body: "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto",
	},
	{
		Id: 2,
		Title: "qui est esse",
		Body: "est rerum tempore vitae\nsequi sint nihil reprehenderit dolor beatae ea dolores neque\nfugiat blanditiis voluptate porro vel nihil molestiae ut reiciendis\nqui aperiam non debitis possimus qui neque nisi nulla",
	},
	{
		Id: 3,
		Title: "ea molestias quasi exercitationem repellat qui ipsa sit aut",
		Body: "et iusto sed quo iure\nvoluptatem occaecati omnis eligendi aut ad\nvoluptatem doloribus vel accusantium quis pariatur\nmolestiae porro eius odio et labore et velit aut",
	},
	{
		Id: 4,
		Title: "eum et est occaecati",
		Body: "ullam et saepe reiciendis voluptatem adipisci\nsit amet autem assumenda provident rerum culpa\nquis hic commodi nesciunt rem tenetur doloremque ipsam iure\nquis sunt voluptatem rerum illo velit",
	},
	{
		Id: 5,
		Title: "nesciunt quas odio",
		Body: "repudiandae veniam quaerat sunt sed\nalias aut fugiat sit autem sed est\nvoluptatem omnis possimus esse voluptatibus quis\nest aut tenetur dolor neque",
	},
	{
		Id: 6,
		Title: "dolorem eum magni eos aperiam quia",
		Body: "ut aspernatur corporis harum nihil quis provident sequi\nmollitia nobis aliquid molestiae\nperspiciatis et ea nemo ab reprehenderit accusantium quas\nvoluptate dolores velit et doloremque molestiae",
	},
	{
		Id: 7,
		Title: "magnam facilis autem",
		Body: "dolore placeat quibusdam ea quo vitae\nmagni quis enim qui quis quo nemo aut saepe\nquidem repellat excepturi ut quia\nsunt ut sequi eos ea sed quas",
	},
	{
		Id: 8,
		Title: "dolorem dolore est ipsam",
		Body: "dignissimos aperiam dolorem qui eum\nfacilis quibusdam animi sint suscipit qui sint possimus cum\nquaerat magni maiores excepturi\nipsam ut commodi dolor voluptatum modi aut vitae",
	},
	{
		Id: 9,
		Title: "nesciunt iure omnis dolorem tempora et accusantium",
		Body: "consectetur animi nesciunt iure dolore\nenim quia ad\nveniam autem ut quam aut nobis\net est aut quod aut provident voluptas autem voluptas",
	},
	{
		Id: 10,
		Title: "optio molestias id quia eum",
		Body: "quo et expedita modi cum officia vel magni\ndoloribus qui repudiandae\nvero nisi sit\nquos veniam quod sed accusamus veritatis error",
	},
}

func expPostsJson() interface{} {
	var res interface{}
	json.Unmarshal([]byte(expPostsStr), &res)
	return res
}


func TestGetJson(t *testing.T) {
	user, status, err := getJson(
		context.TODO(),
		"https://jsonplaceholder.typicode.com/users/1",
	)

	if err != nil {
		t.Fatalf("Non nil error in getJson: %v", err)
	}

	if errorStatus(status) {
		t.Fatalf("Got error status: %d", status)
	}

	exp := expUserJson()
	if !reflect.DeepEqual(exp, user) {
		t.Fatalf("\nExpected:\n%v\nGot:\n%v\n", exp, user)
	}

	posts, status, err := getJson(
		context.TODO(),
		"https://jsonplaceholder.typicode.com/posts?userId=1",
	)

	if err != nil {
		t.Fatalf("Non nil error in getJson: %v", err)
	}

	if errorStatus(status) {
		t.Fatalf("Got error status: %d", status)
	}

	exp = expPostsJson()
	if !reflect.DeepEqual(exp, posts) {
		t.Fatalf("\nExpected:\n%v\nGot:\n%v\n", exp, posts)
	}
}

var indexData = map[string]interface{}{
	"str": "string",
	"int": float64(1),
	"float": float64(1.5),
	"null": nil,
	"list": []interface{}{},
	"obj": map[string]interface{}{},
}

func TestIndexStr(t *testing.T) {
	res, err := indexStr(indexData, "str")
	if err != nil {
		t.Fatalf("Error while indexing a string %v", err)
	}

	if res != "string" {
		t.Fatalf("Got unexpected string: %s", res)
	}

	_, err = indexStr(indexData, "int")
	if err == nil { t.Fatalf("Did not get error indexing int") }

	_, err = indexStr(indexData, "float")
	if err == nil { t.Fatalf("Did not get error indexing float") }

	_, err = indexStr(indexData, "null")
	if err == nil { t.Fatalf("Did not get error indexing null") }

	_, err = indexStr(indexData, "list")
	if err == nil { t.Fatalf("Did not get error indexing list") }

	_, err = indexStr(indexData, "obj")
	if err == nil { t.Fatalf("Did not get error indexing object") }
}

func TestIndexInt(t *testing.T) {
	res, err := indexInt(indexData, "int")
	if err != nil {
		t.Fatalf("Error while indexing a int %v", err)
	}

	if res != 1 {
		t.Fatalf("Got unexpected int: %d", res)
	}

	_, err = indexInt(indexData, "str")
	if err == nil { t.Fatalf("Did not get error indexing str") }

	_, err = indexInt(indexData, "float")
	if err == nil { t.Fatalf("Did not get error indexing float") }

	_, err = indexInt(indexData, "null")
	if err == nil { t.Fatalf("Did not get error indexing null") }

	_, err = indexInt(indexData, "list")
	if err == nil { t.Fatalf("Did not get error indexing list") }

	_, err = indexInt(indexData, "obj")
	if err == nil { t.Fatalf("Did not get error indexing object") }

}

func TestParseUser(t *testing.T) {
	// Happy path
	data := map[string]interface{}{
		"name": "Leanne Graham",
		"username": "Bret",
		"email": "Sincere@april.biz",
		"foo": "bar", // extra field should be ignored
	}

	exp := User{
		Name: "Leanne Graham",
		Username: "Bret",
		Email: "Sincere@april.biz",
	}

	user, err := parseUser(data)
	if err != nil {
		t.Fatalf("Unexpected error converting user: %v", err)
	}


	if !reflect.DeepEqual(exp, user) {
		t.Fatalf("\nExpected:\n%v\nGot:\n%v\n", exp, user)
	}

	// Missing fields
	data = map[string]interface{}{
		"username": "Bret",
		"email": "Sincere@april.biz",
	}

	_, err = parseUser(data)
	if err == nil {
		t.Fatalf(
			"Did not get error parsing data with missing field: %v",
			data,
		)
	}

	// Wrong types
	data = map[string]interface{}{
		"name": "Leanne Graham",
		"username": 10,
		"email": nil,
	}

	_, err = parseUser(data)
	if err == nil {
		t.Fatalf(
			"Did not get error parsing data of wrong type: %v",
			data,
		)
	}
}

func TestParsePost(t *testing.T) {
	// Happy path
	data := map[string]interface{}{
		"id": float64(1),
		"title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
		"body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto",
		"foo": "bar", // extra field should be ignored
	}

	exp := Post{
		Id: 1,
		Title: "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
		Body: "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto",
	}

	post, err := parsePost(data)
	if err != nil {
		t.Fatalf("Unexpected error converting user: %v", err)
	}


	if !reflect.DeepEqual(exp, post) {
		t.Fatalf("\nExpected:\n%v\nGot:\n%v\n", exp, post)
	}

	// Missing fields
	data = map[string]interface{}{
		"title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
		"body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto",
	}

	_, err = parsePost(data)
	if err == nil {
		t.Fatalf(
			"Did not get error parsing data with missing field: %v",
			data,
		)
	}

	// Wrong types
	data = map[string]interface{}{
		"id": float64(1.2),
		"title": nil,
		"body": nil,
	}

	_, err = parsePost(data)
	if err == nil {
		t.Fatalf(
			"Did not get error parsing data of wrong type: %v",
			data,
		)
	}
}

func TestParsePosts(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{
			"userId": float64(1),
			"id": float64(1),
			"title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
			"body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto",
		},
		map[string]interface{}{
			"userId": float64(1),
			"id": float64(2),
			"title": "qui est esse",
			"body": "est rerum tempore vitae\nsequi sint nihil reprehenderit dolor beatae ea dolores neque\nfugiat blanditiis voluptate porro vel nihil molestiae ut reiciendis\nqui aperiam non debitis possimus qui neque nisi nulla",
		},
		map[string]interface{}{
			"userId": float64(1),
			"id": float64(3),
			"title": "ea molestias quasi exercitationem repellat qui ipsa sit aut",
			"body": "et iusto sed quo iure\nvoluptatem occaecati omnis eligendi aut ad\nvoluptatem doloribus vel accusantium quis pariatur\nmolestiae porro eius odio et labore et velit aut",
		},
	}

	exp := []Post{
		{
			Id: 1,
			Title: "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
			Body: "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto",
		},
		{
			Id: 2,
			Title: "qui est esse",
			Body: "est rerum tempore vitae\nsequi sint nihil reprehenderit dolor beatae ea dolores neque\nfugiat blanditiis voluptate porro vel nihil molestiae ut reiciendis\nqui aperiam non debitis possimus qui neque nisi nulla",
		},
		{
			Id: 3,
			Title: "ea molestias quasi exercitationem repellat qui ipsa sit aut",
			Body: "et iusto sed quo iure\nvoluptatem occaecati omnis eligendi aut ad\nvoluptatem doloribus vel accusantium quis pariatur\nmolestiae porro eius odio et labore et velit aut",
		},
	}

	posts, err := parsePosts(data)
	if err != nil {
		t.Fatalf("Unexpected error converting posts: %v", err)
	}

	if !reflect.DeepEqual(exp, posts) {
		t.Fatalf("\nExpected:\n%v\nGot:\n%v\n", exp, posts)
	}

	// Missing field
	data = []interface{}{
		map[string]interface{}{
			"title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
			"body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto",
		},
	}

	_, err = parsePosts(data)
	if err == nil {
		t.Fatalf(
			"Did not get error parsing data with missing field: %v",
			data,
		)
	}

	// Wrong types
	data = []interface{}{
		map[string]interface{}{
			"id": nil,
			"title": 1.5,
			"body": []interface{}{},
		},
	}

	_, err = parsePost(data)
	if err == nil {
		t.Fatalf(
			"Did not get error parsing data of wrong type: %v",
			data,
		)
	}
}

func TestGetUser(t *testing.T) {
	res := getUser(context.TODO(), 1)

	if res.err != nil {
		t.Fatalf("Unexpected getting user: %v", res.err)
	}

	if errorStatus(res.status) {
		t.Fatalf("Got error status: %d", res.status)
	}

	if !reflect.DeepEqual(expUser, res.user) {
		t.Fatalf("\nExpected:\n%v\nGot:\n%v\n", expUser, res.user)
	}
}

func TestGetPosts(t *testing.T) {
	res := getPosts(context.TODO(), 1)

	if res.err != nil {
		t.Fatalf("Unexpected getting posts: %v", res.err)
	}

	if errorStatus(res.status) {
		t.Fatalf("Got error status: %d", res.status)
	}


	if !reflect.DeepEqual(expPosts, res.posts) {
		t.Fatalf("\nExpected:\n%v\nGot:\n%v\n", expPosts, res.posts)
	}
}

func TestMarshal(t *testing.T) {
	userPosts := UserPosts{
		Id: 1,
		UserInfo: User{
			Name: "Leanne Graham",
			Username: "Bret",
			Email: "Sincere@april.biz",
		},
		Posts: []Post{
			{
				Id: 1,
				Title: "title of the user’s post",
				Body: "body of the user’s post",
			},
		},
	}

	var exp interface{}
	json.Unmarshal([]byte(`{
    "id": 1,
    "userInfo": {
        "name": "Leanne Graham",
        "username": "Bret",
        "email": "Sincere@april.biz"
    },
    "posts": [{
        "id": 1,
        "title": "title of the user’s post",
        "body": "body of the user’s post"
    }]
}`), &exp)

	userPostsJson, _ := json.Marshal(userPosts)
	var userPostsData interface{}
	json.Unmarshal(userPostsJson, &userPostsData)

	if !reflect.DeepEqual(exp, userPostsData) {
		t.Fatalf("\nExpected:\n%v\nGot:\n%v\n", exp, userPostsData)
	}
}

func TestServer(t *testing.T) {
	serverExit := &sync.WaitGroup{}
	srv := runServer(serverExit)

	for id := 1; id <= 10; id++ {
		url := fmt.Sprintf("http://localhost:8080/v1/user-posts/%d", id)
		res, status, err := getJson(context.TODO(), url)

		if err != nil {
			t.Fatalf("Failed to get user posts: %v", err)
		}

		if errorStatus(status) {
			t.Fatalf("Unexpected http error status: %d", status)
		}

		userPosts, status, err := getUserPosts(id)
		if err != nil || errorStatus(status) {
			log.Fatalf("Unable to get reference UserPosts")
		}

		userPostsJson, _ := json.Marshal(userPosts)
		var exp interface{}
		json.Unmarshal(userPostsJson, &exp)

		if !reflect.DeepEqual(exp, res) {
			t.Fatalf("\nExpected:\n%v\nGot:\n%v\n", exp, res)
		}
	}

	for id := 11; id <= 20; id++ {
		url := fmt.Sprintf("http://localhost:8080/v1/user-posts/%d", id)
		_, status, err := getJson(context.TODO(), url)
		if status != 404 {
			t.Fatalf("Unexpected http error status: %d", status)
		}

		if err != nil {
			t.Fatalf("Unexpected error getting user posts: %v", err)
		}
	}

	_, status, err := getJson(context.TODO(), "http://localhost:8080/v1/user-posts/-10")
	if status != 404 {
		t.Fatalf("Unexpected http error status: %d", status)
	}

	if err != nil {
		t.Fatalf("Unexpected error getting user posts: %v", err)
	}

	_, status, err = getJson(context.TODO(), "http://localhost:8080/v1/user-posts/asdfqwer")
	if status != 404 {
		t.Fatalf("Unexpected http error status: %d", status)
	}

	if err != nil {
		t.Fatalf("Unexpected error getting user posts: %v", err)
	}


	err = srv.Shutdown(context.TODO())
	if err != nil {
		log.Fatalf("Server failed to shut down")
	}

	serverExit.Wait()
}

func TestServerRemote404(t *testing.T) {
	serverExit := &sync.WaitGroup{}
	srv := runServer(serverExit)

	for id := 11; id <= 20; id++ {
		url := fmt.Sprintf("http://localhost:8080/v1/user-posts/%d", id)
		_, status, err := getJson(context.TODO(), url)
		if status != 404 {
			t.Fatalf("Unexpected http error status: %d", status)
		}

		if err != nil {
			t.Fatalf("Unexpected error getting user posts: %v", err)
		}
	}

	err := srv.Shutdown(context.TODO())
	if err != nil {
		log.Fatalf("Server failed to shut down")
	}

	serverExit.Wait()
}

func TestServerLocal404(t *testing.T) {
	serverExit := &sync.WaitGroup{}
	srv := runServer(serverExit)

	_, status, err := getJson(context.TODO(), "http://localhost:8080/v1/user-posts/-10")
	if status != 404 {
		t.Fatalf("Unexpected http error status: %d", status)
	}

	if err != nil {
		t.Fatalf("Unexpected error getting user posts: %v", err)
	}

	_, status, err = getJson(context.TODO(), "http://localhost:8080/v1/user-posts/asdfqwer")
	if status != 404 {
		t.Fatalf("Unexpected http error status: %d", status)
	}

	if err != nil {
		t.Fatalf("Unexpected error getting user posts: %v", err)
	}


	err = srv.Shutdown(context.TODO())
	if err != nil {
		log.Fatalf("Server failed to shut down")
	}

	serverExit.Wait()
}
