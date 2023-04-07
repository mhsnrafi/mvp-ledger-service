package routes

import (
	"github.com/gin-gonic/gin"
	"ledger-service/controllers"
)

func Legder(router *gin.RouterGroup) {
	auth := router.Group("/")
	{
		auth.POST(
			"users/:uid/add",
			//middlewares.JWTMiddleware(),
			controllers.AddFunds,
		)
		auth.GET(
			"users/:uid/balance",
			//middlewares.JWTMiddleware(),
			controllers.GetBalance,
		)
		auth.GET(
			"users/:uid/history",
			//middlewares.JWTMiddleware(),
			controllers.GetTransactionHistory,
		)
	}
}
