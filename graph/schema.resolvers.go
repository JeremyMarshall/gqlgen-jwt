package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/JeremyMarshall/gqlgen-jwt/graph/generated"
	"github.com/JeremyMarshall/gqlgen-jwt/graph/model"
	jwt "github.com/dgrijalva/jwt-go"
)

func (r *mutationResolver) CreateJwt(ctx context.Context, input model.NewJwt) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":  input.User,
		"roles": input.Roles,

		"iss": Issuer,
		"sub": "gqlgen properties",
		"aud": "gqlgen",
		"exp": time.Now().Add(time.Minute * ExpiryMins).Unix(),
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
	tokenString, err := token.SignedString([]byte(r.JwtSecret))

	return tokenString, err
}

func (r *mutationResolver) UpsertRole(ctx context.Context, input model.AddRole) (*model.Role, error) {
	// If the role exists, update the permissions
	// If the role doesn't exist create it and add the permissions
	role, err := r.Rbac.UpsertRole(&input.Name, input.Permissions, input.Parents)
	if err != nil {
		return nil, err
	}
	return convertRole(input.Name, role), nil
}

func (r *mutationResolver) DeleteRole(ctx context.Context, input model.DeleteRole) (bool, error) {
	return r.Rbac.DeleteRole(&input.Name)
}

func (r *mutationResolver) DeletePermission(ctx context.Context, input model.DeletePermission) (bool, error) {
	return r.Rbac.DeletePermission(&input.Name, &input.Permission)
}

func (r *mutationResolver) Save(ctx context.Context) (bool, error) {
	f, err := os.Create(r.Serialize)
	defer f.Close()
	if err != nil {
		return false, err
	}
	err = r.Rbac.Save(f)
	return err == nil, err
}

func (r *mutationResolver) AddNewspaper(ctx context.Context, name string) (string, error) {
	return "Add Newspaper", nil
}

func (r *mutationResolver) DeleteNewspaper(ctx context.Context, name string) (bool, error) {
	return true, nil
}

func (r *mutationResolver) AddStaff(ctx context.Context, input model.ModStaff) (string, error) {
	return "Add Staff", nil
}

func (r *mutationResolver) AddStory(ctx context.Context, input model.AddStory) (string, error) {
	return "Add Story", nil
}

func (r *mutationResolver) AddPhoto(ctx context.Context, input model.AddPhoto) (string, error) {
	return "Add Photo", nil
}

func (r *mutationResolver) DeleteStaff(ctx context.Context, input model.ModStaff) (bool, error) {
	return true, nil
}

func (r *mutationResolver) DeleteStory(ctx context.Context, input model.DeleteMedia) (bool, error) {
	return true, nil
}

func (r *mutationResolver) DeletePhoto(ctx context.Context, input model.DeleteMedia) (bool, error) {
	return true, nil
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
		return []byte(r.JwtSecret), nil
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

func (r *queryResolver) Permission(ctx context.Context, name *string) ([]*string, error) {
	perm, err := r.Rbac.GetPermissions(name)
	if err != nil {
		return nil, err
	}
	ret := make([]*string, 0)
	for i := range perm {
		ret = append(ret, &perm[i])
	}
	return ret, nil
}

func (r *queryResolver) Role(ctx context.Context, name *string) ([]*model.Role, error) {
	ret := make([]*model.Role, 0)
	roles, err := r.Rbac.GetRoles(name)
	if err != nil {
		return nil, err
	}

	for k, v := range roles {
		ret = append(ret, convertRole(k, v))
	}

	return ret, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
