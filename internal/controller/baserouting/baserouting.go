package baserouting

import (
	"go-template-fiber/internal/config"

	fiber "github.com/gofiber/fiber/v2"
)

type (
	Health struct {
		Status string
	}
)

func getHealth(c *fiber.Ctx) error {
	return c.JSON(Health{Status: "up"})
}

func getInfo(c *fiber.Ctx) error {
	return c.JSON(config.Cfg.Version)
}

func getEnv(c *fiber.Ctx) error {
	return c.JSON(config.Cfg)
}

func InitBaseRouter(app *fiber.App) {
	// K8s probe
	app.Get("/health", getHealth)

	// info about service
	app.Get("/info", getInfo)

	// env
	app.Get("/env", getEnv)
}
