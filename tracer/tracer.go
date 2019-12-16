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
	InjectHTTPHeader(req *http.Request)
	Finish(tags map[string]interface{})
}

type opentracingTracer struct {
	parentContext         context.Context
	parentSpan, childSpan opentracing.Span
	hasRootSpan           bool
}

// StartTrace starting trace child span from parent span
func StartTrace(ctx context.Context, operationName string) Tracer {
	hasRootSpan := true
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan == nil {
		// init new span
		parentSpan, _ = opentracing.StartSpanFromContext(ctx, operationName)
		hasRootSpan = false
	}
	childSpan := opentracing.GlobalTracer().StartSpan(operationName, opentracing.ChildOf(parentSpan.Context()))
	return &opentracingTracer{
		parentContext: ctx,
		parentSpan:    parentSpan,
		childSpan:     childSpan,
		hasRootSpan:   hasRootSpan,
	}
}

// NewChildContext get context from child span
func (t *opentracingTracer) NewChildContext() context.Context {
	return opentracing.ContextWithSpan(t.parentContext, t.childSpan)
}

// InjectHTTPHeader to continue tracer to http request host
func (t *opentracingTracer) InjectHTTPHeader(req *http.Request) {
	ext.SpanKindRPCClient.Set(t.childSpan)
	t.childSpan.Tracer().Inject(
		t.childSpan.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
}

// Finish trace with tags data, must in deferred function
func (t *opentracingTracer) Finish(tags map[string]interface{}) {
	for k, v := range tags {
		t.childSpan.SetTag(k, toString(v))
	}
	if !t.hasRootSpan {
		t.parentSpan.Finish()
	}
	t.childSpan.Finish()
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

// WithTrace closure with child context
func WithTrace(ctx context.Context, operationName string, tags map[string]interface{}, f func(context.Context)) {
	t := StartTrace(ctx, operationName)
	defer func() {
		t.Finish(tags)
	}()

	f(t.NewChildContext())
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
