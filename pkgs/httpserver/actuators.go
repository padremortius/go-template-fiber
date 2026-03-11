package httpserver

import (
	fiberprometheus "github.com/ansrivas/fiberprometheus/v2"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/padremortius/go-template-fiber/pkgs/redoc"
	"github.com/padremortius/go-template-fiber/pkgs/svclogger"
)

type (
	Health struct {
		Status string
	}

	BaseRoutes struct {
		cfg     any
		version any
		log     svclogger.Log
	}
)

func (b *BaseRoutes) getHealth(c *fiber.Ctx) error {
	return c.JSON(Health{Status: "up"})
}

func (b *BaseRoutes) getInfo(c *fiber.Ctx) error {
	return c.JSON(&b.version)
}

func (b *BaseRoutes) getEnv(c *fiber.Ctx) error {
	return c.JSON(&b.cfg)
}

func InitBaseRouter(app *fiber.App, aSrvName string, aCfg any, aVersion any, aLog svclogger.Log) {
	bRoutes := BaseRoutes{cfg: aCfg, version: aVersion, log: aLog}

	// K8s probe
	app.Get("/health", bRoutes.getHealth)

	// info about service
	app.Get("/info", bRoutes.getInfo)

	// env
	app.Get("/env", bRoutes.getEnv)

	// metrics
	prometheus := fiberprometheus.NewWithDefaultRegistry(aSrvName)
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// redoc
	//if bRoutes.cfg.HTTP.SwaggerDisabled {
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

	app.Use(redoc.New(doc))
	//}
}
