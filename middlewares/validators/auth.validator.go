package validators

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"ledger-service/models"
	"net/http"
)

func AuthValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var authRequest models.AuthRequest
		_ = c.ShouldBindBodyWith(&authRequest, binding.JSON)

		if err := authRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}

func RefreshValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var refreshRequest models.RefreshRequest
		_ = c.ShouldBindBodyWith(&refreshRequest, binding.JSON)

		if err := refreshRequest.Validate(); err != nil {
			models.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}
