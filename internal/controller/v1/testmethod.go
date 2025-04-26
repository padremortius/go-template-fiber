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
func getTest(c *fiber.Ctx) error {
	return c.JSON(&JSONResult{Code: http.StatusOK, Message: "Test complete!"})
}
