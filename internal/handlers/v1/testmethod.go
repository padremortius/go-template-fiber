package v1

import (
	"net/http"

	fiber "github.com/gofiber/fiber/v2"
)

// @Summary Test method
// @Description Test method
// @Produce json
// @Success 200 {object} JSONResult
// @Router /go-template-fiber/v1/test [get]
// @Tags v1
func (v1 *v1Routes) getTest(c *fiber.Ctx) error {
	v1.log.Logger.Info().Msg("Start getTest")
	return c.JSON(&JSONResult{Code: http.StatusOK, Message: "Test complete!"})
}
