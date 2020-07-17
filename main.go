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
	"github.com/JeremyMarshall/gqlgen-jwt/rbac/dummy"
	"github.com/JeremyMarshall/gqlgen-jwt/rbac/gorbac"
	"github.com/JeremyMarshall/gqlgen-jwt/rbac/types"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"

	// "github.com/namsral/flag"
	"github.com/integrii/flaggy"
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
		user := GetCurrentUser(ctx)
		if !rbacChecker.Check(user.User, user.Roles, rbac.String()) {
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
				user := GetCurrentUser(ctx)
				if rbacChecker.CheckDomain(user.User, user.Roles, &domain, rbac.String()) {
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

	SubGoRbac    *flaggy.Subcommand
	SubDummyRbac *flaggy.Subcommand
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func NewOpts(argv []string) *opts {
	o := &opts{
		Port:       GetEnv("PORT", graph.DefaultPort),
		JwtSecret:  GetEnv("JWTSECRET", graph.JwtSecret),
		GorbacYaml: GetEnv("GORBACTYAML", graph.GorbacYaml),
	}

	// Set your program's name and description.  These appear in help output.
	flaggy.SetName("gqlgen-jwt")
	flaggy.SetDescription("A little of JWT and gqlgen")

	// Add a flag to the main program (this will be available in all subcommands as well).
	flaggy.String(&o.Port, "p", "port", "The port to listen on")
	flaggy.String(&o.JwtSecret, "j", "jwt-secret", "The secret to seed JWT tokens")

	// Create any subcommands and set their parameters.
	o.SubGoRbac = flaggy.NewSubcommand("gorbac")
	o.SubGoRbac.Description = "Use gorbac for rbac"
	// Add a flag to the subcommand.
	o.SubGoRbac.String(&o.GorbacYaml, "y", "rbac-yaml", "Yaml file for gorbac source")

	// Create any subcommands and set their parameters.
	o.SubDummyRbac = flaggy.NewSubcommand("dummy")
	o.SubDummyRbac.Description = "Use dummy rbac"

	// Set the version and parse all inputs into variables.
	//   flaggy.SetVersion(version)
	flaggy.AttachSubcommand(o.SubGoRbac, 1)
	flaggy.AttachSubcommand(o.SubDummyRbac, 1)
	flaggy.ParseArgs(argv)

	if !(o.SubGoRbac.Used || o.SubDummyRbac.Used) {
		flaggy.ShowHelpAndExit("Please supply a subcommand")
	}

	return o
}

func main() {

	opts := NewOpts(os.Args[1:])

	var rbac types.Rbac

	if opts.SubGoRbac.Used {

		f, err := os.Open(opts.GorbacYaml)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		rbac, err = gorbac.NewRbac(f)
		if err != nil {
			log.Fatal(err)
		}
	} else if opts.SubDummyRbac.Used {
		rbac = &dummy.Dummy{}
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
