package middleware

import (
	"github.com/gin-contrib/cors"
	"net/http"
)

func NewCors() {
	cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
		},
		AllowHeaders: []string{
			"*",
		},
		ExposeHeaders: []string{
			"*",
		},
	})
}
