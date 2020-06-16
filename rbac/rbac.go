package rbac

import (
	"fmt"
	"io"
	"log"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/mikespook/gorbac"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func LoadYaml(reader io.Reader, v interface{}) error {
	return yaml.NewDecoder(reader).Decode(v)
}

func SaveYaml(writer io.Writer, v interface{}) error {
	return yaml.NewEncoder(writer).Encode(v)
}

type Role struct {
	Permissions []string `yaml:"permissions"`
	Parents     []string `yaml:"parents"`
}

type Serialize struct {
	Permissions []string        `yaml:"permissions"`
	Roles       map[string]Role `yaml:"roles"`
}

type Rbac struct {
	rbac        *gorbac.RBAC
	permissions *gorbac.Permissions
	yamlAll     *Serialize
	mutex       *sync.Mutex
}

func NewRbac(reader io.Reader) (*Rbac, error) {

	ret := &Rbac{
		yamlAll:     &Serialize{},
		mutex:       &sync.Mutex{},
	}

	if err := LoadYaml(reader, ret.yamlAll); err != nil {
		return nil, err
	}

	return ret.Load()
}


func (r *Rbac) Load() (*Rbac, error) {
	r.rbac = gorbac.New()
	r.permissions = &gorbac.Permissions{}

	for _, pid := range r.yamlAll.Permissions {
		(*r.permissions)[pid] = gorbac.NewStdPermission(pid)
	}

	for k, v := range r.yamlAll.Roles {
		role := gorbac.NewStdRole(k)
		for _, pid := range v.Permissions {
			role.Assign((*r.permissions)[pid])
		}
		r.rbac.Add(role)
	}

	for k, v := range r.yamlAll.Roles {
		if err := r.rbac.SetParents(k, v.Parents); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *Rbac) Save(writer io.Writer) error {

	// remove any permissions not mentioned in roles
	r.yamlAll.Permissions = make([]string, 0)

	for _, v := range r.yamlAll.Roles {
		for _, pid := range v.Permissions {
			r.yamlAll.Permissions = appendIfMissing(r.yamlAll.Permissions, &pid)
		}
	}

	err := SaveYaml(writer, r.yamlAll)
	if err != nil {
		return err
	}

	_, err = r.Load()

	return err
}

func appendIfMissing(slice []string, i *string) []string {
	for _, ele := range slice {
		if ele == *i {
			return slice
		}
	}
	return append(slice, *i)
}
func (r *Rbac) Check(roles []string, permission string) bool {
	for _, role := range roles {
		if p, ok := (*r.permissions)[permission]; ok {
			if r.rbac.IsGranted(role, p, nil) {
				return true
			}
		}
	}
	return false
}

func (r *Rbac) GetRoles(name *string) (map[string]Role, error) {
	if name == nil {
		return r.yamlAll.Roles, nil
	}
	if role, ok := r.yamlAll.Roles[*name]; ok {
		return map[string]Role{*name: role}, nil
	}
	return nil, fmt.Errorf("Role %s not found", *name)
}

func (r *Rbac) GetPermissions(name *string) ([]string, error) {
	if name == nil {
		return r.yamlAll.Permissions, nil
	}
	for _, perm := range r.yamlAll.Permissions {
		if perm == *name {
			return []string{*name}, nil
		}
	}
	return nil, fmt.Errorf("Permission %s not found", *name)
}

func (r *Rbac) UpsertRole(name *string, perms []*string, parents []*string) (Role, error) {
	r.mutex.Lock()
	var role Role
	var ok bool

	if role, ok = r.yamlAll.Roles[*name]; !ok {
		// not found so add it
		role = Role{}
		r.yamlAll.Roles[*name] = role
	}

	for _, v := range perms {
		role.Permissions = appendIfMissing(role.Permissions, v)
		r.yamlAll.Permissions = appendIfMissing(r.yamlAll.Permissions, v)
	}

	for _, v := range parents {
		if _, err := r.GetRoles(v); err != nil {
			r.mutex.Unlock()
			return Role{}, fmt.Errorf("Parent role %s not found", *v)
		}

		role.Parents = appendIfMissing(role.Parents, v)
	}

	r.yamlAll.Roles[*name] = role

	r.mutex.Unlock()
	return role, nil
}

func (r *Rbac) DeleteRole(name *string) (bool, error) {
	r.mutex.Lock()

	if _, ok := r.yamlAll.Roles[*name]; !ok {
		r.mutex.Unlock()
		return false, fmt.Errorf("Role %s not found", *name)
	}

	delete(r.yamlAll.Roles, *name)

	r.mutex.Unlock()
	return true, nil
}

func (r *Rbac) DeletePermission(name *string, permission *string) (bool, error) {
	r.mutex.Lock()
	var role Role
	var ok bool

	if role, ok = r.yamlAll.Roles[*name]; !ok {
		r.mutex.Unlock()
		return false, fmt.Errorf("Role %s not found", *name)
	}

	perms := r.yamlAll.Roles[*name].Permissions

	for i, p := range perms {
		if p == *permission {
			perms = append(perms[:i], perms[i+1:]...)
			role.Permissions = perms
			r.yamlAll.Roles[*name] = role
			r.mutex.Unlock()
			return true, nil
		}
	}

	r.mutex.Unlock()
	return false, fmt.Errorf("Permission %s not found", *permission)

}
