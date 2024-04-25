package common

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"k8s-delegater/pkg/common/response"
)

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(*response.Result); ok {
				c.JSON(e.Code, e)
			} else {
				c.JSON(http.StatusInternalServerError, e)
			}
			c.Abort()
		}
	}()
	c.Next()
}

func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}
