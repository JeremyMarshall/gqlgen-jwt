package rbac

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/mikespook/gorbac"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func LoadYaml(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewDecoder(f).Decode(v)
}

func SaveYaml(filename string, v interface{}) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(v)
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
}

func NewRbac(yamlFile string) (*Rbac, error) {

	ret := &Rbac{
		rbac:        gorbac.New(),
		permissions: &gorbac.Permissions{},
		yamlAll:     &Serialize{},
	}

	if err := LoadYaml("all.yaml", ret.yamlAll); err != nil {
		return nil, err
	}

	for _, pid := range ret.yamlAll.Permissions {
		(*ret.permissions)[pid] = gorbac.NewStdPermission(pid)
	}

	for k, v := range ret.yamlAll.Roles {
		role := gorbac.NewStdRole(k)
		for _, pid := range v.Permissions {
			role.Assign((*ret.permissions)[pid])
		}
		ret.rbac.Add(role)
	}

	for k, v := range ret.yamlAll.Roles {
		if err := ret.rbac.SetParents(k, v.Parents); err != nil {
			return nil, err
		}
	}

	return ret, nil

	// // Check if `editor` can add text
	// if rbac.IsGranted("editor", permissions["add-text"], nil) {
	// 	log.Println("Editor can add text")
	// }
	// // Check if `chief-editor` can add text
	// if rbac.IsGranted("chief-editor", permissions["add-text"], nil) {
	// 	log.Println("Chief editor can add text")
	// }
	// // Check if `photographer` can add text
	// if !rbac.IsGranted("photographer", permissions["add-text"], nil) {
	// 	log.Println("Photographer can't add text")
	// }
	// // Check if `nobody` can add text
	// // `nobody` is not exist in goRBAC at the moment
	// if !rbac.IsGranted("nobody", permissions["read-text"], nil) {
	// 	log.Println("Nobody can't read text")
	// }
	// // Add `nobody` and assign `read-text` permission
	// nobody := gorbac.NewStdRole("nobody")
	// permissions["read-text"] = gorbac.NewStdPermission("read-text")
	// nobody.Assign(permissions["read-text"])
	// rbac.Add(nobody)

	// yamlAll.Roles["nobody"] = Role{
	// 	Permissions: []string{"read-text"},
	// }

	// // Check if `nobody` can read text again
	// if rbac.IsGranted("nobody", permissions["read-text"], nil) {
	// 	log.Println("Nobody can read text")
	// }

	// if err := SaveYaml("new-all.yaml", &yamlAll); err != nil {
	// 	log.Fatal(err)
	// }
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
