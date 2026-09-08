package httpapi

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type tokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// Client credentials contra el entorno; el reto no pide un IdP.
func (h *Handler) IssueToken(c *fiber.Ctx) error {
	var body tokenRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "cuerpo inválido: se espera client_id y client_secret")
	}
	if body.ClientID != h.cfg.ClientID || body.ClientSecret != h.cfg.ClientSecret {
		return fiber.NewError(fiber.StatusUnauthorized, "credenciales inválidas")
	}

	token, expiresAt, err := h.signToken(body.ClientID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "no se pudo emitir el token")
	}

	return c.JSON(fiber.Map{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   int(time.Until(expiresAt).Seconds()),
	})
}

func (h *Handler) signToken(subject string) (string, time.Time, error) {
	expiresAt := time.Now().Add(h.cfg.JWTTTL)
	claims := jwt.MapClaims{
		"sub": subject,
		"iss": "qr-api-go",
		"iat": time.Now().Unix(),
		"exp": expiresAt.Unix(),
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.cfg.JWTSecret))
	return signed, expiresAt, err
}

func RequireJWT(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "falta el header Authorization: Bearer <token>")
		}
		raw := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))

		token, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "token inválido o expirado")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("subject", claims["sub"])
		}
		return c.Next()
	}
}
