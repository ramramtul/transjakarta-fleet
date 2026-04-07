package handler

import (
	"net/http"
	"strconv"

	"transjakarta-fleet/internal/service"

	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	service *service.LocationService
}

func NewVehicleHandler(service *service.LocationService) *VehicleHandler {
	return &VehicleHandler{service: service}
}

func (h *VehicleHandler) GetLatestLocation(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")

	result, err := h.service.GetLatest(c.Request.Context(), vehicleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vehicle location not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *VehicleHandler) GetHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	startStr := c.Query("start")
	endStr := c.Query("end")

	start, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start"})
		return
	}

	end, err := strconv.ParseInt(endStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end"})
		return
	}

	result, err := h.service.GetHistory(c.Request.Context(), vehicleID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get history"})
		return
	}

	c.JSON(http.StatusOK, result)
}
