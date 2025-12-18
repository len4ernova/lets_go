package main

import (
	"github/len4ernova/lets_go/internal/models"
	"path/filepath"
	"text/template"
	"time"
)

type templateData struct {
	CurrentYear int
	//Snippet     models.Snippet
	//Snippets    []models.Snippet
	Work  models.Work
	Works []models.Work
}

// пользовательская ф-ия
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Инициализируем объект template.FuncMap и сохраняем его в глобальной переменной.
// По сути, это карта со строковыми ключами, которая служит для поиска по именам наших пользовательских функций шаблона и самих функций.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}
	// Используйте функцию filepath.Glob(), чтобы получить список всех путей к файлам,
	// соответствующих шаблону "./ui/html/pages/*.tmpl". По сути, это даст
	// нам список всех путей к файлам для шаблонов страниц нашего приложения,
	// например: [ui/html/pages/home.tmpl, ui/html/pages/view.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	// Перебираем пути к файлам на странице один за другим.
	for _, page := range pages {
		// Извлечь имя файла (например, 'home.tmpl') из полного пути к файлу
		// // и присвоить его переменной name.
		name := filepath.Base(page)
		// Шаблон.FuncMap должен быть зарегистрирован в наборе шаблонов до вызова метода ParseFiles().
		// Это означает, что мы должны использовать template.New() для создания пустого набора шаблонов,
		// использовать метод Funcs() для регистрации шаблона.FuncMap, а затем проанализировать файл как обычно
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		// // Сгруппируйте файлы в набор шаблонов
		// ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}
		// Вызовите ParseGlob() *для этого набора шаблонов*, чтобы добавить частичные шаблоны.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}
		// Чтобы добавить шаблон страницы, вызовите ParseFiles() *для этого набора шаблонов*.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Добавьте набор шаблонов на карту, используя в качестве ключа название страницы
		// (например, 'home.tmpl').
		cache[name] = ts
	}
	return cache, nil
}
