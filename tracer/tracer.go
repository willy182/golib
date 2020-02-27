package tracer

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	opentracing "github.com/opentracing/opentracing-go"
	ext "github.com/opentracing/opentracing-go/ext"
)

// Tracer for trace
type Tracer interface {
	NewChildContext() context.Context
	Context() context.Context
	InjectHTTPHeader(req *http.Request)
	Finish(tags map[string]interface{})
}

type opentracingTracer struct {
	ctx  context.Context
	span opentracing.Span
}

// StartTrace starting trace child span from parent span
func StartTrace(ctx context.Context, operationName string) Tracer {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		// init new span
		span, ctx = opentracing.StartSpanFromContext(ctx, operationName)
	} else {
		span = opentracing.GlobalTracer().StartSpan(operationName, opentracing.ChildOf(span.Context()))
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return &opentracingTracer{
		ctx:  ctx,
		span: span,
	}
}

// NewChildContext get context from child span
func (t *opentracingTracer) NewChildContext() context.Context {
	return opentracing.ContextWithSpan(t.ctx, t.span)
}

// Context get active context
func (t *opentracingTracer) Context() context.Context {
	return t.ctx
}

// InjectHTTPHeader to continue tracer to http request host
func (t *opentracingTracer) InjectHTTPHeader(req *http.Request) {
	ext.SpanKindRPCClient.Set(t.span)
	t.span.Tracer().Inject(
		t.span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
}

// Finish trace with tags data, must in deferred function
func (t *opentracingTracer) Finish(tags map[string]interface{}) {
	for k, v := range tags {
		t.span.SetTag(k, toString(v))
	}
	t.span.Finish()
}

// Log trace
func Log(ctx context.Context, event string, payload ...interface{}) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}

	if payload != nil {
		for _, p := range payload {
			if e, ok := p.(error); ok && e != nil {
				ext.Error.Set(span, true)
			}
			span.LogEventWithPayload(event, toString(p))
		}
	} else {
		span.LogEvent(event)
	}
}

// WithTrace closure with child context (deprecated)
func WithTrace(ctx context.Context, operationName string, tags map[string]interface{}, f func(context.Context)) {
	t := StartTrace(ctx, operationName)
	defer func() {
		t.Finish(tags)
	}()

	f(t.Context())
}

// WithTraceFunc functional with context and tags in function params
func WithTraceFunc(ctx context.Context, operationName string, fn func(context.Context, map[string]interface{})) {
	t := StartTrace(ctx, operationName)
	tags := make(map[string]interface{})

	defer func() {
		t.Finish(tags)
	}()

	fn(t.Context(), tags)
}

func toString(v interface{}) (s string) {
	switch val := v.(type) {
	case error:
		if val != nil {
			s = val.Error()
		}
	case string:
		s = val
	case int:
		s = strconv.Itoa(val)
	default:
		b, _ := json.Marshal(val)
		s = string(b)
	}
	return
}
