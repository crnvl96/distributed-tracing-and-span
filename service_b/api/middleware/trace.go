package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func CallTrace(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		carrier := propagation.HeaderCarrier(r.Header)
		ctx := r.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
		ctx, span := otel.GetTracerProvider().Tracer("service_b").Start(ctx, "service_b")
		defer span.End()
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
