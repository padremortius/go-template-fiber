package actuators

import (
	fiberprometheus "github.com/ansrivas/fiberprometheus/v2"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/mvrilo/go-redoc"
	fiberredoc "github.com/mvrilo/go-redoc/fiber"
	"github.com/padremortius/go-template-fiber/internal/config"
	"github.com/padremortius/go-template-fiber/internal/svclogger"
)

type (
	Health struct {
		Status string
	}

	BaseRoutes struct {
		cfg config.Config
		log svclogger.Log
	}
)

func (b *BaseRoutes) getHealth(c *fiber.Ctx) error {
	return c.JSON(Health{Status: "up"})
}

func (b *BaseRoutes) getInfo(c *fiber.Ctx) error {
	return c.JSON(b.cfg.Version)
}

func (b *BaseRoutes) getEnv(c *fiber.Ctx) error {
	return c.JSON(&b.cfg)
}

func InitBaseRouter(app *fiber.App, aCfg config.Config, aLog svclogger.Log) {
	bRoutes := BaseRoutes{cfg: aCfg, log: aLog}

	// K8s probe
	app.Get("/health", bRoutes.getHealth)

	// info about service
	app.Get("/info", bRoutes.getInfo)

	// env
	app.Get("/env", bRoutes.getEnv)

	// metrics
	prometheus := fiberprometheus.NewWithDefaultRegistry(bRoutes.cfg.BaseApp.Name)
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// redoc
	if bRoutes.cfg.HTTP.SwaggerDisabled {
		doc := redoc.Redoc{
			SpecFile: "spec/docs.json",
			SpecPath: "/docs.json",
			DocsPath: "/docs",
			Options: map[string]any{
				"disableSearch": true,
				"theme": map[string]any{
					"colors":     map[string]any{"primary": map[string]any{"main": "#297b21"}},
					"typography": map[string]any{"headings": map[string]any{"fontWeight": "600"}},
					"sidebar":    map[string]any{"backgroundColor": "lightblue"},
				},
				"decorator": map[string]any{},
			},
		}

		app.Use(fiberredoc.New(doc))
	}
}
