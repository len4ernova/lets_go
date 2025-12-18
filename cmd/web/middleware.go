package main

import (
	"fmt"
	"net/http"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")
		next.ServeHTTP(w, r)
	})
}
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)
		app.logger.Sugar().Infof("received request: ip %v, proto %v, method %v, uri %v", ip, proto, method, uri)
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Создаём отложенную функцию (которая всегда будет запускаться в случае
		// паники, поскольку Go разворачивает стек).
		defer func() {
			// Используйте встроенную функцию recover, чтобы проверить,
			// произошла ли паника. Если произошла...
			if err := recover(); err != nil {
				// Установите в ответе заголовок «Connection: close».
				w.Header().Set("Connection", "close")
				// Вызовите вспомогательный метод app.serverError, чтобы вернуть ответ 500
				// «Внутренняя ошибка сервера».
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
