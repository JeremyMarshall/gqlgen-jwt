package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/JeremyMarshall/gql-jwt/graph"
	"github.com/JeremyMarshall/gql-jwt/graph/generated"
	"github.com/JeremyMarshall/gql-jwt/rbac"
)

const defaultPort = "8088"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	rbac, err := rbac.NewRbac("aa")
	if err != nil {
		log.Fatal(err)
	}

	resolver := &graph.Resolver{
		Rbac: rbac,
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
