package graph

import (
	"context"
	"fmt"

	"github.com/JeremyMarshall/gqlgen-jwt/graph/model"
)

func (r *mutationResolver) UpsertRole(ctx context.Context, input model.AddRole) (*model.Role, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteRole(ctx context.Context, input model.DeleteRole) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeletePermission(ctx context.Context, input model.DeletePermission) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Permission(ctx context.Context, name *string) ([]*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Role(ctx context.Context, name *string) ([]*model.Role, error) {
	ret := make([]*model.Role, 0)
	roles, err := r.Rbac.GetRoles(name)
	if err != nil {
		return nil, err
	}

	for k, v := range roles {
		r := &model.Role{
			Name:        k,
			Permissions: make([]*string, 0),
			Parents:     make([]*string, 0),
		}

		for i := range v.Permissions {
			r.Permissions = append(r.Permissions, &v.Permissions[i])
		}

		for i := range v.Parents {
			r.Parents = append(r.Parents, &v.Parents[i])
		}

		ret = append(ret, r)
	}

	return ret, nil
}
