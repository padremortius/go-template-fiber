package redoc

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

func New(doc Redoc) fiber.Handler {
	return adaptor.HTTPHandlerFunc(doc.Handler())
}
