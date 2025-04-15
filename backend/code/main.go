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
	userManagementRouting(app)
	shiftRequestRouting(app)
	adminSubmissionRouting(app)
	userRequestsRounting(app)
}

func userManagementRouting(app *App) {
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

func shiftRequestRouting(app *App) {
	http.HandleFunc("GET /admin/requests", app.adminMiddleware(app.adminShowRequestsHandler))
	http.HandleFunc("GET /admin/requests/create", app.adminMiddleware(app.requestCreatePageHandler))
	http.HandleFunc("POST /admin/requests", app.adminMiddleware(app.createRequestHandler))
	// http.HandleFunc("GET /admin/requests/{request_id}")
}

func adminSubmissionRouting(app *App) {
	http.HandleFunc("GET /admin/requests/{request_id}/submissions", app.adminMiddleware(app.adminShowSubmissionsHandler))
	http.HandleFunc("GET /admin/requests/{request_id}/submissions/users/{user_id}", app.adminMiddleware(app.showUserSubmissionHandler))
	http.HandleFunc("GET /admin/requests/{request_id}/submissions/dates/{date}", app.adminMiddleware(app.showDateSubmissionHandler))
}

func userRequestsRounting(app *App) {
	http.HandleFunc("GET /requests", app.showRequestsHandler)
	// http.HandleFunc("GET /requests/{request_id}")
	http.HandleFunc("GET /requests/{request_id}/submit", app.submissionPageHandler)
	http.HandleFunc("POST /requests/{request_id}/submissions", app.submitShiftHandler)
	http.HandleFunc("GET /requests/{request_id}/submissions", app.showSubmissionsHandler)
}
