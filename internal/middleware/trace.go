package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const tracerName = "go_sample_code/middleware"

// Trace creates a Fiber middleware that starts an OTel span for every request,
// propagates W3C trace context from incoming headers, and injects the span into
// the request context so downstream handlers and the Logger middleware can use it.
func Trace(tp *sdktrace.TracerProvider) fiber.Handler {
	tracer := tp.Tracer(tracerName)
	propagator := otel.GetTextMapPropagator()

	return func(c *fiber.Ctx) error {
		// Extract W3C trace context from incoming request headers
		ctx := propagator.Extract(c.UserContext(), propagation.HeaderCarrier(c.GetReqHeaders()))

		// Span name: "HTTP POST /api/users"
		spanName := c.Method() + " " + utils.CopyString(c.Path())

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethodKey.String(c.Method()),
				semconv.HTTPTargetKey.String(c.OriginalURL()),
				semconv.HTTPRouteKey.String(c.Route().Path),
				semconv.HTTPSchemeKey.String(c.Protocol()),
				semconv.NetHostNameKey.String(c.Hostname()),
				semconv.HTTPUserAgentKey.String(string(c.Request().Header.UserAgent())),
				semconv.NetSockPeerAddrKey.String(c.IP()),
				attribute.String("http.request_id", c.GetRespHeader(fiber.HeaderXRequestID)),
			),
		)
		defer span.End()

		// Inject span context back into the fiber context
		c.SetUserContext(ctx)

		// Inject trace context into outgoing response headers (for downstream services)
		propagator.Inject(ctx, propagation.HeaderCarrier(c.GetRespHeaders()))

		err := c.Next()

		// Record response status
		status := c.Response().StatusCode()
		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(status))

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else if status >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, http.StatusText(status))
		} else {
			span.SetStatus(codes.Ok, "")
		}

		return err
	}
}
