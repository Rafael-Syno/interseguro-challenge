// Package httpapi expone la capa HTTP (Fiber) de la API en Go.
package httpapi

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/rbravo/interseguro-challenge/api-go/internal/config"
	"github.com/rbravo/interseguro-challenge/api-go/internal/matrix"
	"github.com/rbravo/interseguro-challenge/api-go/internal/statsclient"
)

type Handler struct {
	cfg   config.Config
	stats *statsclient.Client
}

func NewHandler(cfg config.Config, stats *statsclient.Client) *Handler {
	return &Handler{cfg: cfg, stats: stats}
}

type qrRequest struct {
	Matrix [][]float64 `json:"matrix"`
	// SkipStatistics permite obtener solo Q y R sin llamar a la API de Node.
	SkipStatistics bool `json:"skipStatistics"`
}

func (h *Handler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok", "service": "qr-api-go"})
}

// QR factoriza la matriz y pide las estadísticas a Node.
func (h *Handler) QR(c *fiber.Ctx) error {
	var body qrRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "JSON inválido: se espera {\"matrix\": [[...]]}")
	}

	a, err := matrix.New(body.Matrix)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	q, r := matrix.QR(a)
	rows, cols := a.Dims()

	qRounded := q.Round(h.cfg.RoundDecimals)
	rRounded := r.Round(h.cfg.RoundDecimals)

	response := fiber.Map{
		"input": fiber.Map{
			"rows":    rows,
			"columns": cols,
			"matrix":  a,
		},
		"factorization": fiber.Map{
			"method": "householder",
			"Q":      qRounded,
			"R":      rRounded,
		},
	}

	if body.SkipStatistics {
		return c.JSON(response)
	}

	// Propaga el JWT a la API de Node (secreto compartido).
	token, _, err := h.signToken("qr-api-go")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "no se pudo firmar el token de servicio")
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), h.cfg.StatsTimeout)
	defer cancel()

	stats, err := h.stats.Statistics(ctx, token, []statsclient.NamedMatrix{
		{Name: "Q", Data: qRounded},
		{Name: "R", Data: rRounded},
	})
	if err != nil {
		// 502: la factorización sí se calculó, falló la dependencia.
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	response["statistics"] = stats
	return c.JSON(response)
}

func Register(app *fiber.App, h *Handler, jwtSecret string) {
	app.Get("/health", h.Health)

	v1 := app.Group("/api/v1")
	v1.Post("/auth/token", h.IssueToken)

	protected := v1.Group("/matrix", RequireJWT(jwtSecret))
	protected.Post("/qr", h.QR)
}
