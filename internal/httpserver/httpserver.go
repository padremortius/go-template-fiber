package httpserver

import (
	"context"
	"fmt"
	"go-template-fiber/internal/svclogger"
	"time"

	fiberPrometheus "github.com/ansrivas/fiberprometheus/v2"
	gojson "github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mvrilo/go-redoc"
	fiberredoc "github.com/mvrilo/go-redoc/fiber"
	"github.com/valyala/fasthttp"
)

type HTTP struct {
	Cors struct {
		Headers string `yaml:"headers" json:"headers" validate:"required"`
		Methods string `yaml:"methods" json:"methods" validate:"required"`
		Origins string `yaml:"origins" json:"origins" validate:"required"`
	} `yaml:"cors" json:"cors"`
	Port            string `yaml:"port" json:"port"`
	SwaggerDisabled bool   `yaml:"swaggerDisabled" json:"swaggerDisabled" validate:"required"`
	Timeouts        struct {
		Read     time.Duration `yaml:"read" json:"read"`
		Write    time.Duration `yaml:"write" json:"write"`
		Idle     time.Duration `yaml:"idle" json:"idle"`
		Shutdown time.Duration `yaml:"shutdown" json:"shutdown"`
	} `yaml:"timeouts" json:"timeouts"`
	// Token struct {
	// 	PubKeyURL string `yaml:"pubKeyURL" json:"pubKeyURL" validate:"required"`
	// 	PublicKey string `yaml:"publicKey" json:"publicKey" validate:"required"`
	// } `yaml:"token" json:"token"`
}

type Server struct {
	ctx     context.Context
	server  *fasthttp.Server
	Handler *fiber.App
	notify  chan error
}

// New -.
func New(c context.Context, log *svclogger.Log, opts *HTTP) *Server {
	app := fiber.New(fiber.Config{
		JSONEncoder: gojson.Marshal,
		JSONDecoder: gojson.Unmarshal,
	})

	app.Use(recover.New())

	// CORS settings
	app.Use(cors.New(cors.Config{
		AllowHeaders: opts.Cors.Headers,
		AllowMethods: opts.Cors.Methods,
		AllowOrigins: opts.Cors.Origins,
	}))

	// Logger settings
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger:   log.Logger,
		Fields:   []string{"latency", "status", "method", "url", "ua", "ip", "bytesSent"},
		SkipURIs: []string{"/favicon.ico", "/health"},
		Messages: []string{"-"},
	}))

	// metrics settings
	prometheus := fiberPrometheus.New("fiber")
	prometheus.RegisterAt(app, "/prometheus")
	app.Use(prometheus.Middleware)

	// redoc
	if opts.SwaggerDisabled {
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

	s := &Server{
		server: &fasthttp.Server{
			Handler:      app.Handler(),
			IdleTimeout:  opts.Timeouts.Idle,
			ReadTimeout:  opts.Timeouts.Read,
			WriteTimeout: opts.Timeouts.Write,
		},
		notify:  make(chan error, 1),
		Handler: app,
		ctx:     c,
	}

	s.start(fmt.Sprint(":", opts.Port))

	return s
}

func (s *Server) start(aPort string) {
	go func() {
		s.notify <- s.Handler.Listen(aPort)
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown(shutdownTimeout time.Duration) error {
	ctx, cancel := context.WithTimeout(s.ctx, shutdownTimeout)
	defer cancel()

	return s.server.ShutdownWithContext(ctx)
}
