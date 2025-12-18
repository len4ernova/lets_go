package main

import (
	"database/sql"
	"errors"
	"github/len4ernova/lets_go/internal/models"
	"strconv"

	// "html/template"
	"net/http"
)

// home - домашняя страница, с последними 10 работами
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	works, err := app.works.ListWorks(false)
	if err != nil {
		app.serverError(w, r, err)
	}

	data := app.newTemplateData(r)
	data.Works = works
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

// GetWorkList - список всех работ
func (app *application) GetListWorks(w http.ResponseWriter, r *http.Request) {
	works, err := app.works.ListWorks(true)
	if err != nil {
		app.serverError(w, r, err)
	}

	data := app.newTemplateData(r)
	data.Works = works
	app.render(w, r, http.StatusOK, "all_works.tmpl", data)
}

// syncGlab - синхронизация данных с Gitlab
func (app *application) syncGlab(w http.ResponseWriter, r *http.Request) {
	// todo добавить форму с токен
	token := ""
	ip := ""

	groups, err := app.works.Sync(w, r, ip, token)
	if err != nil {
		app.serverError(w, r, err)
	}

	for _, item := range groups {
		_, err := app.works.GetWorkGLabID(item.Id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				_, err := app.works.InsertWork(item.Id, item.Name, item.Path, item.CreatedAt, item.Description, true)
				if err != nil {
					app.serverError(w, r, err)
					continue
				}
			} else {
				app.serverError(w, r, err)
				continue // ? or прерывание
			}
		}
		err = app.works.UpdateWork(item.Id, item.Name, item.Path, item.CreatedAt, item.Description)
		if err != nil {
			app.serverError(w, r, err)
			continue
		}
	}
	// worksGlab, err := app.works.LatestWork(true)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }
	// fmt.Println("обработана отправка в БД")
	// var result models.Glab
	// result.Result = "успешно"
	// result.UpdatedAt = time.Now().Format("20060102")
	// fmt.Println("Result", result)

	// data := app.newTemplateData(r)
	// data.GlabRes = result
	// app.render(w, r, http.StatusOK, "sync.tmpl", data)

	// TODO render
}

// workView - вывод работы по id
func (app *application) workView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	work, err := app.works.GetWork(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Work = work
	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

// func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
// 	// Create some variables holding dummy data. We'll remove these later on
// 	// during the build.
// 	title := "O snail"
// 	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
// 	expires := 7
// 	// Pass the data to the SnippetModel.Insert() method, receiving the
// 	// ID of the new record back.
// 	id, err := app.snippets.Insert(title, content, expires)
// 	if err != nil {
// 		app.serverError(w, r, err)
// 		return
// 	}
// 	// Redirect the user to the relevant page for the snippet.
// 	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
// }

// Change the signature of the snippetCreate handler so it is defined as a method
// against *application.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// workCreate - создать работу.
func (app *application) workCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}
