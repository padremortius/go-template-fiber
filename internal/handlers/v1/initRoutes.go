package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/padremortius/go-template-fiber/internal/config"
	"github.com/padremortius/go-template-fiber/internal/storage"
	"github.com/padremortius/go-template-fiber/pkgs/svclogger"
)

type (
	v1Routes struct {
		cfg   config.Config
		log   svclogger.Log
		store storage.StorageInterface
	}
)

func InitAppRouter(app fiber.Router, aCfg config.Config, aLog svclogger.Log, aStore storage.StorageInterface) {
	v1 := v1Routes{cfg: aCfg, log: aLog, store: aStore}
	app.Add(http.MethodGet, "/v1/test", v1.getTest)
}
