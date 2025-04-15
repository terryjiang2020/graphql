package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/examples/todo/schema"
)

func init() {
	// Initialize sample todos with user IDs
	todo1 := schema.Todo{ID: "a", Text: "A todo not to forget", Done: false, UserID: "sample-user-1"}
	todo2 := schema.Todo{ID: "b", Text: "This is the most important", Done: false, UserID: "sample-user-1"}
	todo3 := schema.Todo{ID: "c", Text: "Please do this or else", Done: false, UserID: "sample-user-2"}
	schema.TodoList = append(schema.TodoList, todo1, todo2, todo3)

	// Initialize sample users with hashed passwords
	pass1, _ := schema.HashPassword("password123")
	pass2, _ := schema.HashPassword("password456")
	user1 := schema.User{ID: "sample-user-1", Email: "user1@example.com", Password: pass1}
	user2 := schema.User{ID: "sample-user-2", Email: "user2@example.com", Password: pass2}
	schema.UserList = append(schema.UserList, user1, user2)

	rand.Seed(time.Now().UnixNano())
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Execute GraphQL query
		result := executeQuery(r.URL.Query().Get("query"), schema.TodoSchema)
		json.NewEncoder(w).Encode(result)
	})

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// Display some basic instructions
	fmt.Println("Now server is running on port 8080")
	fmt.Println("\nAuthentication Operations:")
	fmt.Println("Sign up: curl -g 'http://localhost:8080/graphql?query=mutation+_{signup(email:\"user@example.com\",password:\"password123\"){token,user{id,email}}}'")
	fmt.Println("Login: curl -g 'http://localhost:8080/graphql?query=mutation+_{login(email:\"user@example.com\",password:\"password123\"){token,user{id,email}}}'")

	fmt.Println("\nTodo Operations (require token):")
	fmt.Println("Get single todo: curl -g 'http://localhost:8080/graphql?query={todo(id:\"b\",token:\"YOUR_TOKEN\"){id,text,done}}'")
	fmt.Println("Create new todo: curl -g 'http://localhost:8080/graphql?query=mutation+_{createTodo(text:\"My+new+todo\",token:\"YOUR_TOKEN\"){id,text,done}}'")
	fmt.Println("Update todo: curl -g 'http://localhost:8080/graphql?query=mutation+_{updateTodo(id:\"a\",done:true,token:\"YOUR_TOKEN\"){id,text,done}}'")
	fmt.Println("Load todo list: curl -g 'http://localhost:8080/graphql?query={todoList(token:\"YOUR_TOKEN\"){id,text,done}}'")
	fmt.Println("\nAccess the web app via browser at 'http://localhost:8080'")

	http.ListenAndServe(":8080", nil)
}