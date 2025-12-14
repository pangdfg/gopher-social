package main

import (
	"io"
	"net/http/httptest"
	"time"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pangdfg/gopher-social/internal/env"
	"github.com/pangdfg/gopher-social/internal/ratelimiter"
)

type application_test struct {
	app *fiber.App
	config
}

type ratelimiterConfig struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}


func NewRateLimiter(cfg ratelimiterConfig) fiber.Handler {
	counter := make(map[string]int)
	resetTime := time.Now().Add(cfg.TimeFrame)

	return func(c *fiber.Ctx) error {
		ip := c.Get("X-Forwarded-For", c.IP())

		if time.Now().After(resetTime) {
			counter = make(map[string]int)
			resetTime = time.Now().Add(cfg.TimeFrame)
		}

		counter[ip]++
		if cfg.Enabled && counter[ip] > cfg.RequestsPerTimeFrame {
			return c.SendStatus(fiber.StatusTooManyRequests)
		}

		return c.Next()
	}
}


func (a *application_test) mount() *fiber.App {
	
	a.app.Get("/v1/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	if a.config.rateLimiter.Enabled {
		a.app.Use(NewRateLimiter(ratelimiterConfig{
			RequestsPerTimeFrame: a.config.rateLimiter.RequestsPerTimeFrame,
			TimeFrame:            a.config.rateLimiter.TimeFrame,
			Enabled:              a.config.rateLimiter.Enabled,
		}))
	}

	return a.app
}

func TestApplication(cfg config) *application_test {
	app := fiber.New()
	inst := &application_test{
		app:    app,
		config: cfg,
	}
	inst.mount()
	return inst
}

var _ = Describe("API Test", func() {

	var (
		cfg     config
		appInst *application_test
		mockIP  string
	)

	BeforeEach(func() {

		cfg = config{
			addr:        env.GetString("ADDR", ":8080"),
			apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
			env: env.GetString("ENV", "development"),
			rateLimiter: ratelimiter.Config{
				RequestsPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20),
				TimeFrame:            time.Second * 5,
				Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
			},
		}

		appInst = TestApplication(cfg)
		mockIP = "192.168.1.1"
	})

	It("should allow first N requests then return 429", func() {
		limit := cfg.rateLimiter.RequestsPerTimeFrame
		margin := 2

		for i := 0; i < limit+margin; i++ {
			req := httptest.NewRequest("GET", "/v1/health", nil)
			req.Header.Set("X-Forwarded-For", mockIP)

			resp, err := appInst.app.Test(req, -1)
			Expect(err).To(BeNil())

			body, _ := io.ReadAll(resp.Body)
			_ = body 

			if i < limit {
				Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
			} else {
				Expect(resp.StatusCode).To(Equal(fiber.StatusTooManyRequests))
			}
		}
	})
})
