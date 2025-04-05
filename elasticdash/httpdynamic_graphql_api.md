# GET /graphql (httpdynamic example) - Dynamic User API

## Implementation

### Route Definition and Handler
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/httpdynamic/main.go` (Lines 129-132)
```go
http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
    result := executeQuery(r.URL.Query().Get("query"), schema)
    json.NewEncoder(w).Encode(result)
})
```

### Handler Implementation
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/httpdynamic/main.go` (Lines 50-59)
```go
func executeQuery(query string, schema graphql.Schema) *graphql.Result {
    result := graphql.Do(graphql.Params{
        Schema:        schema,
        RequestString: query,
    })
    if len(result.Errors) > 0 {
        fmt.Printf("wrong result, unexpected errors: %v\n", result.Errors)
    }
    return result
}
```

### Schema Generation
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/httpdynamic/main.go` (Lines 61-115)
```go
func importJSONDataFromFile(fileName string) error {
    content, err := ioutil.ReadFile(fileName)
    if err != nil {
        return err
    }

    var data []map[string]interface{}

    err = json.Unmarshal(content, &data)
    if err != nil {
        return err
    }

    fields := make(graphql.Fields)
    args := make(graphql.FieldConfigArgument)
    for _, item := range data {
        for k := range item {
            fields[k] = &graphql.Field{
                Type: graphql.String,
            }
            args[k] = &graphql.ArgumentConfig{
                Type: graphql.String,
            }
        }
    }

    var userType = graphql.NewObject(
        graphql.ObjectConfig{
            Name:   "User",
            Fields: fields,
        },
    )

    var queryType = graphql.NewObject(
        graphql.ObjectConfig{
            Name: "Query",
            Fields: graphql.Fields{
                "user": &graphql.Field{
                    Type: userType,
                    Args: args,
                    Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                        return filterUser(data, p.Args), nil
                    },
                },
            },
        })

    schema, _ = graphql.NewSchema(
        graphql.SchemaConfig{
            Query: queryType,
        },
    )

    return nil
}
```

### User Filtering Function
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/httpdynamic/main.go` (Lines 36-48)
```go
func filterUser(data []map[string]interface{}, args map[string]interface{}) map[string]interface{} {
    for _, user := range data {
        for k, v := range args {
            if user[k] != v {
                goto nextuser
            }
            return user
        }

    nextuser:
    }
    return nil
}
```

### Dynamic Data Reloading
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/httpdynamic/main.go` (Lines 117-127)
```go
func reloadSchema(sig os.Signal) {
    fmt.Println("Reloading schema...")
    err := importJSONDataFromFile("data.json")
    if err != nil {
        fmt.Println("Error reloading schema:", err)
    } else {
        fmt.Println("Schema reloaded!")
    }
}
```

## Input Format

- **HTTP Method**: GET
- **Endpoint**: `/graphql`
- **Query Parameters**:
  - `query` (required, string): The GraphQL query to execute
  
Example GraphQL query:
```
{user(name:"Dan"){id,surname}}
```

Where:
- `user`: The query field
- Any property from the JSON data can be used as a filter argument
- Any property from the JSON data can be requested as a field

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
curl -g 'http://localhost:8080/graphql?query={user(name:"Dan"){id,surname}}'
```

## Sample Output

```json
{
  "data": {
    "user": {
      "id": "1",
      "surname": "Jones"
    }
  }
}
```

For a different field filter:
```bash
curl -g 'http://localhost:8080/graphql?query={user(id:"2"){name,surname}}'
```

Response:
```json
{
  "data": {
    "user": {
      "name": "Lee",
      "surname": "Brown"
    }
  }
}
```

Note: The schema is dynamically generated from the contents of data.json and can be reloaded by sending a SIGUSR1 signal to the process.