package dummy

import (
	"bytes"
	// "github.com/JeremyMarshall/gqlgen-jwt/rbac/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// "strings"
)

var _ = Describe("Rbac", func() {

	var (
		rbac *Dummy
		err  error
	)

	BeforeEach(func() {})

	Describe("Yaml", func() {

		Context("Valid role and permission", func() {
			It("should succeed", func() {
				Expect(rbac.Check([]string{"editor"}, "add-text")).To(BeTrue())
			})
		})
		Context("Invalid role and valid permission", func() {
			It("should fail", func() {
				Expect(rbac.Check([]string{"error", "invalid2"}, "add-text")).To(BeFalse())
			})
		})
		Context("Valid role and invalid permission", func() {
			It("should fail", func() {
				Expect(rbac.Check([]string{"editor", "photographer"}, "error")).To(BeFalse())
			})
		})
		Context("Get all roles", func() {
			It("should succeed", func() {
				ret, err := rbac.GetRoles(nil)
				Expect(err).To(BeNil())
				Expect(len(ret)).To(Equal(2))
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
				role := "error"
				_, err := rbac.GetRoles(&role)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("Get all permissions", func() {
			It("should succeed", func() {
				ret, err := rbac.GetPermissions(nil)
				Expect(err).To(BeNil())
				Expect(len(ret)).To(Equal(2))
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
				role := "error"
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
			})
		})
		Context("Add new role invalid parent", func() {
			It("should error", func() {
				r := "new"
				p := "new1"
				pa := "error"
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
			})
		})
		Context("Delete unknown role", func() {
			It("should fail", func() {
				r := "error"

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
			})
		})
		Context("Delete  permission on invalid role", func() {
			It("should fail", func() {
				r := "error"
				p := "new1"

				_, err := rbac.DeletePermission(&r, &p)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("Delete  invalid permission on role", func() {
			It("should fail", func() {
				r := "editor"
				p := "error"

				_, err := rbac.DeletePermission(&r, &p)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("Can write yaml", func() {
			It("should succeed", func() {
				buf := new(bytes.Buffer)
				rbac.Save(buf)

				Expect(err).To(BeNil())
			})
		})
	})

})
