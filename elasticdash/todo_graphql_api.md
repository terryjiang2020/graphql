# GET /graphql (todo example) - Todo list API

## Implementation

### Route Definition and Handler
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/todo/main.go` (Lines 35-38)
```go
http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
    result := executeQuery(r.URL.Query().Get("query"), schema.TodoSchema)
    json.NewEncoder(w).Encode(result)
})
```

### Handler Implementation
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/todo/main.go` (Lines 23-32)
```go
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
```

### Todo Schema
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/todo/schema/schema.go`

#### Todo Model
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/todo/schema/schema.go` (Lines 11-15)
```go
type Todo struct {
    ID   string `json:"id"`
    Text string `json:"text"`
    Done bool   `json:"done"`
}
```

#### Schema Implementation (Partial)
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/todo/schema/schema.go`
```go
var TodoType = graphql.NewObject(graphql.ObjectConfig{
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

var TodoSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
    Query: graphql.NewObject(graphql.ObjectConfig{
        Name: "RootQuery",
        Fields: graphql.Fields{
            "todo": &graphql.Field{
                Type: TodoType,
                Args: graphql.FieldConfigArgument{
                    "id": &graphql.ArgumentConfig{
                        Type: graphql.String,
                    },
                },
                Resolve: func(params graphql.ResolveParams) (interface{}, error) {
                    idQuery, isOK := params.Args["id"].(string)
                    if isOK {
                        for _, todo := range TodoList {
                            if todo.ID == idQuery {
                                return todo, nil
                            }
                        }
                    }
                    return nil, nil
                },
            },
            "todoList": &graphql.Field{
                Type: graphql.NewList(TodoType),
                Resolve: func(params graphql.ResolveParams) (interface{}, error) {
                    return TodoList, nil
                },
            },
        },
    }),
    Mutation: graphql.NewObject(graphql.ObjectConfig{
        Name: "RootMutation",
        Fields: graphql.Fields{
            "createTodo": &graphql.Field{
                Type: TodoType,
                Args: graphql.FieldConfigArgument{
                    "text": &graphql.ArgumentConfig{
                        Type: graphql.NewNonNull(graphql.String),
                    },
                },
                Resolve: func(params graphql.ResolveParams) (interface{}, error) {
                    text, _ := params.Args["text"].(string)
                    todo := Todo{
                        ID:   randomID(),
                        Text: text,
                        Done: false,
                    }
                    TodoList = append(TodoList, todo)
                    return todo, nil
                },
            },
            "updateTodo": &graphql.Field{
                Type: TodoType,
                Args: graphql.FieldConfigArgument{
                    "id": &graphql.ArgumentConfig{
                        Type: graphql.NewNonNull(graphql.String),
                    },
                    "done": &graphql.ArgumentConfig{
                        Type: graphql.Boolean,
                    },
                    "text": &graphql.ArgumentConfig{
                        Type: graphql.String,
                    },
                },
                Resolve: func(params graphql.ResolveParams) (interface{}, error) {
                    id, _ := params.Args["id"].(string)
                    done, doneOk := params.Args["done"].(bool)
                    text, textOk := params.Args["text"].(string)
                    
                    // Find and update todo
                    for i, todo := range TodoList {
                        if todo.ID == id {
                            if doneOk {
                                TodoList[i].Done = done
                            }
                            if textOk {
                                TodoList[i].Text = text
                            }
                            return TodoList[i], nil
                        }
                    }
                    return nil, nil
                },
            },
        },
    }),
})
```

## Input Format

- **HTTP Method**: GET
- **Endpoint**: `/graphql`
- **Query Parameters**:
  - `query` (required, string): The GraphQL query to execute
  
Example GraphQL queries:

1. Get a single todo:
```
{todo(id:"b"){id,text,done}}
```

2. Create a todo:
```
mutation {createTodo(text:"My new todo"){id,text,done}}
```

3. Update a todo:
```
mutation {updateTodo(id:"a",done:true){id,text,done}}
```

4. List all todos:
```
{todoList{id,text,done}}
```

## Output Format

- **Content-Type**: application/json
- **Status Codes**:
  - 200: Success (implicit)
  - No explicit error status codes, but errors are returned in the response

- **Response Structure**:
  ```json
  {
    "data": {
      // Requested data based on the query
    },
    "errors": [
      // Array of errors if any occurred
    ]
  }
  ```

## Sample Input

Get a single todo:
```bash
curl -g 'http://localhost:8080/graphql?query={todo(id:"b"){id,text,done}}'
```

Create a new todo:
```bash
curl -g 'http://localhost:8080/graphql?query=mutation+_{createTodo(text:"My+new+todo"){id,text,done}}'
```

Update a todo:
```bash
curl -g 'http://localhost:8080/graphql?query=mutation+_{updateTodo(id:"a",done:true){id,text,done}}'
```

Get all todos:
```bash
curl -g 'http://localhost:8080/graphql?query={todoList{id,text,done}}'
```

## Sample Output

Get a single todo response:
```json
{
  "data": {
    "todo": {
      "id": "b",
      "text": "This is the most important",
      "done": false
    }
  }
}
```

Create a new todo response:
```json
{
  "data": {
    "createTodo": {
      "id": "d",
      "text": "My new todo",
      "done": false
    }
  }
}
```

Update a todo response:
```json
{
  "data": {
    "updateTodo": {
      "id": "a",
      "text": "A todo not to forget",
      "done": true
    }
  }
}
```

List all todos response:
```json
{
  "data": {
    "todoList": [
      {
        "id": "a",
        "text": "A todo not to forget",
        "done": false
      },
      {
        "id": "b",
        "text": "This is the most important",
        "done": false
      },
      {
        "id": "c",
        "text": "Please do this or else",
        "done": false
      }
    ]
  }
}
```