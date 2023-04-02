package utils

import (
	"strings"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test version utils", func() {
	It("Test New version function", func() {
		s := GenerateVersion("")
		Expect(s).ShouldNot(BeNil())

		s2 := GenerateVersion("pre")
		Expect(cmp.Diff(strings.HasPrefix(s2, "pre-"), true)).ShouldNot(BeNil())
	})
})
