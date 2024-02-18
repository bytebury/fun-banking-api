package handlers

import (
	"funbanking/internal/domain/metrics"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MetricHandler struct {
	metricService metrics.MetricService
}

func NewMetricsHandler() MetricHandler {
	return MetricHandler{
		metricService: metrics.NewMetricService(),
	}
}

func (h MetricHandler) GetApplicationInfo(c *gin.Context) {
	appInfo, _ := h.metricService.GetApplicationInfo()
	c.JSON(http.StatusOK, appInfo)
}
