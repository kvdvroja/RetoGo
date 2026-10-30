package handler

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"matrix-qr-api/internal/auth"
	"matrix-qr-api/internal/matrix"
	"matrix-qr-api/internal/stats"
)

const maxDimension = 100

type Options struct {
	JWTSecret     []byte
	TokenTTL      time.Duration
	Username      string
	Password      string
	AllowedOrigin string
	Stats         stats.Calculator
	Logging       bool
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
}

type qrRequest struct {
	Matrix matrix.Matrix `json:"matrix"`
}

type qrResponse struct {
	Q          matrix.Matrix     `json:"q"`
	R          matrix.Matrix     `json:"r"`
	Statistics *stats.Statistics `json:"statistics"`
}

func New(opts Options) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "matrix-qr-api",
		BodyLimit:    1 << 20,
		ErrorHandler: errorHandler,
	})

	app.Use(recover.New())
	if opts.Logging {
		app.Use(logger.New())
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: opts.AllowedOrigin,
		AllowHeaders: "Content-Type, Authorization",
		AllowMethods: "GET, POST, OPTIONS",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "matrix-qr-api",
			"endpoints": []string{
				"GET /health",
				"POST /api/auth/login",
				"POST /api/qr",
			},
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api")
	api.Post("/auth/login", login(opts))
	api.Post("/qr", auth.Middleware(opts.JWTSecret), factorize(opts.Stats))

	return app
}

func login(opts Options) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req loginRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		userOK := subtle.ConstantTimeCompare([]byte(req.Username), []byte(opts.Username)) == 1
		passOK := subtle.ConstantTimeCompare([]byte(req.Password), []byte(opts.Password)) == 1
		if !userOK || !passOK {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}

		token, err := auth.Issue(opts.JWTSecret, req.Username, opts.TokenTTL)
		if err != nil {
			return err
		}
		return c.JSON(loginResponse{Token: token, ExpiresIn: int(opts.TokenTTL.Seconds())})
	}
}

func factorize(calc stats.Calculator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req qrRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "body must be JSON with a numeric \"matrix\" field")
		}
		if len(req.Matrix) > maxDimension || (len(req.Matrix) > 0 && len(req.Matrix[0]) > maxDimension) {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("matrix dimensions must not exceed %dx%d", maxDimension, maxDimension))
		}

		q, r, err := matrix.QR(req.Matrix)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		result, err := calc.Calculate(c.UserContext(), auth.Token(c), []stats.NamedMatrix{
			{Name: "Q", Data: q},
			{Name: "R", Data: r},
		})
		if err != nil {
			return fiber.NewError(fiber.StatusBadGateway, err.Error())
		}

		return c.JSON(qrResponse{Q: q, R: r, Statistics: result})
	}
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "internal server error"

	var fe *fiber.Error
	if errors.As(err, &fe) {
		code = fe.Code
		message = fe.Message
	}
	return c.Status(code).JSON(fiber.Map{"error": message})
}
