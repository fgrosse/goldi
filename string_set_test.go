package goldi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StringSet", func() {
	It("should have a simple method to check if a key is set", func() {
		s := StringSet{}
		key := "test"
		Expect(s.Contains(key)).To(BeFalse())

		s.Set(key)
		Expect(s.Contains(key)).To(BeTrue())
	})
})
