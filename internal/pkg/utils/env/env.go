// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package env

type (
	Mode string
)

const (
	ModeDev  Mode = "dev"  //开发模式
	ModeTest Mode = "test" //测试模式
	ModeProd Mode = "prod" //生产模式
)

func (e Mode) String() string {
	return string(e)
}
