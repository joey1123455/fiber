package subDomains

import "github.com/gofiber/fiber/v2"

// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip middleware.
	// Optional. Default: nil
	Next func(*fiber.Ctx) bool

	// Parse is a bool value that makes use of the params parser
	Parse bool

	// ParamPosition is the index determining which segment of the url contains
	// the param eg ":authors.blog.example.com" in this case the param is authors
	// denoted by the semi colon here its param position is the first element
	//		app.Use(subDomains.New(
	//			fiber.Config{
	//				Parse: true
	//				ParamPosition: 1
	//				}
	//		))
	// eg "blog.:authors.com" its param position is two
	//		app.Use(subDomains.New(
	//			fiber.Config{
	//				Parse: true
	//				ParamPosition: 2
	//				}
	//		))
	ParamPosition int

	// Hosts pointer to a map containig subdomains
	Host *fiber.Hosts
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Parse:         false,
	ParamPosition: 1,
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		panic("ptr to app.Domains must be passed into the config as config.Host")
	}

	// Override default config
	cfg := config[0]

	// set default values
	if cfg.ParamPosition == 0 && cfg.Parse {
		cfg.ParamPosition = ConfigDefault.ParamPosition
	}

	return cfg
}
