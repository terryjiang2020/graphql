package schema

import (
	"errors"
	"math/rand"

	"github.com/graphql-go/graphql"
)

var TodoList []Todo
var UserList []User

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
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
		"isAdmin": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

// root mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
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
		/*
			curl -g 'http://localhost:8080/graphql?query=mutation+_{updateUserPassword(adminId:"admin1",userId:"user1",newPassword:"newpass123"){id,username,isAdmin}}'
		*/
		"updateUserPassword": &graphql.Field{
			Type:        userType,
			Description: "Update user password (admin only)",
			Args: graphql.FieldConfigArgument{
				"adminId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
					Description: "ID of the admin user performing the action",
				},
				"userId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
					Description: "ID of the user whose password will be changed",
				},
				"newPassword": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
					Description: "New password for the user",
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				adminId, _ := params.Args["adminId"].(string)
				userId, _ := params.Args["userId"].(string)
				newPassword, _ := params.Args["newPassword"].(string)
				
				// Verify admin privileges
				var adminUser User
				adminFound := false
				for _, user := range UserList {
					if user.ID == adminId {
						adminUser = user
						adminFound = true
						break
					}
				}
				
				if !adminFound {
					return nil, errors.New("admin user not found")
				}
				
				if !adminUser.IsAdmin {
					return nil, errors.New("insufficient privileges: user is not an admin")
				}
				
				// Find and update the target user
				var updatedUser User
				userFound := false
				for i := 0; i < len(UserList); i++ {
					if UserList[i].ID == userId {
						UserList[i].Password = newPassword
						updatedUser = UserList[i]
						userFound = true
						break
					}
				}
				
				if !userFound {
					return nil, errors.New("target user not found")
				}
				
				return updatedUser, nil
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
		/*
		   curl -g 'http://localhost:8080/graphql?query={user(id:"admin1"){id,username,isAdmin}}'
		*/
		"user": &graphql.Field{
			Type:        userType,
			Description: "Get single user",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				idQuery, isOK := params.Args["id"].(string)
				if isOK {
					for _, user := range UserList {
						if user.ID == idQuery {
							return user, nil
						}
					}
				}
				return nil, errors.New("user not found")
			},
		},
	},
})

// define schema, with our rootQuery and rootMutation
var TodoSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})