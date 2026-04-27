package handler

import (
	"server-watcher-app/src/generic"
	modelsDTOs "server-watcher-app/src/models/dtos"
	servicesAbstarct "server-watcher-app/src/services/abstract"

	"github.com/gofiber/fiber/v2"
)

type MetricViewerHandler struct {
	metricService servicesAbstarct.MetricViewerService
}

func NewMetricViewerHandler(metricService servicesAbstarct.MetricViewerService) *MetricViewerHandler {
	return &MetricViewerHandler{
		metricService: metricService,
	}
}

func (h *MetricViewerHandler) GetMetric(c *fiber.Ctx) error {

	req, err := generic.ParseBody[modelsDTOs.ContainerLogRequestDTO](c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.metricService.GetLogsFromDB(c.Context(), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}
