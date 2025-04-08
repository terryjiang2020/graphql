// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"github.com/graphql-go/graphql"
// )

// func main() {
// 	// Schema
// 	fields := graphql.Fields{
// 		"hello": &graphql.Field{
// 			Type: graphql.String,
// 			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 				return "world", nil
// 			},
// 		},
// 	}
// 	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
// 	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
// 	schema, err := graphql.NewSchema(schemaConfig)
// 	if err != nil {
// 		log.Fatalf("failed to create new schema, error: %v", err)
// 	}

// 	// Query
// 	query := `
// 		{
// 			hello
// 		}
// 	`
// 	params := graphql.Params{Schema: schema, RequestString: query}
// 	r := graphql.Do(params)
// 	if len(r.Errors) > 0 {
// 		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
// 	}
// 	rJSON, _ := json.Marshal(r)
// 	fmt.Printf("%s \n", rJSON) // {“data”:{“hello”:”world”}}
// }

/*
Analyze this codebase and provide a list of all files that contain API endpoints. Ignore any code that has been commented.
Examine all source code and identify any files that define routes, controllers, or API handlers. 
Ignore any commented content.
Don't limit your search to just routes files, also look at controller files or any file that might contain API endpoint definitions.
Your response should include the file path for each file containing API endpoints.
Format your response as a simple list of file paths, one per line, with no additional text.
Save your response to a file named 'api_files.txt' in the current directory.

Analyze this codebase and provide a list of all files that contain usable, uncommented API endpoints. Ignore any code that has been commented.                                                                                                         │
Examine all source code and identify any files that define routes, controllers, or API handlers.                                                                                                                                                       │
Ignore any commented content.                                                                                                                                                                                                                          │
Don't limit your search to just routes files, also look at controller files or any file that might contain API endpoint definitions.                                                                                                                   │
Your response should include the file path for each file containing API endpoints.                                                                                                                                                                     │
Format your response as a simple list of file paths, one per line, with no additional text.                                                                                                                                                            │
Save your response to a file named 'api_files.txt' in the current directory. 
*/