package es

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type (
	Config struct {
		Address    []string
		Username   string
		Password   string
		MaxRetries int
	}
	Es struct {
		*es8.Client
	}
	esTransport struct{}
)

func NewEs(config *Config) (*Es, error) {
	c := es8.Config{
		Addresses:  config.Address,
		Username:   config.Username,
		Password:   config.Password,
		MaxRetries: config.MaxRetries,
		Transport:  &esTransport{},
	}
	client, err := es8.NewClient(c)
	if err != nil {
		return nil, err
	}

	return &Es{Client: client}, nil
}

func MustNewEs(config *Config) *Es {
	es, err := NewEs(config)
	if err != nil {
		panic(err)
	}
	return es
}

func (t *esTransport) RoundTrip(req *http.Request) (rsp *http.Response, err error) {
	var (
		ctx        = req.Context()
		span       oteltrace.Span
		startTime  = time.Now()
		propagator = otel.GetTextMapPropagator()
		indexName  = strings.Split(req.URL.RequestURI(), "/")[1]
		tracer     = trace.TracerFromContext(ctx)
	)
	//初始化trace，获取ctx和span
	spanName := req.URL.Path
	fmt.Println("span name is ", spanName)
	ctx, span = tracer.Start(ctx, spanName,
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		oteltrace.WithAttributes(semconv.HTTPClientAttributesFromHTTPRequest(req)...),
	)
	fmt.Println(ctx)
	defer func() {
		// metric 记录req  计数error
		metricClientReqDur.Observe(time.Since(startTime).Milliseconds(), indexName)

		metricClientReqErrTotal.Inc(indexName, strconv.FormatBool(err != nil))

		span.End()
	}()

	// create a req with context which conclude trace
	req = req.WithContext(ctx)
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
	//在调用defaultTransport的roudtrip方法前初始化metric和trace
	rsp, err = http.DefaultTransport.RoundTrip(req)
	if err != nil {
		//若错误，trace记录error
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
	fmt.Println("span开始正常记录")

	span.SetAttributes(semconv.DBSQLTableKey.String(indexName))
	span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(rsp.StatusCode)...)
	fmt.Println("span正常记录中")
	span.SetStatus(semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(rsp.StatusCode, oteltrace.SpanKindClient))

	return

}
