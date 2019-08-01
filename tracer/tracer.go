package tracer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	ext "github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

// InitOpenTracing with agent and service name
func InitOpenTracing(agentHost, serviceName string) error {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  agentHost,
		},
	}
	tracer, _, err := cfg.New(serviceName, config.Logger(jaeger.StdLogger))
	if err != nil {
		log.Printf("ERROR: cannot init opentracing connection: %v\n", err)
		return err
	}
	opentracing.SetGlobalTracer(tracer)
	return nil
}

// Tracer for trace
type Tracer interface {
	NewChildContext() context.Context
	InjectHTTPHeader(req *http.Request)
	Finish(tags map[string]interface{})
}

type opentracingTracer struct {
	parentSpan, childSpan opentracing.Span
}

// StartTrace starting trace child span from parent span
func StartTrace(ctx context.Context, operationName string) Tracer {
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan == nil {
		// init new span
		parentSpan, _ = opentracing.StartSpanFromContext(context.Background(), operationName)
	}
	childSpan := opentracing.GlobalTracer().StartSpan(operationName, opentracing.ChildOf(parentSpan.Context()))
	return &opentracingTracer{parentSpan, childSpan}
}

func (t *opentracingTracer) NewChildContext() context.Context {
	return opentracing.ContextWithSpan(context.Background(), t.childSpan)
}

// InjectHTTPHeader to continue tracer in http request host
func (t *opentracingTracer) InjectHTTPHeader(req *http.Request) {
	ext.SpanKindRPCClient.Set(t.childSpan)
	t.parentSpan.Tracer().Inject(
		t.parentSpan.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
}

// Finish trace with tags data, must in deferred function
func (t *opentracingTracer) Finish(tags map[string]interface{}) {
	for k, v := range tags {
		t.childSpan.SetTag(k, v)
	}
	t.childSpan.Finish()
}

// Log trace
func Log(ctx context.Context, event string, payload ...interface{}) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		if payload != nil {
			for _, p := range payload {
				span.LogEventWithPayload(event, toString(p))
			}
		} else {
			span.LogEvent(event)
		}
	}
}

// WithTrace closure
func WithTrace(ctx context.Context, operationName string, tags map[string]interface{}, f func()) {
	t := StartTrace(ctx, operationName)
	defer func() {
		if r := recover(); r != nil {
			Log(ctx, operationName, fmt.Errorf("panic: %v", r))
		}
		t.Finish(tags)
	}()

	f()
}

func toString(v interface{}) (s string) {
	switch val := v.(type) {
	case error:
		s = val.Error()
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
