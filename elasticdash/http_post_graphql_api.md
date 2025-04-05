# POST /graphql (http-post example) - Todo list API

## Implementation

### Route Definition and Handler
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/http-post/main.go` (Lines 19-35)
```go
http.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
    var p postData
    if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
        w.WriteHeader(400)
        return
    }
    result := graphql.Do(graphql.Params{
        Context:        req.Context(),
        Schema:         schema.TodoSchema,
        RequestString:  p.Query,
        VariableValues: p.Variables,
        OperationName:  p.Operation,
    })
    if err := json.NewEncoder(w).Encode(result); err != nil {
        fmt.Printf("could not write result to response: %s", err)
    }
})
```

### Request Payload Structure
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/http-post/main.go` (Lines 12-16)
```go
type postData struct {
    Query     string                 `json:"query"`
    Operation string                 `json:"operationName"`
    Variables map[string]interface{} `json:"variables"`
}
```

### Todo Model
File: Referenced from Todo Schema (Lines 11-15)
```go
type Todo struct {
    ID   string `json:"id"`
    Text string `json:"text"`
    Done bool   `json:"done"`
}
```

## Input Format

- **HTTP Method**: POST
- **Endpoint**: `/graphql`
- **Request Body**:
  - Format: JSON
  - Fields:
    - `query` (required, string): The GraphQL query to execute
    - `operationName` (optional, string): Name of the operation to execute
    - `variables` (optional, object): Variables to use with the query

**Example GraphQL Queries:**

1. Get a single todo:
```
{ todo(id:"b") { id text done } }
```

2. Create a todo:
```
mutation { createTodo(text:"My New todo") { id text done } }
```

3. Update a todo:
```
mutation { updateTodo(id:"a", done: true) { id text done } }
```

4. Get all todos:
```
{ todoList { id text done } }
```

## Output Format

- **Content-Type**: application/json
- **Status Codes**:
  - 200: Success
  - 400: Bad Request (when request body cannot be parsed)

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
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{ "query": "{ todo(id:\"b\") { id text done } }" }' \
  http://localhost:8080/graphql
```

Create a new todo:
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{ "query": "mutation { createTodo(text:\"My New todo\") { id text done } }" }' \
  http://localhost:8080/graphql
```

Update a todo:
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{ "query": "mutation { updateTodo(id:\"a\", done: true) { id text done } }" }' \
  http://localhost:8080/graphql
```

Get all todos:
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{ "query": "{ todoList { id text done } }" }' \
  http://localhost:8080/graphql
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
      "text": "My New todo",
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

Get all todos response:
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