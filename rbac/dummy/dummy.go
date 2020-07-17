package dummy

import (
	"fmt"
	"github.com/JeremyMarshall/gqlgen-jwt/rbac/types"
	"io"
)

type Dummy struct {
}

func (d *Dummy) GetPermissions(name *string) ([]string, error) {
	if name == nil {
		return []string{"Perm1", "Perm2"}, nil
	}
	if *name == "error" {
		return make([]string, 0), fmt.Errorf("Permission error")
	}
	return []string{*name}, nil
}

func (d *Dummy) GetRoles(name *string) (map[string]types.Role, error) {
	if name == nil {
		return map[string]types.Role{"role1": {}, "role2": {}}, nil
	}
	if *name == "error" {
		return nil, fmt.Errorf("Role error")
	}
	return map[string]types.Role{*name: {}}, nil
}

func (d *Dummy) UpsertRole(name *string, perms []*string, parents []*string) (types.Role, error) {
	if name == nil {
		return types.Role{}, fmt.Errorf("Upsert error")
	}
	if *name == "error" {
		return types.Role{}, fmt.Errorf("Upsert error")
	}

	ret := types.Role{
		Permissions: make([]string, 0),
		Parents:     make([]string, 0),
	}

	for _, p := range perms {
		if *p == "error" {
			return types.Role{}, fmt.Errorf("Upsert error")
		}
		ret.Permissions = append(ret.Permissions, *p)
	}

	for _, p := range parents {
		if *p == "error" {
			return types.Role{}, fmt.Errorf("Upsert error")
		}
		ret.Parents = append(ret.Parents, *p)
	}

	return ret, nil
}

func (d *Dummy) DeleteRole(name *string) (bool, error) {
	if name == nil {
		return false, fmt.Errorf("Delete error")
	}
	if *name == "error" {
		return false, fmt.Errorf("Delete error")
	}
	return true, nil
}

func (d *Dummy) DeletePermission(name *string, permission *string) (bool, error) {
	if name == nil || permission == nil {
		return false, fmt.Errorf("Delete error")
	}
	if *name == "error" || *permission == "error" {
		return false, fmt.Errorf("Delete error")
	}
	return true, nil
}

func (d *Dummy) Load() error {
	return nil
}

func (d *Dummy) Save(writer io.Writer) error {
	return nil
}

func (d *Dummy) Check(user string, roles []string, permission string) bool {
	if permission == "error" {
		return false
	}
	if roles == nil {
		return false
	}
	if len(roles) == 0 {
		return false
	}
	for _, r := range roles {
		if r == "error" {
			return false
		}
	}
	return true
}

func (d *Dummy) CheckDomain(user string, roles []string, domain *string, permission string) bool {
	if domain == nil {
		return false
	}
	if *domain == "error" {
		return false
	}
	return d.Check(user, roles, fmt.Sprintf("%s-%s", *domain, permission))
}
