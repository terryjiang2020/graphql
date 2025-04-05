# GET /graphql (context example) - Context Demonstration API

## Implementation

### Route Definition
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/context/main.go` (Line 60)
```go
http.HandleFunc("/graphql", graphqlHandler)
```

### Handler Implementation
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/context/main.go` (Lines 42-57)
```go
func graphqlHandler(w http.ResponseWriter, r *http.Request) {
    user := struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    }{1, "cool user"}
    result := graphql.Do(graphql.Params{
        Schema:        Schema,
        RequestString: r.URL.Query().Get("query"),
        Context:       context.WithValue(context.Background(), "currentUser", user),
    })
    if len(result.Errors) > 0 {
        log.Printf("wrong result, unexpected errors: %v", result.Errors)
        return
    }
    json.NewEncoder(w).Encode(result)
}
```

### Schema Definition
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/context/main.go` (Lines 15-40)
```go
var userType = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "User",
        Fields: graphql.Fields{
            "id": &graphql.Field{
                Type: graphql.String,
            },
            "name": &graphql.Field{
                Type: graphql.String,
            },
        },
    },
)

var queryType = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "Query",
        Fields: graphql.Fields{
            "me": &graphql.Field{
                Type: userType,
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    return p.Context.Value("currentUser"), nil
                },
            },
        },
    })

var Schema, _ = graphql.NewSchema(
    graphql.SchemaConfig{
        Query: queryType,
    },
)
```

## Input Format

- **HTTP Method**: GET
- **Endpoint**: `/graphql`
- **Query Parameters**:
  - `query` (required, string): The GraphQL query to execute
  
Example GraphQL query:
```
{me{id,name}}
```

Where:
- `me`: The query field that returns the current user from context
- `id`, `name`: The fields to retrieve from the user object

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

```bash
curl -g 'http://localhost:8080/graphql?query={me{id,name}}'
```

## Sample Output

```json
{
  "data": {
    "me": {
      "id": "1",
      "name": "cool user"
    }
  }
}
```

Note: This example demonstrates how to use context to pass user information to resolvers. The "current user" is hardcoded in the handler and passed to the GraphQL execution context. In a real application, this would typically come from authentication middleware.