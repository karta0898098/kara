package grpc

import "go.uber.org/fx"

// Module for uber fx constructor
var Module = fx.Provide(
	NewGRPC,
)