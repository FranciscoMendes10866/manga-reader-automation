package helpers

import "github.com/go-chi/cors"

var corsOrigin = "*"
var AllowedHeaderValue = "superSecretKey"

var CorsConfig = cors.Options{
	AllowedOrigins:   []string{corsOrigin},
	AllowedMethods:   []string{"POST"},
	AllowedHeaders:   []string{"x-dango-manga-key"},
	AllowCredentials: false,
}
