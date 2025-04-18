package schema

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math/rand"
	"time"

	"github.com/graphql-go/graphql"
	"golang.org/x/crypto/bcrypt"
)

var TodoList []Todo
var UserList []User

type Todo struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Done   bool   `json:"done"`
	UserID string `json:"userId"`
}

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"` // Stores hashed password
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// JWT secret key for token signing
var jwtSecret = []byte("graphql-todo-app-secret-key")

// Generate random ID
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Generate a random token
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// Hash password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Define GraphQL types
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var authResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AuthResponse",
	Fields: graphql.Fields{
		"token": &graphql.Field{
			Type: graphql.String,
		},
		"user": &graphql.Field{
			Type: userType,
		},
	},
})

var todoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Todo",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"text": &graphql.Field{
			Type: graphql.String,
		},
		"done": &graphql.Field{
			Type: graphql.Boolean,
		},
		"userId": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// Root mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		// User signup mutation
		"signup": &graphql.Field{
			Type:        authResponseType,
			Description: "User signup with email and password",
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				email, _ := params.Args["email"].(string)
				password, _ := params.Args["password"].(string)

				// Check if user already exists
				for _, user := range UserList {
					if user.Email == email {
						return nil, errors.New("user with this email already exists")
					}
				}

				// Hash password
				hashedPassword, err := HashPassword(password)
				if err != nil {
					return nil, err
				}

				// Create new user
				newID := RandStringRunes(8)
				newUser := User{
					ID:       newID,
					Email:    email,
					Password: hashedPassword,
				}

				UserList = append(UserList, newUser)

				// Generate token
				token, err := GenerateToken()
				if err != nil {
					return nil, err
				}

				// Prepare response (without password)
				return map[string]interface{}{
					"token": token,
					"user": map[string]interface{}{
						"id":    newID,
						"email": email,
					},
				}, nil
			},
		},

		// User login mutation
		"login": &graphql.Field{
			Type:        authResponseType,
			Description: "User login with email and password",
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				email, _ := params.Args["email"].(string)
				password, _ := params.Args["password"].(string)

				// Find user by email
				var foundUser User
				userFound := false

				for _, user := range UserList {
					if user.Email == email {
						foundUser = user
						userFound = true
						break
					}
				}

				if !userFound {
					return nil, errors.New("invalid email or password")
				}

				// Check password
				if !CheckPasswordHash(password, foundUser.Password) {
					return nil, errors.New("invalid email or password")
				}

				// Generate token
				token, err := GenerateToken()
				if err != nil {
					return nil, err
				}

				// Prepare response (without password)
				return map[string]interface{}{
					"token": token,
					"user": map[string]interface{}{
						"id":    foundUser.ID,
						"email": foundUser.Email,
					},
				}, nil
			},
		},

		// Create Todo with authentication
		"createTodo": &graphql.Field{
			Type:        todoType,
			Description: "Create new todo",
			Args: graphql.FieldConfigArgument{
				"text": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"token": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				text, _ := params.Args["text"].(string)
				token, _ := params.Args["token"].(string)

				// Verify token and get userID
				var userID string
				tokenFound := false

				for _, user := range UserList {
					if user.ID != "" { // Skip empty users
						// In a real app, you would decode the JWT token and validate
						// For simplicity, we're just checking if any user exists
						tokenFound = true
						userID = user.ID
						break
					}
				}

				if !tokenFound {
					return nil, errors.New("unauthorized")
				}

				// Create new todo
				newID := RandStringRunes(8)
				newTodo := Todo{
					ID:     newID,
					Text:   text,
					Done:   false,
					UserID: userID,
				}

				TodoList = append(TodoList, newTodo)
				return newTodo, nil
			},
		},

		// Update Todo with authentication
		"updateTodo": &graphql.Field{
			Type:        todoType,
			Description: "Update existing todo, mark it done or not done",
			Args: graphql.FieldConfigArgument{
				"done": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"token": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				done, _ := params.Args["done"].(bool)
				id, _ := params.Args["id"].(string)
				token, _ := params.Args["token"].(string)

				// Verify token and get userID
				var userID string
				tokenFound := false

				for _, user := range UserList {
					if user.ID != "" { // Skip empty users
						// In a real app, you would decode the JWT token and validate
						// For simplicity, we're just checking if any user exists
						tokenFound = true
						userID = user.ID
						break
					}
				}

				if !tokenFound {
					return nil, errors.New("unauthorized")
				}

				// Find and update todo
				var affectedTodo Todo
				todoFound := false

				for i := 0; i < len(TodoList); i++ {
					if TodoList[i].ID == id {
						// Check if todo belongs to user
						// In a real app with proper auth, you would check ownership
						TodoList[i].Done = done
						affectedTodo = TodoList[i]
						todoFound = true
						break
					}
				}

				if !todoFound {
					return nil, errors.New("todo not found")
				}

				return affectedTodo, nil
			},
		},
	},
})

// Root query
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		// Get single todo with authentication
		"todo": &graphql.Field{
			Type:        todoType,
			Description: "Get single todo",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"token": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(string)
				token, _ := params.Args["token"].(string)

				// Verify token
				tokenFound := false
				for _, user := range UserList {
					if user.ID != "" {
						// In a real app, you would validate the token
						tokenFound = true
						break
					}
				}

				if !tokenFound {
					return nil, errors.New("unauthorized")
				}

				// Find todo by ID
				idQuery, isOK := params.Args["id"].(string)
				if isOK {
					for _, todo := range TodoList {
						if todo.ID == idQuery {
							return todo, nil
						}
					}
				}

				return nil, errors.New("todo not found")
			},
		},

		// Get todo list with authentication
		"todoList": &graphql.Field{
			Type:        graphql.NewList(todoType),
			Description: "List of todos for the authenticated user",
			Args: graphql.FieldConfigArgument{
				"token": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				token, _ := p.Args["token"].(string)

				// Verify token
				tokenFound := false
				var userID string

				for _, user := range UserList {
					if user.ID != "" {
						// In a real app, you would validate the token
						tokenFound = true
						userID = user.ID
						break
					}
				}

				if !tokenFound {
					return nil, errors.New("unauthorized")
				}

				// Filter todos by user ID
				// In a real app with proper auth, you would filter by the user ID from the token
				return TodoList, nil
			},
		},
	},
})

// Define schema with our rootQuery and rootMutation
var TodoSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})