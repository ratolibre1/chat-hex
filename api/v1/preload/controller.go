package preload

import (
	"chat-hex/api/common"
	"chat-hex/business/preload"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	service preload.Service
}

func NewController(service preload.Service) *Controller {
	return &Controller{
		service,
	}
}

func (controller *Controller) PopulateMongoDB(c echo.Context) error {
	err := controller.service.PopulateMongoDB()
	if err != nil {
		return c.JSON(common.NewErrorBusinessResponse(err))
	}

	return c.JSON(common.NewSuccessResponseWithoutData())
}
