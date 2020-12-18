package metrics

import (
	"context"
	"strings"
)

const DefaultTraceID = "X-Request-ID"

func GetTraceID(ctx context.Context) string {
	v,ok := ctx.Value(DefaultTraceID).(string)
	if ok {
		return v
	}

	v, ok = ctx.Value(strings.ToLower(DefaultTraceID)).(string)
	if ok && v != "" {
		return v
	}

	return ""
}


