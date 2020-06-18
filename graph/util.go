package graph

import (
	"github.com/JeremyMarshall/gqlgen-jwt/graph/model"
	"github.com/JeremyMarshall/gqlgen-jwt/rbac"
)

const (
	JwtSecret     = "secret"
	Issuer        = "issuer"
	ExpiryMins    = 60
	JwtTokenField = "user"
	DefaultPort   = "8088"
	GorbacYaml    = "./all.yaml"
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
