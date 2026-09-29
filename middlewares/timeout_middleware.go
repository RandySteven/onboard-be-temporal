package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (s *ServerMiddleware) TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeoutTime := 30
		if raw := os.Getenv("SERVER_TIMEOUT"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil {
				log.Printf("Could not parse timeout value: %v", err)
			} else {
				timeoutTime = parsed
			}
		}

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeoutTime)*time.Second)
		defer cancel()

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		if err := ctx.Err(); err != nil && errors.Is(err, context.DeadlineExceeded) {
			log.Println(err)
		}
	})
}
