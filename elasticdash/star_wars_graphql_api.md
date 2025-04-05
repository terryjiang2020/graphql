# GET /graphql (star-wars example) - Star Wars API

## Implementation

### Route Definition and Handler
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/star-wars/main.go` (Lines 13-20)
```go
http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("query")
    result := graphql.Do(graphql.Params{
        Schema:        testutil.StarWarsSchema,
        RequestString: query,
    })
    json.NewEncoder(w).Encode(result)
})
```

### Schema Implementation
This example uses a predefined Star Wars schema from the testutil package:
File: `/Users/jiangjiahao/Documents/GitHub/graphql/testutil/testutil.go`

The schema contains various Star Wars entities such as:
- Human
- Droid
- Episode
- Character interfaces
- Relationships between characters
- Query fields for heroes, humans, and droids

Note: The full schema implementation is quite extensive and available in the testutil package.

## Input Format

- **HTTP Method**: GET
- **Endpoint**: `/graphql`
- **Query Parameters**:
  - `query` (required, string): The GraphQL query to execute
  
Example GraphQL queries:

1. Get the hero:
```
{hero{name}}
```

2. Get a specific human by ID:
```
{human(id:"1000"){name,friends{name}}}
```

3. Get a specific droid by ID:
```
{droid(id:"2000"){name,primaryFunction}}
```

4. Query with episode argument:
```
{hero(episode:EMPIRE){name,friends{name}}}
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

Get the hero:
```bash
curl -g 'http://localhost:8080/graphql?query={hero{name}}'
```

Get a human with friends:
```bash
curl -g 'http://localhost:8080/graphql?query={human(id:"1000"){name,friends{name}}}'
```

Get a droid with primary function:
```bash
curl -g 'http://localhost:8080/graphql?query={droid(id:"2001"){name,primaryFunction}}'
```

Query hero by episode:
```bash
curl -g 'http://localhost:8080/graphql?query={hero(episode:EMPIRE){name,appearsIn}}'
```

## Sample Output

Hero response:
```json
{
  "data": {
    "hero": {
      "name": "R2-D2"
    }
  }
}
```

Human with friends response:
```json
{
  "data": {
    "human": {
      "name": "Luke Skywalker",
      "friends": [
        {
          "name": "Han Solo"
        },
        {
          "name": "Leia Organa"
        },
        {
          "name": "C-3PO"
        },
        {
          "name": "R2-D2"
        }
      ]
    }
  }
}
```

Droid with primary function response:
```json
{
  "data": {
    "droid": {
      "name": "C-3PO",
      "primaryFunction": "Protocol"
    }
  }
}
```

Hero by episode response:
```json
{
  "data": {
    "hero": {
      "name": "Luke Skywalker",
      "appearsIn": [
        "NEWHOPE",
        "EMPIRE",
        "JEDI"
      ]
    }
  }
}
```