package graph_test

import (
	"context"
	"github.com/JeremyMarshall/gqlgen-jwt/graph"
	"github.com/JeremyMarshall/gqlgen-jwt/graph/model"
	"github.com/JeremyMarshall/gqlgen-jwt/rbac/dummy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Schema.Resolvers", func() {
	var (
		resolver *graph.Resolver
	)
	BeforeEach(func() {
		resolver = &graph.Resolver{
			Rbac: &dummy.Dummy{},
		}
	})

	Describe("Jwt", func() {
		var (
			jwt string
			err error
		)
		Context("Can create jwt", func() {
			It("should succeed", func() {
				jwt, err = resolver.Mutation().CreateJwt(context.Background(), model.NewJwt{User: "test", Roles: []string{"role1", "role2"}})

				Expect(err).To(BeNil())
				Expect(jwt).ShouldNot(Equal(""))
			})
		})
		Context("Can decode jwt", func() {
			It("should succeed", func() {
				decode, err := resolver.Query().Jwt(context.Background(), jwt)

				Expect(err).To(BeNil())
				Expect(decode.User).To(Equal("test"))
			})
		})
	})
	Describe("Rbac", func() {
		Context("Can upsert valid role", func() {
			It("should succeed", func() {
				role, err := resolver.Mutation().UpsertRole(context.Background(), model.AddRole{})

				Expect(err).To(BeNil())
				Expect(role.Name).To(Equal(""))
			})
		})

		Context("Cannot upsert invalid role", func() {
			It("should fail", func() {
				_, err := resolver.Mutation().UpsertRole(
					context.Background(), 
					model.AddRole{Name: "error"},)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("Can delete role", func() {
			It("should succeed", func() {
				ok, err := resolver.Mutation().DeleteRole(context.Background(), model.DeleteRole{})

				Expect(err).To(BeNil())
				Expect(ok).To(BeTrue())
			})
		})

		Context("Cannot delete invalid role", func() {
			It("should fail", func() {
				_, err := resolver.Mutation().DeleteRole(
					context.Background(), 
					model.DeleteRole{
						Name: "error",
					})

				Expect(err).To(HaveOccurred())
			})
		})
		Context("Can delete permission", func() {
			It("should succeed", func() {
				ok, err := resolver.Mutation().DeletePermission(context.Background(), model.DeletePermission{})

				Expect(err).To(BeNil())
				Expect(ok).To(BeTrue())
			})
		})

		Context("Cannot delete invalid permission", func() {
			It("should fail", func() {
				_, err := resolver.Mutation().DeletePermission(
					context.Background(), 
					model.DeletePermission{
						Name: "error",
					})

				Expect(err).To(HaveOccurred())
			})
		})

		Context("Can get all permissions", func() {
			It("should succeed", func() {
				perms, err := resolver.Query().Permission(context.Background(), nil)

				Expect(err).To(BeNil())
				Expect(len(perms)).To(Equal(2))
			})
		})
		Context("Can get one permission", func() {
			It("should succeed", func() {
				perm := "perm1"
				perms, err := resolver.Query().Permission(context.Background(), &perm)

				Expect(err).To(BeNil())
				Expect(len(perms)).To(Equal(1))
			})
		})
		Context("Can't get invalid permission", func() {
			It("should succeed", func() {
				perm := "error"
				_, err := resolver.Query().Permission(context.Background(), &perm)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("Cannot delete invalid permission", func() {
			It("should fail", func() {
				_, err := resolver.Mutation().DeletePermission(
					context.Background(), 
					model.DeletePermission{
						Name: "error",
					})

				Expect(err).To(HaveOccurred())
			})
		})

		Context("Can get all roles", func() {
			It("should succeed", func() {
				roles, err := resolver.Query().Role(context.Background(), nil)

				Expect(err).To(BeNil())
				Expect(len(roles)).To(Equal(2))
			})
		})
		Context("Can get one role", func() {
			It("should succeed", func() {
				role := "role1"
				roles, err := resolver.Query().Role(context.Background(), &role)

				Expect(err).To(BeNil())
				Expect(len(roles)).To(Equal(1))
			})
		})
		Context("Can't get invalid role", func() {
			It("should succeed", func() {
				role := "error"
				_, err := resolver.Query().Role(context.Background(), &role)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
