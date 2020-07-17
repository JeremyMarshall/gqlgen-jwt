package gorbac

import (
	"bytes"
	"github.com/JeremyMarshall/gqlgen-jwt/rbac/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("Rbac", func() {

	var (
		yaml string
		rbac *Rbac
		err  error
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
					Roles: map[string]types.Role{
						"role1": types.Role{Permissions: []string{"perm1"}},
						"role2": types.Role{
							Permissions: []string{"perm2"},
							Parents:     []string{"role1"},
						},
					}}
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
				Expect(rbac.Check("user", []string{"editor"}, "add-text")).To(BeTrue())
			})
		})
		Context("Invalid role and valid permission", func() {
			It("should fail", func() {
				Expect(rbac.Check("user", []string{"invalid", "invalid2"}, "add-text")).To(BeFalse())
			})
		})
		Context("Valid role and invalid permission", func() {
			It("should fail", func() {
				Expect(rbac.Check("user", []string{"editor", "photographer"}, "invalid")).To(BeFalse())
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
		Context("Get one permission", func() {
			It("should succeed", func() {
				p := "add-text"
				ret, err := rbac.GetPermissions(&p)
				Expect(err).To(BeNil())
				Expect(len(ret)).To(Equal(1))
			})
		})
		Context("Get invalid persmission", func() {
			It("should fail", func() {
				role := "invalid"
				_, err := rbac.GetPermissions(&role)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("Add new role with perms and parent", func() {
			It("should succeed", func() {
				r := "new"
				p := "new1"
				pa := "editor"
				perms := []*string{&p}
				parents := []*string{&pa}

				new, err := rbac.UpsertRole(&r, perms, parents)
				Expect(err).To(BeNil())
				Expect(len(new.Parents)).To(Equal(1))
				Expect(len(new.Permissions)).To(Equal(1))
				Expect(len(rbac.yamlAll.Permissions)).To(Equal(8))
				Expect(len(rbac.yamlAll.Roles)).To(Equal(4))

			})
		})
		Context("Add new role invalid parent", func() {
			It("should error", func() {
				r := "new"
				p := "new1"
				pa := "invalid"
				perms := []*string{&p}
				parents := []*string{&pa}

				_, err := rbac.UpsertRole(&r, perms, parents)
				Expect(err).To(HaveOccurred())

			})
		})
		Context("Delete a role", func() {
			It("should succeed", func() {
				r := "new"

				ok, err := rbac.DeleteRole(&r)
				Expect(err).To(BeNil())
				Expect(ok).To(BeTrue())
				// The permissions aren't reset until a save
				Expect(len(rbac.yamlAll.Permissions)).To(Equal(8))
				Expect(len(rbac.yamlAll.Roles)).To(Equal(3))

			})
		})
		Context("Delete unknown role", func() {
			It("should fail", func() {
				r := "invalid"

				_, err := rbac.DeleteRole(&r)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("Delete a permission", func() {
			It("should succeed", func() {
				r := "editor"
				p := "add-text"

				ok, err := rbac.DeletePermission(&r, &p)
				Expect(err).To(BeNil())
				Expect(ok).To(BeTrue())
				// The permissions aren't reset until a save
				Expect(len(rbac.yamlAll.Permissions)).To(Equal(8))
				Expect(len(rbac.yamlAll.Roles)).To(Equal(3))

			})
		})
		Context("Delete  permission on invalid role", func() {
			It("should fail", func() {
				r := "invalid"
				p := "new1"

				_, err := rbac.DeletePermission(&r, &p)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("Delete  invalid permission on role", func() {
			It("should fail", func() {
				r := "editor"
				p := "invalid"

				_, err := rbac.DeletePermission(&r, &p)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("Can write yaml", func() {
			It("should succeed", func() {
				buf := new(bytes.Buffer)
				rbac.Save(buf)

				Expect(err).To(BeNil())
				Expect(len(rbac.yamlAll.Permissions)).To(Equal(6))
				Expect(len(rbac.yamlAll.Roles)).To(Equal(3))
			})
		})
	})

})
