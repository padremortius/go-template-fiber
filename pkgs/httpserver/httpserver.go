package httpserver

import (
	"context"
	"time"

	"github.com/padremortius/go-template-fiber/pkgs/svclogger"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type HTTP struct {
	Port            string `yaml:"port" json:"port" default:"8080"`
	SwaggerDisabled bool   `yaml:"swaggerDisabled" json:"swaggerDisabled" validate:"required"`
	Timeouts        struct {
		Read     time.Duration `yaml:"read" json:"read" default:"30s"`
		Write    time.Duration `yaml:"write" json:"write" default:"30s"`
		Idle     time.Duration `yaml:"idle" json:"idle" default:"30s"`
		Shutdown time.Duration `yaml:"shutdown" json:"shutdown" default:"30s"`
	} `yaml:"timeouts" json:"timeouts"`
}

type AppServer struct {
	ctx     context.Context
	server  *fasthttp.Server
	Handler *fiber.App
	notify  chan error
}

// New -.
func New(c context.Context, log *svclogger.Log, opts *HTTP) *AppServer {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Logger settings
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger:   log.Logger,
		Fields:   []string{"latency", "status", "method", "url", "ua", "ip", "bytesSent"},
		SkipURIs: mySkipper(),
		Messages: []string{"-"},
	}))

	s := &AppServer{
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

	return s
}

func (s *AppServer) Start(aPort string) {
	go func() {
		s.notify <- s.Handler.Listen(":" + aPort)
		close(s.notify)
	}()
}

// Notify -.
func (s *AppServer) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *AppServer) Shutdown(shutdownTimeout time.Duration) error {
	ctx, cancel := context.WithTimeout(s.ctx, shutdownTimeout)
	defer cancel()

	return s.server.ShutdownWithContext(ctx)
}
