package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestGoldi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Goldi Test Suite")
}
