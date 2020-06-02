package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/JeremyMarshall/gql-jwt/graph/generated"
	"github.com/JeremyMarshall/gql-jwt/graph/model"
	jwt "github.com/dgrijalva/jwt-go"
)

const (
	JWT_SECRET  = "secret"
	ISSUER      = "issuer"
	EXPIRY_MINS = 5
)

func (r *mutationResolver) CreateJwt(ctx context.Context, input model.NewJwt) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":  input.User,
		"roles": input.Roles,

		"iss": ISSUER,
		"sub": "gqlgen properties",
		"aud": "gqlgen",
		"exp": time.Now().Add(time.Minute * EXPIRY_MINS).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		// iss	Issuer			Identifies principal that issued the JWT.
		// sub	Subject			Identifies the subject of the JWT.
		// aud	Audience		Identifies the recipients that the JWT is intended for. Each principal intended to process the JWT must identify itself with a value in the audience claim. If the principal processing the claim does not identify itself with a value in the aud claim when this claim is present, then the JWT must be rejected.
		// exp	Expiration Time	Identifies the expiration time on and after which the JWT must not be accepted for processing. The value must be a NumericDate:[9] either an integer or decimal, representing seconds past 1970-01-01 00:00:00Z.
		// nbf	Not Before		Identifies the time on which the JWT will start to be accepted for processing. The value must be a NumericDate.
		// iat	Issued at		Identifies the time at which the JWT was issued. The value must be a NumericDate.
		// jti	JWT ID			Case sensitive unique identifier of the token even among different issuers.
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

		ret := &model.Jwt{Properties: make([]*model.Property, 0), Roles: make([]string, 0)}

		for k, v := range claims {
			switch k {
			case "user":
				val := fmt.Sprint(v)
				ret.User = val
			case "roles":
				for _, r := range v.([]interface{}) {
					ret.Roles = append(ret.Roles, fmt.Sprint(r))
				}
			case "exp", "nbf", "iat":
				t := fmt.Sprint(time.Unix(int64(v.(float64)), 0))
				ret.Properties = append(ret.Properties, &model.Property{Name: k, Value: t})
			default:
				val := fmt.Sprint(v)
				ret.Properties = append(ret.Properties, &model.Property{Name: k, Value: val})
			}
		}
		return ret, nil
	}
	return nil, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

