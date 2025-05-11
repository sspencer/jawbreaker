package main

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func slogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t1 := time.Now()

		remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			remoteIP = r.RemoteAddr
		}

		defer func() {
			slog.Info("Served request",
				slog.String("method", r.Method),
				slog.String("proto", r.Proto),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", remoteIP),
				slog.Duration("duration", time.Since(t1)),
				slog.Int("status", ww.Status()),
				slog.Int("size", ww.BytesWritten()),
			)
		}()
		next.ServeHTTP(ww, r)
	})
}

// recoverMiddleware catches panics and logs them using slog.
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				slog.Error("Panic recovered",
					slog.Any("panic", rvr),
					slog.String("stack", string(debug.Stack())),
					slog.String("requestID", middleware.GetReqID(r.Context())),
				)

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func compressionMiddleware(r chi.Router) {
	compressor := middleware.NewCompressor(5) // Default compression level for gzip
	compressor.SetEncoder("br", func(w io.Writer, level int) io.Writer {
		return brotli.NewWriterLevel(w, brotli.DefaultCompression)
	})
	r.Use(compressor.Handler)
}
