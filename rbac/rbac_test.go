package rbac_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/JeremyMarshall/gqlgen-jwt/rbac"
	"strings"
	"bytes"
)

var _ = Describe("Rbac", func() {
	yaml := `
permissions: 
- add-text
- edit-text
- insert-photo
- add-photo
- edit-photo
- del-text
- del-photo
roles: 
 chief-editor: 
  parents: 
  - editor
  - photographer
  permissions: 
  - del-text
  - del-photo
 editor: 
  permissions: 
  - add-text
  - edit-text
  - insert-photo
 photographer: 
  permissions: 
  - add-photo
  - edit-photo`

	Describe("Yaml", func() {
		Context("Can read yaml", func() {
			It("should succeed", func() {
				s := &Serialize{}
				err := LoadYaml(strings.NewReader(yaml), s)

				Expect(err).To(BeNil())
				Expect(len(s.Permissions)).To(Equal(7))
				Expect(len(s.Roles)).To(Equal(3))
			})
		})
		Context("Can write yaml", func() {
			It("should succeed", func() {
				s := &Serialize{
					Permissions: []string{"perm1", "perm2"},
					Roles: map[string]Role{
						"role1": Role{Permissions: []string{"perm1"}},
						"role2": Role{
							Permissions: []string{"perm2"},
							Parents: []string{"role1"},
					},
				},}
				buf := new(bytes.Buffer)
				err := SaveYaml(buf, s)

				Expect(err).To(BeNil())
			})
		})
	})

})
