package main

import "net/http"

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Передайте servemux в качестве параметра next промежуточному программному обеспечению commonHeaders. // Поскольку commonHeaders — это просто функция, а функция возвращает
	// http.Handler, нам больше ничего не нужно делать.
	//return app.recoverPanic(app.logRequest(commonHeaders(mux)))

	// Создаём цепочку промежуточного программного обеспечения, содержащую наше «стандартное» промежуточное ПО,
	// которое будет использоваться для каждого запроса, поступающего в наше приложение
	standart := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	// Возвращает «стандартную» цепочку промежуточного программного обеспечения, за которой следует servemux.
	return standart.Then(mux)
}
