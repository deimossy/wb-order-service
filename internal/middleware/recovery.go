package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func Recovery(next http.Handler, logg *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logg.Error("panic recovered",
					zap.Any("error", rec),
					zap.String("url", r.URL.Path),
				)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
