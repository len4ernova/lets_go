package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

// The serverError helper writes a log entry at Error level (including the request
// method and URI as attributes), then sends a generic 500 Internal Server Error
// response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Извлекаем из кэша соответствующий набор шаблонов по имени страницы
	// (например, 'home.tmpl'). Если в кэше нет записи с указанным
	//  именем, создаем новую ошибку и вызываем вспомогательный метод
	// serverError(), который мы создали ранее, и возвращаем результат.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}
	// Initialize a new buffer.
	buf := new(bytes.Buffer)
	// Запишите шаблон в буфер, а не сразу в
	// http.ResponseWriter. Если возникнет ошибка, вызовите нашу вспомогательную функцию serverError()
	// и затем вернитесь.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Если шаблон записан в буфер без ошибок, можно смело
	// записывать код состояния HTTP в http.ResponseWriter.
	w.WriteHeader(status)
	// Запишите содержимое буфера в http.ResponseWriter.
	// Примечание: это ещё один случай, когда мы передаём http.ResponseWriter в функцию,
	// которая принимает io.Writer
	buf.WriteTo(w)
}

// Создайте вспомогательную функцию newTemplateData(), которая возвращает структуру templateData,
// инициализированную текущим годом.
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}
