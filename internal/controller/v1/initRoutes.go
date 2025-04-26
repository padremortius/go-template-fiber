package v1

import (
	"fmt"
	"go-template-fiber/internal/config"

	"github.com/gofiber/fiber/v2"
)

type appRoutes struct {

}

func InitAppRouter(app *fiber.App) {
    v1App := app.Group(fmt.Sprint("/", config.Cfg.BaseApp.Name, "/v1"))

	v1App.Get("test", getTest)
}
