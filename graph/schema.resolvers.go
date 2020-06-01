package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/JeremyMarshall/gql-jwt/graph/generated"
	"github.com/JeremyMarshall/gql-jwt/graph/model"

	"github.com/dgrijalva/jwt-go"
)

const (
	JWT_SECRET = "secret"
)

func (r *mutationResolver) CreateJwt(ctx context.Context, input model.NewJwt) (string, error) {
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"foo": "bar",
			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(JWT_SECRET))

		return tokenString, err
}

func (r *queryResolver) Jwt(ctx context.Context, token string) (*model.Jwt, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(JWT_SECRET), nil
	})

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {

		ret := &model.Jwt{User:"", Properties: make([]*string, 0)}

		for _,v := range claims {
			val :=  fmt.Sprint(v)
			ret.Properties = append(ret.Properties, &val )
		}
		return ret, nil
		// fmt.Println(claims["foo"], claims["nbf"])
	} else {
		return nil, err
	}
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
