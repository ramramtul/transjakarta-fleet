package router

import (
	"transjakarta-fleet/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(vehicleHandler *handler.VehicleHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/vehicles/:vehicle_id/location", vehicleHandler.GetLatestLocation)
	r.GET("/vehicles/:vehicle_id/history", vehicleHandler.GetHistory)

	return r
}
