package main

import (
	"compress/gzip" // Import the gzip package
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
)

// encodingSpec represents an Accept-Encoding header value with its quality factor.
type encodingSpec struct {
	name string
	q    float64
}

// compressResponseWriter wraps the http.ResponseWriter to provide compression.
// It holds a generic io.WriteCloser for either gzip or brotli.
type compressResponseWriter struct {
	io.Writer // This will be the compression writer (gzip or brotli)
	io.Closer // To close the compression writer
	http.ResponseWriter
}

// CompressionMiddleware selects the best encoding (Brotli or Gzip) supported by the client.
func CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		specs := parseAcceptEncoding(r.Header.Get("Accept-Encoding"))

		// Default to no compression
		chosenEncoding := ""
		highestQ := 0.0

		// Prefer Brotli ("br") if supported
		for _, spec := range specs {
			if (spec.name == "br" || spec.name == "*") && spec.q > highestQ {
				chosenEncoding = "br"
				highestQ = spec.q
				break // Found preferred Brotli
			}
		}

		// If Brotli not found or acceptable, check for Gzip ("gzip")
		if chosenEncoding == "" {
			highestQ = 0.0 // Reset Q for Gzip check
			for _, spec := range specs {
				// Only consider gzip if Brotli wasn't chosen
				// Also consider '*' if gzip is acceptable and q > 0
				if (spec.name == "gzip" || spec.name == "*") && spec.q > highestQ {
					// If '*' was already chosen for Brotli, don't override with gzip unless explicit
					if chosenEncoding == "br" && spec.name == "*" {
						continue
					}
					chosenEncoding = "gzip"
					highestQ = spec.q
					// Check if a specific gzip q-value is higher than a wildcard *
					if spec.name == "gzip" {
						break // Found explicit gzip, prioritize it over '*'
					}
				}
			}
		}

		// If no acceptable encoding found or q=0 for the chosen one, pass through
		if highestQ == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Set Vary header to indicate response varies based on Accept-Encoding
		w.Header().Add("Vary", "Accept-Encoding")

		var cw io.WriteCloser
		var writer io.Writer

		switch chosenEncoding {
		case "br":
			w.Header().Set("Content-Encoding", "br")
			brWriter := brotli.NewWriter(w)
			cw = brWriter
			writer = brWriter
		case "gzip":
			w.Header().Set("Content-Encoding", "gzip")
			// Use gzip.NewWriterLevel for potentially better control, or just NewWriter
			gzWriter, _ := gzip.NewWriterLevel(w, gzip.DefaultCompression) // Or e.g., gzip.BestSpeed
			cw = gzWriter
			writer = gzWriter
		default:
			// Should not happen due to highestQ check, but fallback just in case
			next.ServeHTTP(w, r)
			return
		}

		// Wrap the response writer
		crw := &compressResponseWriter{
			Writer:         writer,
			Closer:         cw,
			ResponseWriter: w,
		}
		defer crw.Close() // Ensure the writer is closed

		next.ServeHTTP(crw, r)
	})
}

// parseAcceptEncoding parses the Accept-Encoding header string into a slice of encodingSpecs.
// (This function remains the same as the previous version)
func parseAcceptEncoding(header string) []encodingSpec {
	var specs []encodingSpec
	parts := strings.Split(header, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		params := strings.Split(part, ";")
		name := strings.TrimSpace(params[0])
		q := 1.0 // Default quality factor

		if len(params) > 1 {
			for _, p := range params[1:] {
				p = strings.TrimSpace(p)
				if strings.HasPrefix(p, "q=") {
					qStr := strings.TrimPrefix(p, "q=")
					qVal, err := strconv.ParseFloat(qStr, 64)
					if err == nil && qVal >= 0 && qVal <= 1 { // q must be between 0.0 and 1.0
						q = qVal
					} else {
						q = 0 // Invalid q-value means unacceptable
					}
					break // Only one q value is expected per encoding
				}
			}
		}

		// Ignore identity encoding explicitly unless it's the only one or q=0
		if name == "identity" && q > 0 {
			continue
		}
		// Ignore empty names
		if name == "" {
			continue
		}

		specs = append(specs, encodingSpec{name: name, q: q})
	}

	// Sort by q-value (descending) - useful for debugging or more complex selection logic
	sort.SliceStable(specs, func(i, j int) bool {
		// Prioritize specific encodings over wildcard '*' if q is equal
		if specs[i].q == specs[j].q {
			if specs[i].name == "*" {
				return false
			}
			if specs[j].name == "*" {
				return true
			}
		}
		return specs[i].q > specs[j].q
	})

	return specs
}

// Write delegates the Write call to the underlying compression writer.
func (w *compressResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Close closes the underlying compression writer.
func (w *compressResponseWriter) Close() error {
	if w.Closer != nil {
		return w.Closer.Close()
	}
	return nil
}

// Flush implements the http.Flusher interface.
func (w *compressResponseWriter) Flush() {
	// Flush the compression writer if it supports it.
	if flusher, ok := w.Writer.(http.Flusher); ok {
		flusher.Flush()
	}
	// Also flush the underlying ResponseWriter if it supports flushing.
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Header returns the header map of the underlying ResponseWriter.
func (w *compressResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// WriteHeader sends an HTTP response header with the provided status code.
func (w *compressResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}
