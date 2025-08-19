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

	query, args, err := psql.Select("title", "width", "height").From("window").ToSql()
	E(err)
	rows, err := db.Query(query, args...)
	E(err)
	defer rows.Close()

	var w_title string = "No Title"
	var w_width int = 0
	var w_height int = 0
	if rows.Next() {
		err = rows.Scan(&w_title, &w_width, &w_height)
		E(err)
	}

	a := app.New()
	w := a.NewWindow(w_title)
	log.Println(w_width, w_height)
	w.Resize(fyne.NewSize(float32(w_width), float32(w_height)))
	

	query, args, err = psql.Select("title").From("label").ToSql()
	E(err)
	rows, err = db.Query(query, args...)
	E(err)
	defer rows.Close()

	w_content := []fyne.CanvasObject{};
	for rows.Next() {
		var l_title string = "No Title"
		err = rows.Scan(&l_title)
		E(err)
		w_content = append(w_content, widget.NewLabel(l_title))
	}

	w.SetContent(container.NewVBox(
		w_content...
	))

	w.ShowAndRun()
}