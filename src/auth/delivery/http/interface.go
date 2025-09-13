package http

import "github.com/gofiber/fiber/v2"

type (
	Handler interface {
		Login(c *fiber.Ctx) error
		Register(c *fiber.Ctx) error
	}
)
