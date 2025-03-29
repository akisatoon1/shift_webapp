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
	http.HandleFunc("GET /home", app.homeHandler)
	http.HandleFunc("GET /login", app.loginHandler)
	http.HandleFunc("POST /login", app.loginHandler)
	http.HandleFunc("POST /logout", app.logoutHandler)

	// adminのみアクセスできる
	// /admin/*という形の存在しないリソースは、adminには404notfound, 他へはアクセスを禁じる
	http.HandleFunc("GET /admin/", app.adminMiddleware(func(w http.ResponseWriter, r *http.Request, usr user) {
		http.NotFound(w, r)
	}))
	http.HandleFunc("GET /admin/home", app.adminMiddleware(app.adminHomeHandler))
	http.HandleFunc("GET /admin/register", app.adminMiddleware(app.adminRegisterHandler))
	http.HandleFunc("POST /admin/register", app.adminMiddleware(app.adminRegisterHandler))
	http.HandleFunc("GET /admin/users", app.adminMiddleware(app.adminUsersHandler))
	http.HandleFunc("POST /admin/delete", app.adminMiddleware(app.adminDeleteHandler))
}
