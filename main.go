package main

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// User represents a simple User struct
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// In-memory data store
var users = []User{
	{ID: "1", Name: "John Doe", Email: "john.doe@example.com"},
	{ID: "2", Name: "Jane Doe", Email: "jane.doe@example.com"},
}

// GraphQL schema definitions
func getSchema() (graphql.Schema, error) {
	// Define the User type
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":    &graphql.Field{Type: graphql.String},
			"name":  &graphql.Field{Type: graphql.String},
			"email": &graphql.Field{Type: graphql.String},
		},
	})

	// Query: Fetch all users
	queryFields := graphql.Fields{
		"users": &graphql.Field{
			Type: graphql.NewList(userType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return users, nil
			},
		},
		"user": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(string)
				for _, user := range users {
					if user.ID == id {
						return user, nil
					}
				}
				return nil, nil
			},
		},
	}

	// Mutation: Create, Update, and Delete
	mutationFields := graphql.Fields{
		"createUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.String},
				"name":  &graphql.ArgumentConfig{Type: graphql.String},
				"email": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user := User{
					ID:    p.Args["id"].(string),
					Name:  p.Args["name"].(string),
					Email: p.Args["email"].(string),
				}
				users = append(users, user)
				return user, nil
			},
		},
		"updateUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.String},
				"name":  &graphql.ArgumentConfig{Type: graphql.String},
				"email": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(string)
				name, nameOk := p.Args["name"].(string)
				email, emailOk := p.Args["email"].(string)

				for i, user := range users {
					if user.ID == id {
						if nameOk {
							users[i].Name = name
						}
						if emailOk {
							users[i].Email = email
						}
						return users[i], nil
					}
				}
				return nil, nil
			},
		},
		"deleteUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(string)
				for i, user := range users {
					if user.ID == id {
						deletedUser := users[i]
						users = append(users[:i], users[i+1:]...)
						return deletedUser, nil
					}
				}
				return nil, nil
			},
		},
	}

	// Create the schema
	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: queryFields}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{Name: "Mutation", Fields: mutationFields}),
	}
	return graphql.NewSchema(schemaConfig)
}

func main() {
	// Initialize the schema
	schema, err := getSchema()
	if err != nil {
		panic(err)
	}

	// Create a GraphQL handler
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true, // Enable GraphiQL interface
	})

	// Serve the GraphQL endpoint
	http.Handle("/graphql", h)

	// Start the server
	port := ":8080"
	println("GraphQL server running at http://localhost" + port + "/graphql")
	http.ListenAndServe(port, nil)
}
