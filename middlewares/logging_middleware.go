package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/RandySteven/onboard-be/enums"
	"github.com/google/uuid"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	Body       *bytes.Buffer
	StatusCode int
}

func (rw *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriterWrapper) Write(b []byte) (int, error) {
	rw.Body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (s *ServerMiddleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), enums.RequestID, uuid.NewString())
		requestTime := time.Now()

		var body []byte
		if r.Method != http.MethodGet {
			var err error
			body, err = io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Failed to read request body: %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		rw := &ResponseWriterWrapper{
			ResponseWriter: w,
			Body:           &bytes.Buffer{},
			StatusCode:     http.StatusOK,
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

		log.Printf("%s %s %s %d %s", requestTime.Format(time.RFC3339), r.Method, r.URL.Path, rw.StatusCode, ctx.Value(enums.RequestID))
		if len(body) > 0 {
			var requestBody any
			if err := json.Unmarshal(body, &requestBody); err == nil {
				if m, ok := requestBody.(map[string]any); ok {
					if _, exists := m["password"]; exists {
						m["password"] = "[redacted]"
					}
					if _, exists := m["token"]; exists {
						m["token"] = "[redacted]"
					}
				}
				log.Printf("request: %v", requestBody)
			}
		}
	})
}
