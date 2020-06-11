package graph

import (
	"github.com/JeremyMarshall/gqlgen-jwt/graph/model"
	"github.com/JeremyMarshall/gqlgen-jwt/rbac"
)

const (
	JWT_SECRET      = "secret"
	ISSUER          = "issuer"
	EXPIRY_MINS     = 60
	JWT_TOKEN_FIELD = "user"
)

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
