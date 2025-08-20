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

func insertCanvasObject(arr []fyne.CanvasObject, index int, value fyne.CanvasObject) []fyne.CanvasObject {
	if len(arr) == index {
		return append(arr, value)
	}
	arr = append(arr[:index+1], arr[index:]...)
	arr[index] = value
	return arr
}

func database_connection() (sq.StatementBuilderType, *sql.DB) {
	db, err := sql.Open("sqlite3", "./gui.db")
	E(err)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Question)
	return psql, db
}

func insertCanvasObjectOfType(psql sq.StatementBuilderType, db *sql.DB, w_content *[]fyne.CanvasObject, window_id int) {
	query, args, err := psql.Select("L.title", "W.widget_order").From("widget AS W").Join("label AS L ON W.widget_id = L.widget_id").Where(sq.Eq{"W.window_id": window_id}).OrderBy("W.widget_order ASC").ToSql()
	E(err)
	rows, err := db.Query(query, args...)
	E(err)
	defer rows.Close()

	for rows.Next() {
		l_title := "No Title"
		widget_order := 0
		err = rows.Scan(&l_title, &widget_order)
		E(err)
		*w_content = insertCanvasObject(*w_content, widget_order, widget.NewLabel(l_title))
	}
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

	w_content := []fyne.CanvasObject{}
	insertCanvasObjectOfType(psql, db, &w_content, window_id)

	w.SetContent(container.NewVBox(
		w_content...,
	))

	w.ShowAndRun()
}
