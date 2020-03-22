package main

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rotscher/cruds/internal/cards"
	"github.com/rotscher/cruds/internal/route"
	"github.com/uber/jaeger-lib/metrics"
	"io"
	"log"
	"net/http"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func main() {
	handler := &route.RegexpHandler{}

	//https://stackoverflow.com/questions/6564558/wildcards-in-the-pattern-for-http-handlefunc
	handler.HandleFunc("/cards", cards.GetAll).Methods("GET")
	handler.HandleFunc("/cards/{cardId}", cards.GetById).Methods("GET")
	handler.HandleFunc("/cards", cards.Insert).Methods("POST")

	closer := initTracer()
	defer closer.Close()
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func initTracer() io.Closer {
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := jaegercfg.Configuration{
		ServiceName: "cruds",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 0,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, _ := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(tracer)
	return closer
	// continue main()
}
