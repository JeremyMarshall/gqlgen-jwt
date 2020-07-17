package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/JeremyMarshall/gqlgen-jwt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"
	"github.com/JeremyMarshall/gqlgen-jwt/graph"
	"github.com/JeremyMarshall/gqlgen-jwt/graph/model"
	"github.com/JeremyMarshall/gqlgen-jwt/rbac/dummy"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Main", func() {
	var (
		resolver    *graph.Resolver
		tokenString string
		// token *jwt.Token
	)

	convertBody := func(input *bytes.Buffer) map[string]string {
		m := make(map[string]string)
		err := json.Unmarshal(input.Bytes(), &m)
		Expect(err).To(BeNil())
		return m
	}

	BeforeEach(func() {
		// get a new token as they expire
		var err error
		resolver = &graph.Resolver{
			JwtSecret: graph.JwtSecret,
		}
		tokenString, err = resolver.Mutation().CreateJwt(context.Background(), model.NewJwt{User: "aa", Roles: []string{"jwt", "rbac-rw"}})
		Expect(err).To(BeNil())
		Expect(tokenString).NotTo(BeNil())
	})

	Describe("jwt middleware", func() {

		Context("Can process bearer token", func() {
			It("should succeed", func() {
				next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					val := r.Context().Value("user")
					Expect(val).NotTo(BeNil())

					user := GetCurrentUser(r.Context())

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(user)
				})

				req, err := http.NewRequest("GET", "/query", nil)
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
				Expect(err).To(BeNil())

				// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
				rr := httptest.NewRecorder()
				handler := AuthMiddleware(next, graph.JwtSecret)

				// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
				// directly and pass in our Request and ResponseRecorder.
				handler.ServeHTTP(rr, req)

				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusOK))

				// Check the response body is what we expect.
				Expect(rr.Body).NotTo(BeNil())
			})
		})

		Context("error invalid bearer token", func() {
			It("should return error code", func() {
				next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// this shouldn't get called
					Expect(true).To(BeFalse())
				})

				req, err := http.NewRequest("GET", "/query", nil)
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "invalid"))
				Expect(err).To(BeNil())

				// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
				rr := httptest.NewRecorder()
				handler := AuthMiddleware(next, graph.JwtSecret)

				// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
				// directly and pass in our Request and ResponseRecorder.
				handler.ServeHTTP(rr, req)

				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusUnauthorized))

				// Check the response body is what we expect.
				Expect(convertBody(rr.Body)).To(Equal(map[string]string{"message": "JWT Auth token contains an invalid number of segments"}))
				Expect(rr.Body).NotTo(BeNil())
			})
		})

		Context("No bearer token", func() {
			It("should succeed", func() {
				next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					val := r.Context().Value("user")
					Expect(val).To(BeNil())

					w.WriteHeader(http.StatusOK)
				})

				req, err := http.NewRequest("GET", "/query", nil)
				Expect(err).To(BeNil())

				// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
				rr := httptest.NewRecorder()
				handler := AuthMiddleware(next, graph.JwtSecret)

				// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
				// directly and pass in our Request and ResponseRecorder.
				handler.ServeHTTP(rr, req)

				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("gql rbac middleware", func() {
		Context("Role fulfils permission", func() {
			It("should succeed", func() {
				rbac := &dummy.Dummy{}
				rbw := RbacMiddleware(rbac)

				next := func(ctx context.Context) (res interface{}, err error) {
					return true, nil
				}

				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return []byte(graph.JwtSecret), nil
				})
				Expect(err).To(BeNil())

				ctx := context.WithValue(context.Background(), graph.JwtTokenField, token)

				ok, err := rbw(ctx, nil, next, "RBAC_MUTATE")
				Expect(err).To(BeNil())
				Expect(ok).To(BeTrue())

			})
		})

		Context("Role doesn't fulfil permission", func() {
			It("should fail", func() {
				rbac := &dummy.Dummy{}
				rbw := RbacMiddleware(rbac)

				next := func(ctx context.Context) (res interface{}, err error) {
					return true, nil
				}

				tokenString2, err := resolver.Mutation().CreateJwt(context.Background(), model.NewJwt{User: "aa", Roles: []string{}})
				Expect(err).To(BeNil())
				Expect(tokenString).NotTo(BeNil())

				token, err := jwt.Parse(tokenString2, func(token *jwt.Token) (interface{}, error) {
					return []byte(graph.JwtSecret), nil
				})
				Expect(err).To(BeNil())

				ctx := context.WithValue(context.Background(), graph.JwtTokenField, token)

				_, err = rbw(ctx, nil, next, "RBAC_MUTATE")
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("gql rbac domain middleware", func() {
		Context("Role fulfils permission", func() {
			It("should succeed", func() {
				rbac := &dummy.Dummy{}
				rbw := RbacDomainMiddleware(rbac)

				next := func(ctx context.Context) (res interface{}, err error) {
					return true, nil
				}

				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return []byte(graph.JwtSecret), nil
				})
				Expect(err).To(BeNil())

				ctx := context.WithValue(context.Background(), graph.JwtTokenField, token)

				ok, err := rbw(ctx, map[string]interface{}{"newspaper": "test"}, next, "RBAC_MUTATE", model.DomainNewspaper)
				Expect(err).To(BeNil())
				Expect(ok).To(BeTrue())

			})
		})

		Context("Role doesn't fulfil permission", func() {
			It("should fail", func() {
				rbac := &dummy.Dummy{}
				rbw := RbacDomainMiddleware(rbac)

				next := func(ctx context.Context) (res interface{}, err error) {
					return true, nil
				}

				tokenString2, err := resolver.Mutation().CreateJwt(context.Background(), model.NewJwt{User: "aa", Roles: []string{}})
				Expect(err).To(BeNil())
				Expect(tokenString).NotTo(BeNil())

				token, err := jwt.Parse(tokenString2, func(token *jwt.Token) (interface{}, error) {
					return []byte(graph.JwtSecret), nil
				})
				Expect(err).To(BeNil())

				ctx := context.WithValue(context.Background(), graph.JwtTokenField, token)

				_, err = rbw(ctx, map[string]interface{}{"newspaper": "test"}, next, "RBAC_MUTATE", model.DomainNewspaper)
				Expect(err).To(HaveOccurred())
			})
		})
		Describe("options", func() {
			Context("load from defaults", func() {
				It("should use consts", func() {
					opts := NewOpts([]string{"dummy"})
					Expect(opts.Port).To(Equal(graph.DefaultPort))
					Expect(opts.JwtSecret).To(Equal(graph.JwtSecret))
					Expect(opts.GorbacYaml).To(Equal(graph.GorbacYaml))
				})
			})
		})
	})
})
