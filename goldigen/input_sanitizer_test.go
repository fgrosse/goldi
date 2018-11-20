package main

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type BufferMock struct {
}

func (*BufferMock) WriteString(string) (int, error) {
	return 0, errors.New("err")
}

func (*BufferMock) WriteByte(c byte) error {
	return errors.New("err")
}

func (*BufferMock) Bytes() []byte {
	return nil
}

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

		It("should return error when cannot write string to buffer", func() {
			s = &sanitizer{
				buf:      &BufferMock{},
				inQuotes: false,
			}
			_, err := s.Write([]byte("@"))
			Expect(err).To(MatchError("err"))
		})

		It("should return error when cannot write byte to buffer", func() {
			s = &sanitizer{
				buf:      &BufferMock{},
				inQuotes: false,
			}
			_, err := s.Write([]byte("doe@"))
			Expect(err).To(MatchError("err"))
		})
	})
})
