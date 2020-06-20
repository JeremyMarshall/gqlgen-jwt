package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/JeremyMarshall/gqlgen-jwt"
	"context"
	"net/http"
	"net/http/httptest"
)

func NewContextWithRequestID(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, "reqId", "1234")
}

var _ = Describe("Main", func() {
	var ()
	BeforeEach(func() {

	})

	Describe("jwt middleware", func() {

		Context("Can process bearer token", func() {
			It("should succeed", func() {

				xxx := func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						// var ctx context.Context
						// ctx = NewContextWithRequestID(ctx, r)
						// next.ServeHTTP(w, r.WithContext(ctx))
						AuthMiddleware(next, "1234")
					})
				}

				// create a handler to use as "next" which will verify the request
				nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					Expect(r.Context().Value("reqId")).To(Equal("1234"))
					// val := r.Context().Value("reqId")
					// if val == nil {
					// 	t.Error("reqId not present")
					// }
					// valStr, ok := val.(string)
					// if !ok {
					// 	t.Error("not string")
					// }
					// if valStr != "1234" {
					// 	t.Error("wrong reqId")
					// }
				})

				// create the handler to test, using our custom "next" handler
				handlerToTest := xxx(nextHandler)

				// create a mock request to use
				req := httptest.NewRequest("GET", "http://testing", nil)
				req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJncWxnZW4iLCJleHAiOjE1OTI1NDY5MTcsImlhdCI6MTU5MjU0MzMxNywiaXNzIjoiaXNzdWVyIiwibmJmIjoxNTkyNTQzMzE3LCJyb2xlcyI6WyJqd3QiLCJyYmFjLXJ3Il0sInN1YiI6ImdxbGdlbiBwcm9wZXJ0aWVzIiwidXNlciI6ImFhIn0.3q7y_NSfuVaJ4AKuCV0be3GXrbZhvL9RqZGMfrfWcBI")

				rec := httptest.NewRecorder()

				// call the handler using a mock response recorder (we'll not use that anyway)
				handlerToTest.ServeHTTP(rec, req)

				Expect(rec).To(BeNil())
				// Expect(jwt).ShouldNot(Equal(""))
			})
		})
	})
})
