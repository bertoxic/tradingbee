package routes

import (
	"net/http"

	handler "github.com/bertoxic/tradingbee/internal/handlers"
	"github.com/bertoxic/tradingbee/internal/models"
	"github.com/bertoxic/tradingbee/internal/transport/httputil"

	"github.com/gin-gonic/gin"
)

func Router() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/authenticate", func(ctx *gin.Context) {
		httputil.WriteJson(ctx, true, 200, &models.JsonResponse{
			Success: true,
			Message: "duely authenticated",
			Data:    "no data currently",
		})
	})
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "authenticated no wow",
		})
	})
	
	router.POST("/signup", handler.SignUp)
	router.POST("/login", handler.Login)
	router.POST("/otp", handler.GenerateOTPResponse)
	

	return router
}
