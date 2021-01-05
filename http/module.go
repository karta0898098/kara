package http

import "go.uber.org/fx"

var Module = fx.Provide(
	NewEcho,
)
