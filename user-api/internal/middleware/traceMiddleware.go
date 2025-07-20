package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

// TraceMiddleware 统一处理trace ID响应头设置
func TraceMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取trace ID并设置到响应头
		spanCtx := trace.SpanContextFromContext(r.Context())
		if spanCtx.IsValid() {
			traceID := spanCtx.TraceID().String()
			w.Header().Set("X-Trace-ID", traceID)
			w.Header().Set("X-Request-ID", traceID)

			// 也可以设置span ID
			spanID := spanCtx.SpanID().String()
			w.Header().Set("X-Span-ID", spanID)
		}

		// 继续处理请求
		next(w, r)
	}
}
