package main

import (
	"database/sql"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
)

func E(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func database_connection() (sq.StatementBuilderType, *sql.DB) {
	db, err := sql.Open("sqlite3", "./gui.db")
	E(err)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Question)
	return psql, db
}

func main() {
	psql, db := database_connection()
	defer db.Close()

	query, args, err := psql.Select("window_id", "title", "width", "height").From("window").ToSql()
	E(err)
	rows, err := db.Query(query, args...)
	E(err)
	defer rows.Close()

	var window_id int = 0
	var w_title string = "No Title"
	var w_width int = 0
	var w_height int = 0
	if rows.Next() {
		err = rows.Scan(&window_id, &w_title, &w_width, &w_height)
		E(err)
	}

	a := app.New()
	w := a.NewWindow(w_title)
	w.Resize(fyne.NewSize(float32(w_width), float32(w_height)))

	query, args, err = psql.Select("L.title", "W.widget_order").From("widget AS W").Join("label AS L ON W.widget_id = L.widget_id").Where(sq.Eq{"W.window_id": window_id}).OrderBy("W.widget_order ASC").ToSql()
	E(err)
	rows, err = db.Query(query, args...)
	E(err)
	defer rows.Close()

	w_content := []fyne.CanvasObject{}
	for rows.Next() {
		var l_title string = "No Title"
		var widget_order int = 0
		err = rows.Scan(&l_title, &widget_order)
		E(err)
		w_content = append(w_content, widget.NewLabel(l_title))
	}

	w.SetContent(container.NewVBox(
		w_content...,
	))

	w.ShowAndRun()
}
