package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("sanitizer", func() {
	var s *sanitizer
	BeforeEach(func() {
		s = newSanitizer()
	})

	Describe("escaping @", func() {
		It("should escape @ if not inside q quoted string", func() {
			s.Write([]byte(`john.doe@example.com`))
			Expect(string(s.Bytes())).To(Equal(`john.doe\@example.com`))
		})

		It("should not escape @ inside single quoted string ", func() {
			s.Write([]byte(`'john.doe@example.com'`))
			Expect(string(s.Bytes())).To(Equal(`'john.doe@example.com'`))
		})

		It("should not escape @ inside double quoted string ", func() {
			s.Write([]byte(`"john.doe@example.com"`))
			Expect(string(s.Bytes())).To(Equal(`"john.doe@example.com"`))
		})

		It("should not escape @ inside quoted string with escaped single quotes", func() {
			s.Write([]byte(`'john''.doe@example.com'`))
			Expect(string(s.Bytes())).To(Equal(`'john''.doe@example.com'`))
		})

		It("should not escape @ inside quoted string with escaped double quotes", func() {
			s.Write([]byte(`"john"".doe@example.com"`))
			Expect(string(s.Bytes())).To(Equal(`"john"".doe@example.com"`))
		})

		It("should not escape @ inside single quoted string with newlines", func() {
			s.Write([]byte("'User:\njohn.doe@example.com'"))
			Expect(string(s.Bytes())).To(Equal("'User:\njohn.doe@example.com'"))
		})

		It("should not escape @ inside double quoted string with newlines", func() {
			s.Write([]byte("\"User:\njohn.doe@example.com\""))
			Expect(string(s.Bytes())).To(Equal("\"User:\njohn.doe@example.com\""))
		})
	})
})
