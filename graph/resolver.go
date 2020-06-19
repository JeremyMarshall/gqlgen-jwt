package graph

import "github.com/JeremyMarshall/gqlgen-jwt/rbac/gorbac"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Rbac      gorbac.Rbac
	JwtSecret string
	Serialize string
}
