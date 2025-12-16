package main

import (
	"github/len4ernova/lets_go/internal/models"
	"path/filepath"
	"text/template"
)

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
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
		// Создаём срез, содержащий пути к файлам для нашего базового шаблона, всех
		// частичных шаблонов и страницы
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}
		// Сгруппируйте файлы в набор шаблонов
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		// Добавьте набор шаблонов на карту, используя в качестве ключа название страницы
		// (например, 'home.tmpl').
		cache[name] = ts
		return cache, nil
	}
}
