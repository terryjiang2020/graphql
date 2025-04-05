# GET /graphql (http example) - User query API

## Implementation

### Route Definition and Handler
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/http/main.go` (Lines 90-93)
```go
http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
    result := executeQuery(r.URL.Query().Get("query"), schema)
    json.NewEncoder(w).Encode(result)
})
```

### Handler Implementation
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/http/main.go` (Lines 76-85)
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

### Data Model
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/http/main.go` (Lines 12-15)
```go
type user struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

### Schema Definition
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/http/main.go` (Lines 25-74)
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
            "user": &graphql.Field{
                Type: userType,
                Args: graphql.FieldConfigArgument{
                    "id": &graphql.ArgumentConfig{
                        Type: graphql.String,
                    },
                },
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    idQuery, isOK := p.Args["id"].(string)
                    if isOK {
                        return data[idQuery], nil
                    }
                    return nil, nil
                },
            },
        },
    })

var schema, _ = graphql.NewSchema(
    graphql.SchemaConfig{
        Query: queryType,
    },
)
```

### Helper Function for Loading Data
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/http/main.go` (Lines 101-114)
```go
func importJSONDataFromFile(fileName string, result interface{}) (isOK bool) {
    isOK = true
    content, err := ioutil.ReadFile(fileName)
    if err != nil {
        fmt.Print("Error:", err)
        isOK = false
        return
    }
    err = json.Unmarshal(content, result)
    if err != nil {
        isOK = false
        fmt.Print("Error:", err)
        return
    }
    return
}
```

## Input Format

- **HTTP Method**: GET
- **Endpoint**: `/graphql`
- **Query Parameters**:
  - `query` (required, string): The GraphQL query to execute
  
Example GraphQL query:
```
{user(id:"1"){name}}
```

Where:
- `user`: The query field
- `id`: Required parameter for the user query (string)
- `name`: The field to retrieve from the user object

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
curl -g 'http://localhost:8080/graphql?query={user(id:"1"){name}}'
```

## Sample Output

```json
{
  "data": {
    "user": {
      "name": "Dan"
    }
  }
}
```

For a more complex query like `{user(id:"1"){id,name}}`, the response would be:

```json
{
  "data": {
    "user": {
      "id": "1",
      "name": "Dan"
    }
  }
}
```