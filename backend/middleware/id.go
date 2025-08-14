package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func requireUintParam(param, ctxKey, label string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := strings.TrimSpace(c.Param(param))
		u64, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || u64 == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + label})
			c.Abort()
			return
		}
		u := uint(u64)
		if uint64(u) != u64 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + label})
			c.Abort()
			return
		}
		c.Set(ctxKey, u)
		c.Next()
	}
}

const (
	CtxKeyUserID           = "userID"
	CtxKeyJobApplicationID = "jobApplicationID"
)

func RequireUserID() gin.HandlerFunc           { return requireUintParam("id", CtxKeyUserID, "user id") }
func RequireJobApplicationID() gin.HandlerFunc { return requireUintParam("id", CtxKeyJobApplicationID, "job application id") }
