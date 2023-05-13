// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"math/rand"
	"time"

	"github.com/coding-hui/iam/cmd/iam-authzserver/app"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	app.NewAuthzServerAPP("iam-authzserver").Run()
}
