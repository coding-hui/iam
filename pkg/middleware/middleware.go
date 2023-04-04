package middleware

import "github.com/gin-gonic/gin"

func InitMiddleware(r *gin.Engine) {
	// Custom Error Wrapper
	r.Use(CustomError)
	// NoCache is a middleware function that appends headers
	r.Use(NoCache)
	// 跨域处理
	r.Use(Options)
	// Secure is a middleware function that appends security
	r.Use(Secure)
	// request log
	r.Use(requestLog())
}
