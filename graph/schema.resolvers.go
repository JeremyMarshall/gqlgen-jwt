package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/JeremyMarshall/gqlgen-jwt/graph/generated"
	"github.com/JeremyMarshall/gqlgen-jwt/graph/model"
)

func (r *mutationResolver) Jwt(ctx context.Context) (*model.JwtMutation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Rbac(ctx context.Context) (*model.RbacMutation, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Jwt(ctx context.Context) (*model.JwtQuery, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Rbac(ctx context.Context) (*model.RbacQuery, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
