package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/JeremyMarshall/gqlgen-jwt/graph"
	"github.com/JeremyMarshall/gqlgen-jwt/graph/generated"
	"github.com/JeremyMarshall/gqlgen-jwt/graph/model"
	"github.com/JeremyMarshall/gqlgen-jwt/rbac"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func AuthMiddleware(next http.Handler) http.Handler {
	if len(graph.JWT_SECRET) == 0 {
		log.Fatal("HTTP server unable to start, expected an APP_KEY for JWT auth")
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(graph.JWT_SECRET), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		Debug:         true,
		// Set this to false if you always want a bearer token present
		CredentialsOptional: true,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
			data := gqlerror.Error{
				Message: fmt.Sprintf("JWT Auth %s", err),
			}
			w.Header().Set("Content-Type", "application/json")
			// w.WriteHeader(http.StatusCreated)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(data)
			// w.Write([]byte(fmt.Sprintf("401 - %s", err)))
		},
	})
	return jwtMiddleware.Handler(next)
}

const defaultPort = "8088"

type User struct {
	User  string
	Roles []string
}

func (u *User) HasRbac(rbac model.Rbac) bool {
	return true
}

func getCurrentUser(ctx context.Context) *User {
	if aaa := ctx.Value("user"); aaa != nil {
		bbb := aaa.(*jwt.Token)

		if claims, ok := bbb.Claims.(jwt.MapClaims); ok && bbb.Valid {
			u := &User{
				User:  claims["user"].(string),
				Roles: make([]string, 0),
			}
			for _, r := range claims["roles"].([]interface{}) {
				u.Roles = append(u.Roles, fmt.Sprint(r))
			}
			return u
		}
	}
	return &User{}
}

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

	c := generated.Config{
		Resolvers: resolver,
		Directives: generated.DirectiveRoot{
			HasRbac: func(ctx context.Context, obj interface{}, next graphql.Resolver, rbac model.Rbac) (res interface{}, err error) {
				if !getCurrentUser(ctx).HasRbac(rbac) {
					// block calling the next resolver
					return nil, fmt.Errorf("Access denied")
				}

				// or let it pass through
				return next(ctx)
			},
		},
	}
	// c.Directives.HasRbac = func(ctx context.Context, obj interface{}, next graphql.Resolver, rbac model.Rbac) (res interface{}, err error) {
	// 	if !getCurrentUser(ctx).HasRbac(rbac) {
	// 		// block calling the next resolver
	// 		return nil, fmt.Errorf("Access denied")
	// 	}

	// 	// or let it pass through
	// 	return next(ctx)
	// }

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", AuthMiddleware(handlers.LoggingHandler(os.Stdout, srv)))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
