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
	"github.com/JeremyMarshall/gqlgen-jwt/rbac/gorbac"
	"github.com/JeremyMarshall/gqlgen-jwt/rbac/types"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/namsral/flag"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func AuthMiddleware(next http.Handler, secret string) http.Handler {
	if len(secret) == 0 {
		log.Fatal("HTTP server unable to start, expected an APP_KEY for JWT auth")
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		Debug:         true,
		// Set this to false if you always want a bearer token present
		CredentialsOptional: true,
		UserProperty:        graph.JwtTokenField,
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

type User struct {
	User  string
	Roles []string
}

func GetCurrentUser(ctx context.Context) *User {
	if rawToken := ctx.Value(graph.JwtTokenField); rawToken != nil {
		token := rawToken.(*jwt.Token)

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
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

type RbacMiddlewareFunc func(ctx context.Context, obj interface{}, next graphql.Resolver, rbac model.Rbac) (res interface{}, err error)
type RbacDomainMiddlewareFunc func(ctx context.Context, obj interface{}, next graphql.Resolver, rbac model.Rbac, domainFiled model.Domain) (res interface{}, err error)

func RbacMiddleware(rbacChecker types.Rbac) RbacMiddlewareFunc {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, rbac model.Rbac) (res interface{}, err error) {
		if !rbacChecker.Check(GetCurrentUser(ctx).Roles, rbac.String()) {
			// block calling the next resolver
			return nil, fmt.Errorf("Access denied")
		}

		// or let it pass through
		return next(ctx)
	}
}

func RbacDomainMiddleware(rbacChecker types.Rbac) RbacDomainMiddlewareFunc {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, rbac model.Rbac, domainString model.Domain) (res interface{}, err error) {

		if args, ok := obj.(map[string]interface{}); ok {
			if domain, ok := args[domainString.String()].(string); ok {
				if rbacChecker.CheckDomain(GetCurrentUser(ctx).Roles, &domain, rbac.String()) {
					return next(ctx)
				}
			}
		}
		return nil, fmt.Errorf("Access denied")
	}
}

type opts struct {
	Port       string
	GorbacYaml string
	JwtSecret  string
}

func NewOpts() *opts {
	o := &opts{}

	flag.StringVar(&o.Port, "port", graph.DefaultPort, "Port number")
	flag.StringVar(&o.JwtSecret, "jwtSecret", graph.JwtSecret, "JWT Secret")
	flag.StringVar(&o.GorbacYaml, "gorbacYaml", graph.GorbacYaml, "RBAC yaml")

	flag.Parse()

	return o
}

func main() {

	opts := NewOpts()

	f, err := os.Open(opts.GorbacYaml)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rbac, err := gorbac.NewRbac(f)
	if err != nil {
		log.Fatal(err)
	}

	resolver := &graph.Resolver{
		Rbac:      rbac,
		JwtSecret: opts.JwtSecret,
		Serialize: opts.GorbacYaml,
	}

	c := generated.Config{
		Resolvers: resolver,
		Directives: generated.DirectiveRoot{
			HasRbac:       RbacMiddleware(rbac),
			HasRbacDomain: RbacDomainMiddleware(rbac),
		},
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", AuthMiddleware(handlers.LoggingHandler(os.Stdout, srv), opts.JwtSecret))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", opts.Port)
	log.Fatal(http.ListenAndServe(":"+opts.Port, nil))
}
