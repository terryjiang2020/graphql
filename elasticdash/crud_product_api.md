# GET /product (crud example) - Product CRUD API

## Implementation

### Route Definition and Handler
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/crud/main.go` (Lines 228-231)
```go
http.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
    result := executeQuery(r.URL.Query().Get("query"), schema)
    json.NewEncoder(w).Encode(result)
})
```

### Handler Implementation
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/crud/main.go` (Lines 216-225)
```go
func executeQuery(query string, schema graphql.Schema) *graphql.Result {
    result := graphql.Do(graphql.Params{
        Schema:        schema,
        RequestString: query,
    })
    if len(result.Errors) > 0 {
        fmt.Printf("errors: %v", result.Errors)
    }
    return result
}
```

### Product Model
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/crud/main.go` (Lines 14-19)
```go
type Product struct {
    ID    int64   `json:"id"`
    Name  string  `json:"name"`
    Info  string  `json:"info,omitempty"`
    Price float64 `json:"price"`
}
```

### Product Type Definition
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/crud/main.go` (Lines 42-60)
```go
var productType = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "Product",
        Fields: graphql.Fields{
            "id": &graphql.Field{
                Type: graphql.Int,
            },
            "name": &graphql.Field{
                Type: graphql.String,
            },
            "info": &graphql.Field{
                Type: graphql.String,
            },
            "price": &graphql.Field{
                Type: graphql.Float,
            },
        },
    },
)
```

### Query Type Definition
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/crud/main.go` (Lines 62-101)
```go
var queryType = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "Query",
        Fields: graphql.Fields{
            "product": &graphql.Field{
                Type:        productType,
                Description: "Get product by id",
                Args: graphql.FieldConfigArgument{
                    "id": &graphql.ArgumentConfig{
                        Type: graphql.Int,
                    },
                },
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    id, ok := p.Args["id"].(int)
                    if ok {
                        // Find product
                        for _, product := range products {
                            if int(product.ID) == id {
                                return product, nil
                            }
                        }
                    }
                    return nil, nil
                },
            },
            "list": &graphql.Field{
                Type:        graphql.NewList(productType),
                Description: "Get product list",
                Resolve: func(params graphql.ResolveParams) (interface{}, error) {
                    return products, nil
                },
            },
        },
    })
```

### Mutation Type Definition
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/crud/main.go` (Lines 103-208)
```go
var mutationType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Mutation",
    Fields: graphql.Fields{
        "create": &graphql.Field{
            Type:        productType,
            Description: "Create new product",
            Args: graphql.FieldConfigArgument{
                "name": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.String),
                },
                "info": &graphql.ArgumentConfig{
                    Type: graphql.String,
                },
                "price": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.Float),
                },
            },
            Resolve: func(params graphql.ResolveParams) (interface{}, error) {
                rand.Seed(time.Now().UnixNano())
                product := Product{
                    ID:    int64(rand.Intn(100000)), // generate random ID
                    Name:  params.Args["name"].(string),
                    Info:  params.Args["info"].(string),
                    Price: params.Args["price"].(float64),
                }
                products = append(products, product)
                return product, nil
            },
        },
        "update": &graphql.Field{
            Type:        productType,
            Description: "Update product by id",
            Args: graphql.FieldConfigArgument{
                "id": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.Int),
                },
                "name": &graphql.ArgumentConfig{
                    Type: graphql.String,
                },
                "info": &graphql.ArgumentConfig{
                    Type: graphql.String,
                },
                "price": &graphql.ArgumentConfig{
                    Type: graphql.Float,
                },
            },
            Resolve: func(params graphql.ResolveParams) (interface{}, error) {
                id, _ := params.Args["id"].(int)
                name, nameOk := params.Args["name"].(string)
                info, infoOk := params.Args["info"].(string)
                price, priceOk := params.Args["price"].(float64)
                
                // Update product
                for i, product := range products {
                    if int64(id) == product.ID {
                        if nameOk {
                            products[i].Name = name
                        }
                        if infoOk {
                            products[i].Info = info
                        }
                        if priceOk {
                            products[i].Price = price
                        }
                        return products[i], nil
                    }
                }
                return nil, nil
            },
        },
        "delete": &graphql.Field{
            Type:        productType,
            Description: "Delete product by id",
            Args: graphql.FieldConfigArgument{
                "id": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.Int),
                },
            },
            Resolve: func(params graphql.ResolveParams) (interface{}, error) {
                id, _ := params.Args["id"].(int)
                
                // Delete product
                for i, product := range products {
                    if int64(id) == product.ID {
                        // Return deleted product
                        deletedProduct := products[i]
                        
                        // Remove from product list
                        products = append(products[:i], products[i+1:]...)
                        
                        return deletedProduct, nil
                    }
                }
                return nil, nil
            },
        },
    },
})
```

### Schema Definition
File: `/Users/jiangjiahao/Documents/GitHub/graphql/examples/crud/main.go` (Lines 210-214)
```go
var schema, _ = graphql.NewSchema(
    graphql.SchemaConfig{
        Query:    queryType,
        Mutation: mutationType,
    },
)
```

## Input Format

- **HTTP Method**: GET
- **Endpoint**: `/product`
- **Query Parameters**:
  - `query` (required, string): The GraphQL query to execute
  
Example GraphQL queries:

1. Get a product by ID:
```
{product(id:1){name,info,price}}
```

2. List all products:
```
{list{id,name,info,price}}
```

3. Create a product:
```
mutation {create(name:"Inca Kola",info:"Inca Kola is a soft drink that was created in Peru in 1935",price:1.99){id,name,info,price}}
```

4. Update a product:
```
mutation {update(id:1,price:3.95){id,name,info,price}}
```

5. Delete a product:
```
mutation {delete(id:1){id,name,info,price}}
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

Get a product by ID:
```bash
curl -g 'http://localhost:8080/product?query={product(id:1){name,info,price}}'
```

List all products:
```bash
curl -g 'http://localhost:8080/product?query={list{id,name,info,price}}'
```

Create a product:
```bash
curl -g 'http://localhost:8080/product?query=mutation+_{create(name:"Inca+Kola",info:"Inca+Kola+is+a+soft+drink+that+was+created+in+Peru+in+1935",price:1.99){id,name,info,price}}'
```

Update a product:
```bash
curl -g 'http://localhost:8080/product?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}'
```

Delete a product:
```bash
curl -g 'http://localhost:8080/product?query=mutation+_{delete(id:1){id,name,info,price}}'
```

## Sample Output

Get a product response:
```json
{
  "data": {
    "product": {
      "name": "Chicha Morada",
      "info": "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level",
      "price": 7.99
    }
  }
}
```

List products response:
```json
{
  "data": {
    "list": [
      {
        "id": 1,
        "name": "Chicha Morada",
        "info": "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level",
        "price": 7.99
      },
      {
        "id": 2,
        "name": "Pisco Sour",
        "info": "Pisco Sour is an alcoholic cocktail of Peruvian origin that is traditional to Peru and Chile",
        "price": 9.95
      }
    ]
  }
}
```

Create product response:
```json
{
  "data": {
    "create": {
      "id": 12345,
      "name": "Inca Kola",
      "info": "Inca Kola is a soft drink that was created in Peru in 1935",
      "price": 1.99
    }
  }
}
```

Update product response:
```json
{
  "data": {
    "update": {
      "id": 1,
      "name": "Chicha Morada",
      "info": "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level",
      "price": 3.95
    }
  }
}
```

Delete product response:
```json
{
  "data": {
    "delete": {
      "id": 1,
      "name": "Chicha Morada",
      "info": "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level",
      "price": 3.95
    }
  }
}
```