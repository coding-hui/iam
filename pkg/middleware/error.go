// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"net"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/pkg/api"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
)

// GinRecovery custom error output.
func GinRecovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			customRecovery(err, c)
		}
	}()
	c.Next()
}

func customRecovery(err any, c *gin.Context) {
	if c.IsAborted() {
		c.Status(200)
	}
	// Check for a broken connection, as it is not really a
	// condition that warrants a panic stack trace.
	var brokenPipe bool
	var ne *net.OpError
	var se *os.SyscallError
	if ok := errors.As(err.(error), &ne); ok {
		if ok := errors.As(err.(error), &se); ok {
			if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
				strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
				brokenPipe = true
			}
		}
	}

	httpRequest, _ := httputil.DumpRequest(c.Request, false)
	log.Errorf("[Recovery] %s panic recovered. err: %s. request: %s.",
		timeFormat(time.Now()), err, httpRequest)
	if brokenPipe {
		// If the connection is dead, we can't write a status to it.
		_ = c.Error(err.(error))
		c.Abort()
	}

	api.FailWithErrCode(err.(error), c)
}

func timeFormat(t time.Time) string {
	timeString := t.Format("2006/01/02 - 15:04:05")
	return timeString
}
