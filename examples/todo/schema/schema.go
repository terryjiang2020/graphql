package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/joho/godotenv"
)

var TodoList []Todo
var UserList []User

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Token         string `json:"token"`
	SessionToken  string `json:"sessionToken"`
	Email         string `json:"email"`
	GitlabID      int    `json:"gitlabID"`
	AvatarURL     string `json:"avatarURL"`
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// getGitLabToken gets an OAuth token from GitLab API using password grant type
func getGitLabToken(baseURL, clientID, clientSecret, username, password string) (string, error) {
	// Build the token request URL
	tokenURL := fmt.Sprintf("%s/oauth/token", baseURL)
	
	// Create the request body
	requestBody := map[string]string{
		"grant_type":    "password",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"username":      username,
		"password":      password,
		"scope":         "api read_user",
	}
	
	// Convert request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	
	// Make the HTTP request
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get token, status: %d, response: %s", resp.StatusCode, string(body))
	}
	
	// Parse the response
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}
	
	return tokenResponse.AccessToken, nil
}

// getGitLabUserInfo gets user information from GitLab API using an access token
func getGitLabUserInfo(baseURL, accessToken string) (map[string]interface{}, error) {
	// Build the user info request URL
	userInfoURL := fmt.Sprintf("%s/user", baseURL)
	
	// Create the HTTP request
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	
	// Set the authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	
	// Make the HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info, status: %d, response: %s", resp.StatusCode, string(body))
	}
	
	// Parse the response
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	
	return userInfo, nil
}

// define custom GraphQL ObjectType `todoType` for our Golang struct `Todo`
// Note that
// - the fields in our todoType maps with the json tags for the fields in our struct
// - the field type matches the field type in our struct
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
	},
})

// define custom GraphQL ObjectType for our User struct
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"username": &graphql.Field{
			Type: graphql.String,
		},
		"sessionToken": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"gitlabID": &graphql.Field{
			Type: graphql.Int,
		},
		"avatarURL": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// root mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		/*
			curl -g 'http://localhost:8080/graphql?query=mutation+_{gitlabLogin(username:"user",password:"pass"){id,username,sessionToken}}'
		*/
		"gitlabLogin": &graphql.Field{
			Type:        userType,
			Description: "Login with GitLab credentials",
			Args: graphql.FieldConfigArgument{
				"username": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"token": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// Load environment variables from .env file
				if err := godotenv.Load(); err != nil {
					return nil, fmt.Errorf("error loading .env file: %v", err)
				}

				// Get GitLab API credentials from environment variables
				gitlabAPIURL := os.Getenv("GITLAB_API_URL")
				if gitlabAPIURL == "" {
					gitlabAPIURL = "https://gitlab.com/api/v4" // Default GitLab API URL
				}
				clientID := os.Getenv("GITLAB_CLIENT_ID")
				clientSecret := os.Getenv("GITLAB_CLIENT_SECRET")
				
				if clientID == "" || clientSecret == "" {
					return nil, fmt.Errorf("missing GitLab credentials in environment variables")
				}

				username, _ := params.Args["username"].(string)
				password, passwordOK := params.Args["password"].(string)
				token, tokenOK := params.Args["token"].(string)

				// Check if user provided password or token
				if !passwordOK && !tokenOK {
					return nil, fmt.Errorf("either password or token must be provided")
				}

				var userInfo map[string]interface{}
				var accessToken string
				var err error

				if tokenOK {
					// Use provided personal access token
					accessToken = token
					// Get user info using the token
					userInfo, err = getGitLabUserInfo(gitlabAPIURL, accessToken)
				} else {
					// Request OAuth token using password grant type
					accessToken, err = getGitLabToken(gitlabAPIURL, clientID, clientSecret, username, password)
					if err != nil {
						return nil, fmt.Errorf("gitlab authentication failed: %v", err)
					}
					// Get user info using the obtained token
					userInfo, err = getGitLabUserInfo(gitlabAPIURL, accessToken)
				}

				if err != nil {
					return nil, fmt.Errorf("failed to get user info: %v", err)
				}

				// Generate a session token for our app
				sessionToken := RandStringRunes(32)
				
				// Extract user information from GitLab response
				gitlabID := int(userInfo["id"].(float64))
				gitlabUsername := userInfo["username"].(string)
				email := userInfo["email"].(string)
				avatarURL := userInfo["avatar_url"].(string)
				
				// Check if user already exists
				for i, user := range UserList {
					if user.GitlabID == gitlabID {
						// Update user information and session token
						UserList[i].SessionToken = sessionToken
						UserList[i].Token = accessToken
						UserList[i].Username = gitlabUsername
						UserList[i].Email = email
						UserList[i].AvatarURL = avatarURL
						return UserList[i], nil
					}
				}

				// Create new user
				newID := RandStringRunes(8)
				newUser := User{
					ID:           newID,
					Username:     gitlabUsername,
					Token:        accessToken,
					SessionToken: sessionToken,
					Email:        email,
					GitlabID:     gitlabID,
					AvatarURL:    avatarURL,
				}

				UserList = append(UserList, newUser)
				return newUser, nil
			},
		},
		
		/*
			curl -g 'http://localhost:8080/graphql?query=mutation+_{createTodo(text:"My+new+todo"){id,text,done}}'
		*/
		"createTodo": &graphql.Field{
			Type:        todoType, // the return type for this field
			Description: "Create new todo",
			Args: graphql.FieldConfigArgument{
				"text": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				// marshall and cast the argument value
				text, _ := params.Args["text"].(string)

				// figure out new id
				newID := RandStringRunes(8)

				// perform mutation operation here
				// for e.g. create a Todo and save to DB.
				newTodo := Todo{
					ID:   newID,
					Text: text,
					Done: false,
				}

				TodoList = append(TodoList, newTodo)

				// return the new Todo object that we supposedly save to DB
				// Note here that
				// - we are returning a `Todo` struct instance here
				// - we previously specified the return Type to be `todoType`
				// - `Todo` struct maps to `todoType`, as defined in `todoType` ObjectConfig`
				return newTodo, nil
			},
		},
		/*
			curl -g 'http://localhost:8080/graphql?query=mutation+_{updateTodo(id:"a",done:true){id,text,done}}'
		*/
		"updateTodo": &graphql.Field{
			Type:        todoType, // the return type for this field
			Description: "Update existing todo, mark it done or not done",
			Args: graphql.FieldConfigArgument{
				"done": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// marshall and cast the argument value
				done, _ := params.Args["done"].(bool)
				id, _ := params.Args["id"].(string)
				affectedTodo := Todo{}

				// Search list for todo with id and change the done variable
				for i := 0; i < len(TodoList); i++ {
					if TodoList[i].ID == id {
						TodoList[i].Done = done
						// Assign updated todo so we can return it
						affectedTodo = TodoList[i]
						break
					}
				}
				// Return affected todo
				return affectedTodo, nil
			},
		},
	},
})

// root query
// we just define a trivial example here, since root query is required.
// Test with curl
// curl -g 'http://localhost:8080/graphql?query={lastTodo{id,text,done}}'
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		/*
		   curl -g 'http://localhost:8080/graphql?query={currentUser(sessionToken:"token"){id,username}}'
		*/
		"currentUser": &graphql.Field{
			Type:        userType,
			Description: "Get current authenticated user",
			Args: graphql.FieldConfigArgument{
				"sessionToken": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				tokenQuery, _ := params.Args["sessionToken"].(string)
				
				// Search for user with matching session token
				for _, user := range UserList {
					if user.SessionToken == tokenQuery {
						return user, nil
					}
				}
				
				return nil, fmt.Errorf("invalid session token")
			},
		},

		/*
		   curl -g 'http://localhost:8080/graphql?query={todo(id:"b"){id,text,done}}'
		*/
		"todo": &graphql.Field{
			Type:        todoType,
			Description: "Get single todo",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				idQuery, isOK := params.Args["id"].(string)
				if isOK {
					// Search for el with id
					for _, todo := range TodoList {
						if todo.ID == idQuery {
							return todo, nil
						}
					}
				}

				return Todo{}, nil
			},
		},

		"lastTodo": &graphql.Field{
			Type:        todoType,
			Description: "Last todo added",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return TodoList[len(TodoList)-1], nil
			},
		},

		/*
		   curl -g 'http://localhost:8080/graphql?query={todoList{id,text,done}}'
		*/
		"todoList": &graphql.Field{
			Type:        graphql.NewList(todoType),
			Description: "List of todos",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return TodoList, nil
			},
		},
	},
})

// define schema, with our rootQuery and rootMutation
var TodoSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})