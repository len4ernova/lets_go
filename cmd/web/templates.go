package main

import (
	"github/len4ernova/lets_go/internal/models"
	"github/len4ernova/lets_go/ui"
	"io/fs"
	"path/filepath"
	"text/template"
	"time"
)

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string // Add a CSRFToken field.
}

// пользовательская ф-ия
func humanDate(t time.Time) string {
	// Return the empty string if time has the zero value.
	if t.IsZero() {
		return ""
	}
	// Convert the time to UTC before formatting it.
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// Инициализируем объект template.FuncMap и сохраняем его в глобальной переменной.
// По сути, это карта со строковыми ключами, которая служит для поиска по именам наших пользовательских функций шаблона и самих функций.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}
	// Используйте функцию fs.Glob(), чтобы получить список всех путей к файлам,
	// соответствующих шаблону "./ui/html/pages/*.tmpl". По сути, это даст
	// нам список всех путей к файлам для шаблонов страниц нашего приложения,
	// например: [ui/html/pages/home.tmpl, ui/html/pages/view.tmpl]
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	// Перебираем пути к файлам на странице один за другим.
	for _, page := range pages {
		// Извлечь имя файла (например, 'home.tmpl') из полного пути к файлу
		// // и присвоить его переменной name.
		name := filepath.Base(page)

		// Создайте срез, содержащий шаблоны путей к файлам для шаблонов, которые мы хотим проанализировать
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}
		//  Используйте ParseFS() вместо ParseFiles() для анализа файлов шаблонов
		//  из встроенной файловой системы ui.Files.

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Добавьте набор шаблонов на карту, используя в качестве ключа название страницы
		// (например, 'home.tmpl').
		cache[name] = ts
	}
	return cache, nil
}
