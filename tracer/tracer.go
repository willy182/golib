package tracer

import (
	"context"
	"errors"
	"log"
	"net/http"
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
	InjectHTTPHeader(req *http.Request)
	Finish(tags map[string]interface{})
}

type opentracingTracer struct {
	parentSpan, childSpan opentracing.Span
}

// StartTrace starting trace child span from parent span
func StartTrace(ctx context.Context, operationName string) (Tracer, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan == nil {
		return nil, errors.New("no span in context")
	}
	childSpan := opentracing.GlobalTracer().StartSpan(operationName, opentracing.ChildOf(parentSpan.Context()))
	return &opentracingTracer{parentSpan, childSpan}, nil
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

// Finish trace with tags data, must in defered function
func (t *opentracingTracer) Finish(tags map[string]interface{}) {
	for k, v := range tags {
		t.childSpan.SetTag(k, v)
	}
	t.childSpan.Finish()
}
