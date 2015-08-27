package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"

	"github.com/fgrosse/goldi/goldigen"
)

var _ = Describe("GoPathChecker", func() {
	var verbose = false

	Context("without GOPATH environment variable", func() {
		It("should return an empty string if $GOPATH is not set", func() {
			checker := main.NewGoPathChecker(verbose)
			Expect(checker.PackageName("")).To(BeEmpty())
			Expect(checker.PackageName("/home/fgrosse/go/src/foo/bar/baz/goldi.go")).To(BeEmpty())
		})
	})

	Context("with GOPATH environment variable", func() {
		var (
			goPath         = "/home/fgrosse/go"
			originalGoPath = os.Getenv("GOPATH")
		)

		BeforeEach(func() {
			os.Setenv("GOPATH", goPath)
		})

		AfterEach(func() {
			os.Clearenv()
			os.Setenv("GOPATH", originalGoPath)
		})

		It("should return an empty string if the given output path is empty", func() {
			checker := main.NewGoPathChecker(verbose)
			Expect(checker.PackageName("")).To(BeEmpty())
		})

		It("should return an empty string if the given output path is not inside the $GOPATH", func() {
			checker := main.NewGoPathChecker(verbose)
			Expect(checker.PackageName("/usr/lib/go/src/foo/bar/baz/goldi.go")).To(BeEmpty())
		})

		It("should return the correct package name if the given output path is inside $GOPATH", func() {
			checker := main.NewGoPathChecker(verbose)
			Expect(checker.PackageName("/home/fgrosse/go/src/foo/bar/baz/goldi.go")).To(Equal("foo/bar/baz"))
		})

		Context("with relative output dirs inside the $GOPATH", func() {
			BeforeEach(func() {
				os.Setenv("GOPATH", originalGoPath)
			})

			It("should return the correct package name", func() {
				checker := main.NewGoPathChecker(verbose)
				Expect(checker.PackageName("some_file.go")).To(Equal("github.com/fgrosse/goldi/goldigen"))
			})

			It("should return the correct package name when navigating up the file tree", func() {
				checker := main.NewGoPathChecker(verbose)
				Expect(checker.PackageName("../some_file.go")).To(Equal("github.com/fgrosse/goldi"))
			})

			It("should return the correct package name when navigating up and down the file tree", func() {
				checker := main.NewGoPathChecker(verbose)
				Expect(checker.PackageName("../goldigen/some_file.go")).To(Equal("github.com/fgrosse/goldi/goldigen"))
			})

			It("should return the correct package name when navigating into different directories of the file tree", func() {
				checker := main.NewGoPathChecker(verbose)
				Expect(checker.PackageName("../some_dir/some_file.go")).To(Equal("github.com/fgrosse/goldi/some_dir"))
			})
		})
	})
})
