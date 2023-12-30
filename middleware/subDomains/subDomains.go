// note : to be added as the last handler
package subDomains

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	// Return new handler
	return func(c *fiber.Ctx) (err error) { //nolint:nonamedreturns // Uses recover() to overwrite the error
		// Don't execute middleware if Next returns true

		// if cfg.Next != nil && cfg.Next(c) {
		// 	return c.Next()
		// }

		hosts := *cfg.Host
		var host *fiber.SubApps
		var nxtCtxPtr = c.Context()
		if cfg.Parse {
			hostPath, hostParam := urlParser(c.Hostname(), cfg.ParamPosition)
			host = hosts[hostPath]
			c.SubdomainParam(hostParam)
		} else {
			host = hosts[c.Hostname()]
		}
		if host == nil {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			host.FiberApp.Handler()(nxtCtxPtr)
			return nil
		}
	}
}

func urlParser(hostname string, paramPosition int) (fullHost, hostParam string) {
	var builder strings.Builder
	domSegs := strings.Split(hostname, ".")

	for idx, seg := range domSegs {
		if idx == paramPosition-1 {
			hostParam = seg
			continue
		}
		builder.WriteString(seg)
		builder.WriteString(".")
	}

	fullHost = builder.String()
	if len(fullHost) > 0 {
		fullHost = fullHost[:len(fullHost)-1] // Remove the trailing dot
	}

	return
}
