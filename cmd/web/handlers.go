package main

import (
	"errors"
	"fmt"
	"github/len4ernova/lets_go/internal/models"
	"strconv"

	// "html/template"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// panic("oops! something went wrong") // Deliberate panic

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Вызовите вспомогательную функцию newTemplateData(), чтобы получить структуру templateData, содержащую
	// данные по умолчанию (которые на данный момент представляют собой только текущий год), и добавьте к ней
	// фрагмент данных
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, r, http.StatusOK, "view.tmpl", data)
}
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Сначала мы вызываем r.ParseForm(), который добавляет все данные из тела запроса POST
	// в карту r.PostForm. То же самое происходит с запросами PUT и PATCH
	// . Если возникают какие-либо ошибки, мы используем вспомогательную функцию app.ClientError()
	// для отправки пользователю ответа 400 Bad Request
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Используйте метод r.PostForm.Get() для извлечения заголовка и содержимого
	// из карты r.PostForm.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	// Метод r.PostForm.Get() всегда возвращает данные формы в виде *строки*.
	// Однако мы ожидаем, что значение expires будет числом, и хотим
	//  представить его в нашем коде Go как целое число. Поэтому нам нужно
	// вручную преобразовать данные формы в целое число с помощью
	// strconv.Atoi(), и мы отправим ответ с кодом 400 «Неверный запрос»,
	// если преобразование не удастся
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// Change the signature of the snippetCreate handler so it is defined as a method
// against *application.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}
