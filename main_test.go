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
	"net/http"
	"net/http/httptest"
	// "github.com/JeremyMarshall/gqlgen-jwt/rbac/dummy"
)

// func NewContextWithRequestID(ctx context.Context, r *http.Request) context.Context {
// 	return context.WithValue(ctx, "reqId", "1234")
// }

var _ = Describe("Main", func() {
	var (
		resolver *graph.Resolver
		token string
	)

	convertBody := func(input *bytes.Buffer) map[string]string {
		m := make(map[string]string)
		err := json.Unmarshal(input.Bytes(), &m)
		Expect(err).To(BeNil())
		return m
	}

	// nextHandlerNoJwt := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	val := r.Context().Value("reqId")
	//     Expect(val).To(BeNil())

	//     valStr, ok := val.(string)
	// 	Expect(ok).To(BeTrue())
	//     Expect(valStr).To(Equal("1234"))
	// })

	// nextHandlerInvalidJwt := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	val := r.Context().Value("reqId")
	//     Expect(val).To(BeNil())

	//     valStr, ok := val.(string)
	// 	Expect(ok).To(BeTrue())
	//     Expect(valStr).To(Equal("1234"))
	// })

	BeforeEach(func() {
		// get a new token as they expire
		var err error
		resolver = &graph.Resolver{
			JwtSecret: graph.JwtSecret,
		}
		token, err = resolver.Mutation().CreateJwt(context.Background(), model.NewJwt{User: "aa", Roles: []string{"jwt", "rbac-rw"}})
		fmt.Println(token)
		Expect(err).To(BeNil())
		Expect(token).NotTo(BeNil())
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
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
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
})
