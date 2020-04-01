package tracer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"
	ext "github.com/opentracing/opentracing-go/ext"
)

// Tracer for trace
type Tracer interface {
	NewChildContext() context.Context
	Context() context.Context
	Tags() map[string]interface{}
	InjectHTTPHeader(req *http.Request)
	Finish(tags ...map[string]interface{})
}

type opentracingTracer struct {
	ctx  context.Context
	span opentracing.Span
	tags map[string]interface{}
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

// NewChildContext get context from child span (deprecated)
func (t *opentracingTracer) NewChildContext() context.Context {
	return opentracing.ContextWithSpan(t.ctx, t.span)
}

// Context get active context
func (t *opentracingTracer) Context() context.Context {
	return t.ctx
}

// Tags create tags in tracer span
func (t *opentracingTracer) Tags() map[string]interface{} {
	t.tags = make(map[string]interface{})
	return t.tags
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

// Finish trace with additional tags data, must in deferred function
func (t *opentracingTracer) Finish(tags ...map[string]interface{}) {
	defer t.span.Finish()

	if tags != nil && t.tags == nil {
		t.tags = make(map[string]interface{})
	}

	for _, tag := range tags {
		for k, v := range tag {
			t.tags[k] = v
		}
	}

	for k, v := range t.tags {
		t.span.SetTag(k, toString(v))
	}
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
	defer t.Finish()

	fn(t.Context(), t.Tags())
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

// GetTraceID func
func GetTraceID(ctx context.Context) string {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	traceID := fmt.Sprintf("%+v", span)
	splits := strings.Split(traceID, ":")
	if len(splits) > 0 {
		return splits[0]
	}

	return traceID
}
