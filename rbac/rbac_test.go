package rbac

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
	"bytes"
)

var _ = Describe("Rbac", func() {

    var (
		yaml  string
		rbac *Rbac
		err error
    )

    BeforeEach(func() {
		yaml = `
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
	})

	Describe("Yaml", func() {
		Context("Can read yaml", func() {
			It("should succeed", func() {
				s := &Serialize{}
				err = LoadYaml(strings.NewReader(yaml), s)

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
				err = SaveYaml(buf, s)

				Expect(err).To(BeNil())
			})
		})
		Context("Can create Rbac struct", func() {
			It("should succeed", func() {
				rbac, err = NewRbac(strings.NewReader(yaml))

				Expect(err).To(BeNil())
				Expect(len(rbac.yamlAll.Permissions)).To(Equal(7))
				Expect(len(rbac.yamlAll.Roles)).To(Equal(3))
			})
		})		
		Context("Valid role and permission", func() {
			It("should succeed", func() {
				Expect(rbac.Check([]string{"editor"}, "add-text")).To(Equal(true))
			})
		})	
		Context("Invalid role and valid permission", func() {
			It("should fail", func() {
				Expect(rbac.Check([]string{"invalid", "invalid2"}, "add-text")).To(Equal(false))
			})
		})	
		Context("Valid role and invalid permission", func() {
			It("should fail", func() {
				Expect(rbac.Check([]string{"editor", "photographer"}, "invalid")).To(Equal(false))
			})
		})
		Context("Get all roles", func() {
			It("should succeed", func() {
				ret, err := rbac.GetRoles(nil)
				Expect(err).To(BeNil())
				Expect(len(ret)).To(Equal(3))
			})
		})
		Context("Get existing role", func() {
			It("should succeed", func() {
				role := "editor"
				ret, err := rbac.GetRoles(&role)
				Expect(err).To(BeNil())
				Expect(len(ret)).To(Equal(1))
			})
		})
		Context("Get invalid role", func() {
			It("should fail", func() {
				role := "invalid"
				_, err := rbac.GetRoles(&role)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("Get all permissions", func() {
			It("should succeed", func() {
				ret, err := rbac.GetPermissions(nil)
				Expect(err).To(BeNil())
				Expect(len(ret)).To(Equal(7))
			})
		})
	})

})
