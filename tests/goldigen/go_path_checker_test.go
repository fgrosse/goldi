package goldigen

import (
	"github.com/fgrosse/goldi/generator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("GoPathChecker", func() {
	Context("without GOPATH environment variable", func() {
		It("should return an empty string if $GOPATH is not set", func() {
			checker := generator.NewGoPathChecker()
			Expect(checker.PackageName("")).To(BeEmpty())
			Expect(checker.PackageName("/home/fgrosse/go/src/foo/bar/baz/goldi.go")).To(BeEmpty())
		})
	})

	Context("without GOPATH environment variable", func() {
		var (
			goPath         = "/home/fgrosse/go"
			originalGoPath = os.Getenv("GOPATH")
		)

		BeforeEach(func() {
			os.Setenv("GOPATH", goPath)
		})

		AfterEach(func() {
			os.Clearenv()
		})

		It("should return an empty string if the given output path is empty", func() {
			checker := generator.NewGoPathChecker()
			Expect(checker.PackageName("")).To(BeEmpty())
		})

		It("should return an empty string if the given output path is not inside the $GOPATH", func() {
			checker := generator.NewGoPathChecker()
			Expect(checker.PackageName("/usr/lib/go/src/foo/bar/baz/goldi.go")).To(BeEmpty())
		})

		It("should return the correct package name if the given output path is inside $GOPATH", func() {
			checker := generator.NewGoPathChecker()
			Expect(checker.PackageName("/home/fgrosse/go/src/foo/bar/baz/goldi.go")).To(Equal("foo/bar/baz"))
		})

		Context("with relative output dirs inside the $GOPATH", func() {
			BeforeEach(func() {
				os.Setenv("GOPATH", originalGoPath)
			})

			It("should return the correct package name", func() {
				checker := generator.NewGoPathChecker()
				Expect(checker.PackageName("goldi.go")).To(Equal("github.com/fgrosse/goldi/tests/goldigen"))
			})

			It("should return the correct package name when navigating up the file tree", func() {
				checker := generator.NewGoPathChecker()
				Expect(checker.PackageName("../goldi.go")).To(Equal("github.com/fgrosse/goldi/tests"))
			})

			It("should return the correct package name when navigating up and down the file tree", func() {
				checker := generator.NewGoPathChecker()
				Expect(checker.PackageName("../goldigen/goldi.go")).To(Equal("github.com/fgrosse/goldi/tests/goldigen"))
			})

			It("should return the correct package name when navigating into different directories of the file tree", func() {
				checker := generator.NewGoPathChecker()
				Expect(checker.PackageName("../testAPI/goldi.go")).To(Equal("github.com/fgrosse/goldi/tests/testAPI"))
			})

		})
	})
})
