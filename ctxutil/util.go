package ctxutil

import "context"

// GetRequestID get x-request-id from context
func GetRequestID(ctx context.Context) string {
	v, ok := ctx.Value("X-Request-ID").(string)
	if !ok {
		return ""
	}
	return v
}
