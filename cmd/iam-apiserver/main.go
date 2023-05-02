// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/coding-hui/iam/cmd/iam-apiserver/app"
)

func main() {
	cmd := app.NewAPIServerCommand()
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
