package main_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fgrosse/goldi/goldigen"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoPathChecker", func() {
	newPathChecker := func() *main.GoPathChecker {
		checker := main.NewGoPathChecker(true)
		checker.Logger = GinkgoWriter
		return checker
	}

	originalGoPath := os.Getenv("GOPATH")
	AfterEach(func() {
		os.Setenv("GOPATH", originalGoPath)
	})

	Context("without GOPATH environment variable", func() {
		BeforeEach(func() {
			_ = os.Unsetenv("GOPATH")
		})

		It("should return an empty string if $GOPATH is not set", func() {
			checker := newPathChecker()
			Expect(checker.PackageName("")).To(BeEmpty())
			Expect(checker.PackageName("/home/fgrosse/go/src/foo/bar/baz/goldi.go")).To(BeEmpty())
		})
	})

	Context("with GOPATH environment variable", func() {
		goPath := "/home/fgrosse/go"

		BeforeEach(func() {
			os.Setenv("GOPATH", goPath)
		})

		It("should return an empty string if the given output path is empty", func() {
			checker := newPathChecker()
			Expect(checker.PackageName("")).To(BeEmpty())
		})

		It("should return an empty string if the given output path is not inside the $GOPATH", func() {
			checker := newPathChecker()
			Expect(checker.PackageName("/usr/lib/go/src/foo/bar/baz/goldi.go")).To(BeEmpty())
		})

		It("should return the correct package name if the given output path is inside $GOPATH", func() {
			checker := newPathChecker()
			Expect(checker.PackageName("/home/fgrosse/go/src/foo/bar/baz/goldi.go")).To(Equal("foo/bar/baz"))
		})

		It("should log message in verbose mode", func() {
			logger := new(bytes.Buffer)
			checker := main.NewGoPathChecker(true)
			checker.Logger = logger

			checker.PackageName("/home/fgrosse/go/src/foo/bar/baz/goldi.go")
			Expect(logger.String()).NotTo(BeEmpty())
		})

		Specify("behaviour with relative output dirs inside the $GOPATH", func() {
			os.Setenv("GOPATH", originalGoPath)

			// If this repository is not checked out inside the GOPATH we are
			// skipping these tests. On Travis we should always execute them
			// however.
			if !isInGoPath(pwd()) {
				Skip("Skipping test because repository is not checked out inside $GOPATH")
			}

			checker := newPathChecker()
			log := func(s string) {
				fmt.Fprintln(GinkgoWriter, s)
			}

			log("It should return the correct package name")
			Expect(checker.PackageName("some_file.go")).To(Equal("github.com/fgrosse/goldi/goldigen"))

			log("It should return the correct package name when navigating up the file tree")
			Expect(checker.PackageName("../some_file.go")).To(Equal("github.com/fgrosse/goldi"))

			log("It should return the correct package name when navigating up and down the file tree")
			Expect(checker.PackageName("../goldigen/some_file.go")).To(Equal("github.com/fgrosse/goldi/goldigen"))

			log("It should return the correct package name when navigating into different directories of the file tree")
			Expect(checker.PackageName("../some_dir/some_file.go")).To(Equal("github.com/fgrosse/goldi/some_dir"))
		})
	})
})

func pwd() string {
	path, _ := filepath.Abs("main.go")
	return filepath.Dir(path)
}

func isInGoPath(p string) bool {
	for _, goPath := range strings.Split(os.Getenv("GOPATH"), ":") {
		goPath = goPath + "/src/"
		if strings.Contains(p, goPath) {
			return true
		}
	}

	return false
}
