package types

import (
	"io"
)

type Role struct {
	Permissions []string `yaml:"permissions"`
	Parents     []string `yaml:"parents"`
}

type Rbac interface {
	RbacQuery
	RbacMutate
	Check(roles []string, permission string) bool
}

type RbacQuery interface {
	GetPermissions(name *string) ([]string, error)
	GetRoles(name *string) (map[string]Role, error)
}
type RbacMutate interface {
	UpsertRole(name *string, perms []*string, parents []*string) (Role, error)
	DeleteRole(name *string) (bool, error)
	DeletePermission(name *string, permission *string) (bool, error)
	Load() error
	Save(writer io.Writer) error
}
