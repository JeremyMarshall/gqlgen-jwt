package graph

import (
	"context"
	"fmt"

	"github.com/JeremyMarshall/gqlgen-jwt/graph/model"
	"github.com/JeremyMarshall/gqlgen-jwt/rbac"
)

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
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeletePermission(ctx context.Context, input model.DeletePermission) (bool, error) {
	panic(fmt.Errorf("not implemented"))
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

func convertRole(k string, v rbac.Role) *model.Role {
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

	return r
}
