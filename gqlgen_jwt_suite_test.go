package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGqlgenJwt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GqlgenJwt Suite")
}
