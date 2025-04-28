package v1

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type appRoutes struct {

}

func InitAppRouter(app *fiber.App, appName string) {
	v1App := app.Group(fmt.Sprint("/", appName, "/v1"))

	v1App.Get("test", getTest)
}
