// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utils

import (
	"regexp"
	"strings"
)

func CamelToUnderscore(s string) string {
	// 使用正则表达式将驼峰字符串转换为下划线分隔的形式
	reg := regexp.MustCompile("([a-z0-9])([A-Z])")
	underscoreStr := reg.ReplaceAllString(s, "${1}_${2}")

	// 将转换后的字符串全部转为小写
	underscoreStr = strings.ToLower(underscoreStr)

	return underscoreStr
}
