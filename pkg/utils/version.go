package utils

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
