package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

// LogRequests Мидлвайр, который логирует все входящие запросы
// Этот мидлвайр должен использоватся для каждого роута
func (m *Middleware) LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := time.Now()
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			m.logger.ErrWithContext(r.Context(), err, "")
		}
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		defer func() {
			byteDump, err := httputil.DumpRequest(r, false)

			if err != nil {
				m.logger.ErrWithContext(ctx, err, "")
			} else {
				msg := string(byteDump)

				if len(bodyBytes) > 0 {
					msg += fmt.Sprintf("\n Body: %s", string(bodyBytes))
				}

				msg += fmt.Sprintf("\n Duration: %v", time.Since(t))
				m.logger.DebugWithContext(ctx, msg)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
