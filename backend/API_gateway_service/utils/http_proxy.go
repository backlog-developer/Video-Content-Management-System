package utils

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// ForwardRequest handles GET/DELETE without body (dynamic method)
func ForwardRequest(c *fiber.Ctx, targetURL string) error {
	req, err := http.NewRequest(c.Method(), targetURL, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Request creation failed")
	}
	copyHeaders(c, req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).SendString("Failed to reach service")
	}
	defer resp.Body.Close()

	return relayResponse(c, resp)
}

// ForwardRequestWithBody handles POST/PUT with JSON body
func ForwardRequestWithBody(c *fiber.Ctx, targetURL string) error {
	reqBody := bytes.NewBuffer(c.Body())

	req, err := http.NewRequest(c.Method(), targetURL, reqBody)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Request creation failed")
	}
	copyHeaders(c, req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).SendString("Failed to reach service")
	}
	defer resp.Body.Close()

	return relayResponse(c, resp)
}

// Copies all request headers (e.g., Authorization) to the outgoing request
func copyHeaders(c *fiber.Ctx, req *http.Request) {
	c.Request().Header.VisitAll(func(k, v []byte) {
		req.Header.Set(string(k), string(v))
	})
}

// Sends the response back to the client
func relayResponse(c *fiber.Ctx, resp *http.Response) error {
	for k, v := range resp.Header {
		if len(v) > 0 {
			c.Set(k, v[0])
		}
	}
	body, _ := io.ReadAll(resp.Body)
	return c.Status(resp.StatusCode).Send(body)
}
