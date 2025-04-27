package context

import "backend/db"

type AppContext struct {
	DB db.DB
}
