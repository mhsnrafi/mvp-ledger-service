package routes

import (
	"github.com/gin-gonic/gin"
	"ledger-service/controllers"
	"ledger-service/middlewares/validators"
)

func AuthRoute(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST(
			"/generate_access_token",
			validators.AuthValidator(),
			controllers.GenerateAccessToken,
		)

		auth.POST(
			"/refresh",
			validators.RefreshValidator(),
			controllers.Refresh,
		)
	}
}
