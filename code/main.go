package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	db *sql.DB
}

func main() {
	fmt.Println("Server Start")

	// access db
	db, err := sql.Open("sqlite3", "./data/shift_webapp.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	app := &App{db: db}

	routing(app)
	http.ListenAndServe(":80", nil)
}

func routing(app *App) {
	http.HandleFunc("/home", app.homeHandler)
	http.HandleFunc("/login", app.loginHandler)
	http.HandleFunc("/logout", app.logoutHandler)

	// admin
	http.HandleFunc("/admin/home", app.adminHandler(app.adminHomeHandler))
	http.HandleFunc("/admin/register", app.adminHandler(app.adminRegisterHandler))
	http.HandleFunc("/admin/users", app.adminHandler(app.adminUsersHandler))
	http.HandleFunc("/admin/delete", app.adminHandler(app.adminDeleteHandler))
}
