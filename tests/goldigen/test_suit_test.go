package goldigen

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestGoldiGen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoldiGen Test Suite")
}
