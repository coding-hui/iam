// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package version

import (
	"fmt"
	"strings"
	"time"
)

// GenerateVersion Generate version numbers by time
func GenerateVersion(pre string) string {
	timeStr := time.Now().Format("20060102150405.000")
	timeStr = strings.Replace(timeStr, ".", "", 1)
	if pre != "" {
		return fmt.Sprintf("%s-%s", pre, timeStr)
	}
	return timeStr
}
