package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RequiredField struct {
	Value any
	Name  string
}

func ValidateRequiredFields(c *gin.Context, fields []RequiredField) bool {
	for _, field := range fields {
		if isEmpty(field.Value) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("A %s is required.", field.Name),
			})
			return false
		}
	}
	return true
}

func isEmpty(value any) bool {
	switch v := value.(type) {
	case string:
		return v == ""
	case *string:
		return v == nil || *v == ""
	default:
		return value == nil
	}
}
