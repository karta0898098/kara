package tracer

type tracer string

func (t tracer) ToString() string  {
	return string(t)
}

const (
	TraceIDKey tracer = "trace_id"
)
