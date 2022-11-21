package middleware

import (
	"fmt"
	"net/http"
	"reflect"
)

// Мидлвер, который отлавливает панику и не дает процессу сервиса упасть с ошибкой.
// Этот мидлвайр должен использоватся для каждого роута
func (m *Middleware) RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errType := reflect.TypeOf(err)
				fmt.Println("recovered: ", errType)

				if convertedError, ok := err.(interface{ Error() string }); ok {
					fmt.Println("Error: ", convertedError.Error())
					m.logger.ErrWithContext(r.Context(), convertedError, "Internal error")
				}
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
