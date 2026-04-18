// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/coding-hui/iam/internal/apiserver"
	"github.com/coding-hui/iam/pkg/app"
)

func main() {
	app.NewApp(
		"apiserver",
		"apiserver",
		app.WithDescription("WeCoding IAM API Server"),
		app.WithOptions(apiserver.NewOptions()),
		app.WithRunFunc(apiserver.Run),
	).Run()
}
