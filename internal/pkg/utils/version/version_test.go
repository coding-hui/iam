// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package version

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
