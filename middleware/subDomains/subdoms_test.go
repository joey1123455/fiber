package subDomains

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// go test -run Benchmark_URLParser -v
func Benchmark_URLParser(b *testing.B) {
	domSegs := "example.com.path.to.resource"
	paramPosition := 3
	var res, hostParam string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, hostParam = urlParser(domSegs, paramPosition)
	}
	utils.AssertEqual(b, "example.com.to.resource", res)
	utils.AssertEqual(b, "path", hostParam)
}

// go test -run Test_URLParser -v
func Test_URLParser(t *testing.T) {
	domSegs := "example.com.path.to.resource"
	res, hostParam := urlParser(domSegs, 3)
	utils.AssertEqual(t, "example.com.to.resource", res)
	utils.AssertEqual(t, "path", hostParam)
	res1, hostParam1 := urlParser(domSegs, 1)
	utils.AssertEqual(t, "com.path.to.resource", res1)
	utils.AssertEqual(t, "example", hostParam1)

}

// go test -run Test_SubDomain_Middleware_Without_Parse -v
func Test_SubDomain_Middleware_Without_Parse(t *testing.T) {
	t.Parallel()
	app := fiber.New(fiber.Config{
		EnableSubDomains: true,
	})
	// app.Use(New(Config{
	// 	Host: &app.SubDomains,
	// }))
	app.Domain("author.localhost:3000").Get("/chain", func(c *fiber.Ctx) error {
		var err error
		c.SendString(c.Hostname())
		return err
	})
	app.Get("/", func(c *fiber.Ctx) error {
		c.SendString(c.Hostname())
		return nil
	})
	app.Use(New(Config{
		Host: &app.SubDomains,
	}))
	// test main domain
	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "http://localhost:3000", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertNotEqual(t, nil, resp)

	body, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertNotEqual(t, nil, body)
	utils.AssertEqual(t, "localhost:3000", string(body))

	// test sub domain
	subResp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "http://author.localhost:3000/chain", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertNotEqual(t, nil, resp)

	body2, err := io.ReadAll(subResp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertNotEqual(t, nil, body2)
	utils.AssertEqual(t, "author.localhost:3000", string(body2))
}

// go test -run Benchmark_SubDomain_Middleware_Without_Parse -v
func Benchmark_SubDomain_Middleware_Without_Parse(b *testing.B) {
	app := fiber.New(fiber.Config{
		EnableSubDomains: true,
	})
	app.Domain("author.localhost:3000").Get("/chain", func(c *fiber.Ctx) error {
		var err error
		c.SendString(c.Hostname())
		return err
	})
	app.Get("/", func(c *fiber.Ctx) error {
		c.SendString(c.Hostname())
		return nil
	})
	app.Use(New(Config{
		Host: &app.SubDomains,
	}))

	var subResp *http.Response
	var err error
	b.ResetTimer()
	for idx := 0; idx < b.N; idx++ {
		subResp, err = app.Test(httptest.NewRequest(fiber.MethodGet, "http://author.localhost:3000/chain", nil))
	}

	// test sub domain
	utils.AssertEqual(b, nil, err)
	utils.AssertNotEqual(b, nil, subResp)

	body2, err := io.ReadAll(subResp.Body)
	utils.AssertEqual(b, nil, err)
	utils.AssertNotEqual(b, nil, body2)
	utils.AssertEqual(b, "author.localhost:3000", string(body2))
}

// go test -run Test_SubDomain_Middleware_With_Parse -v
func Test_SubDomain_Middleware_With_Parse(t *testing.T) {
	t.Parallel()
	app := fiber.New(fiber.Config{
		EnableSubDomains: true,
	})
	app.Get("/", func(c *fiber.Ctx) error {
		c.SendString(c.Hostname())
		return nil
	})
	app.Domain(":author.localhost:3000").Get("/chain", func(c *fiber.Ctx) error {
		var err error
		c.SendString(c.RetriveValIdx0())
		return err
	})
	app.Use(New(Config{
		Host:  &app.SubDomains,
		Parse: true,
	}))

	// test main domain
	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "http://localhost:3000", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertNotEqual(t, nil, resp)

	body, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err)
	utils.AssertNotEqual(t, nil, body)
	utils.AssertEqual(t, "localhost:3000", string(body))

	// test sub domain
	subResp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "http://author.localhost:3000/chain", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertNotEqual(t, nil, subResp)

	body2, err := io.ReadAll(subResp.Body)
	t.Log(string(body2))
	utils.AssertEqual(t, nil, err)
	utils.AssertNotEqual(t, nil, body2)
	utils.AssertEqual(t, "author.localhost:3000", string(body2))
}

func Benchmark_SubDomain_Middleware_With_Parse(b *testing.B) {

	app := fiber.New(fiber.Config{
		EnableSubDomains: true,
	})
	app.Get("/", func(c *fiber.Ctx) error {
		c.SendString(c.Hostname())
		return nil
	})
	app.Domain(":author.localhost:3000").Get("/chain", func(c *fiber.Ctx) error {
		var err error
		c.SendString(c.Hostname())
		return err
	})
	app.Use(New(Config{
		Host:  &app.SubDomains,
		Parse: true,
	}))

	// test sub domain
	var subResp *http.Response
	var err error
	b.ResetTimer()
	for idx := 0; idx < b.N; idx++ {
		subResp, err = app.Test(httptest.NewRequest(fiber.MethodGet, "http://author.localhost:3000/chain", nil))
	}

	utils.AssertEqual(b, nil, err)
	utils.AssertNotEqual(b, nil, subResp)

	body2, err := io.ReadAll(subResp.Body)
	utils.AssertEqual(b, nil, err)
	utils.AssertNotEqual(b, nil, body2)
	utils.AssertEqual(b, "author.localhost:3000", string(body2))
}
