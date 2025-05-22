package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Sqlite3DBはDBインターフェースのsqlite3実装
// フィールドConnは*sql.DB型
type Sqlite3DB struct {
	Conn *sql.DB
}

// NewSqlite3DBはSqlite3DBの初期化関数
func NewSqlite3DB(dataSourceName string) (*Sqlite3DB, error) {
	conn, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Sqlite3DB{Conn: conn}, nil
}

// アプリ終了時にDBを閉じるため
func (db *Sqlite3DB) Close() error {
	return db.Conn.Close()
}

// ユーザーIDでユーザーを取得
func (db *Sqlite3DB) GetUserByID(id int) (User, error) {
	var user User
	row := db.Conn.QueryRow("SELECT id, login_id, password, name, role, created_at FROM users WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.LoginID, &user.Password, &user.Name, &user.Role, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// login_idでユーザーを取得
func (db *Sqlite3DB) GetUserByLoginID(loginID string) (User, error) {
	var user User
	row := db.Conn.QueryRow("SELECT id, login_id, password, name, role, created_at FROM users WHERE login_id = ?", loginID)
	err := row.Scan(&user.ID, &user.LoginID, &user.Password, &user.Name, &user.Role, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// 全リクエストを取得
func (db *Sqlite3DB) GetRequests() ([]Request, error) {
	rows, err := db.Conn.Query("SELECT id, creator_id, start_date, end_date, deadline, created_at FROM requests")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []Request
	for rows.Next() {
		var req Request
		err := rows.Scan(&req.ID, &req.CreatorID, &req.StartDate, &req.EndDate, &req.Deadline, &req.CreatedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

// 指定リクエストIDのリクエストを取得
func (db *Sqlite3DB) GetRequestByID(id int) (Request, error) {
	var req Request
	row := db.Conn.QueryRow("SELECT id, creator_id, start_date, end_date, deadline, created_at FROM requests WHERE id = ?", id)
	err := row.Scan(&req.ID, &req.CreatorID, &req.StartDate, &req.EndDate, &req.Deadline, &req.CreatedAt)
	if err != nil {
		return Request{}, err
	}
	return req, nil
}

// 指定リクエストIDのエントリー一覧を取得
func (db *Sqlite3DB) GetEntriesByRequestID(requestID int) ([]Entry, error) {
	rows, err := db.Conn.Query("SELECT id, request_id, user_id, date, hour FROM entries WHERE request_id = ?", requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		err := rows.Scan(&entry.ID, &entry.RequestID, &entry.UserID, &entry.Date, &entry.Hour)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// 新しいシフトリクエストを作成
func (db *Sqlite3DB) CreateRequest(creatorID int, startDate string, endDate string, deadline string) (int, error) {
	res, err := db.Conn.Exec(
		"INSERT INTO requests (creator_id, start_date, end_date, deadline) VALUES (?, ?, ?, ?)",
		creatorID, startDate, endDate, deadline,
	)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

// 新しい1つのエントリーを作成
func (db *Sqlite3DB) createEntry(requestID int, userID int, date string, hour int) (int, error) {
	res, err := db.Conn.Exec(
		"INSERT INTO entries (request_id, user_id, date, hour) VALUES (?, ?, ?, ?)",
		requestID, userID, date, hour,
	)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

// 新しいエントリーを作成
func (db *Sqlite3DB) CreateEntries(entries []Entry) ([]int, error) {
	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var ids []int
	for _, entry := range entries {
		id, err := db.createEntry(entry.RequestID, entry.UserID, entry.Date, entry.Hour)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return ids, nil
}
