package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// mux.HandleFunc("GET /{$}", app.home)
	// mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	// mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	// mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Update these routes to use the new dynamic middleware chain followed by
	// the appropriate handler function. Note that because the alice ThenFunc()
	// method returns a http.Handler (rather than a http.HandlerFunc) we also
	// need to switch to registering the route using the mux.Handle() method.

	// Unprotected application routes using the "dynamic" middleware chain.
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected (authenticated-only) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// Передайте servemux в качестве параметра next промежуточному программному обеспечению commonHeaders. // Поскольку commonHeaders — это просто функция, а функция возвращает
	// http.Handler, нам больше ничего не нужно делать.
	// 	mux := http.NewServeMux()
	// mux.Handle("GET /{$}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.home)))
	// mux.Handle("GET /snippet/view/:id", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetView)))

	// mux.Handle("POST /snippet/create", app.sessionManager.LoadAndSave(app.requireAuthentication(http.HandlerFunc(app.snippetCreat)))
	// ... etc
	//return app.recoverPanic(app.logRequest(commonHeaders(mux)))

	// Создаём цепочку промежуточного программного обеспечения, содержащую наше «стандартное» промежуточное ПО,
	// которое будет использоваться для каждого запроса, поступающего в наше приложение
	standart := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	// Возвращает «стандартную» цепочку промежуточного программного обеспечения, за которой следует servemux.
	return standart.Then(mux)
}
