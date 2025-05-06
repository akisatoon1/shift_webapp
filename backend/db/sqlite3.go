package db

import (
	"database/sql"
	"time"

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

// sqlite3では時間は文字列型で保存されるため、
// 共通フォーマットを用いてtime.Time型に変換する
func parseTime(t string) time.Time {
	parsed, _ := time.Parse(time.DateTime, t)
	return parsed
}

// sqlite3の保存する文字列型のtimeデータのフォーマット
func formatTime[T DateOnly | DateTime](t T) string {
	return time.Time(t).Format(time.DateTime)
}

// ユーザーIDでユーザーを取得
func (db *Sqlite3DB) GetUserByID(id int) (User, error) {
	var user User
	row := db.Conn.QueryRow("SELECT id, login_id, password, name, role, created_at FROM users WHERE id = ?", id)
	var createdAt string
	err := row.Scan(&user.ID, &user.LoginID, &user.Password, &user.Name, &user.Role, &createdAt)
	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}
	user.CreatedAt = DateTime(parseTime(createdAt))
	return user, nil
}

// login_idでユーザーを取得
func (db *Sqlite3DB) GetUserByLoginID(loginID string) (User, error) {
	var user User
	row := db.Conn.QueryRow("SELECT id, login_id, password, name, role, created_at FROM users WHERE login_id = ?", loginID)
	var createdAt string
	err := row.Scan(&user.ID, &user.LoginID, &user.Password, &user.Name, &user.Role, &createdAt)
	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}
	user.CreatedAt = DateTime(parseTime(createdAt))
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
		var startDate, endDate, deadline, createdAt string
		err := rows.Scan(&req.ID, &req.CreatorID, &startDate, &endDate, &deadline, &createdAt)
		if err != nil {
			return nil, err
		}
		req.StartDate = DateOnly(parseTime(startDate))
		req.EndDate = DateOnly(parseTime(endDate))
		req.Deadline = DateTime(parseTime(deadline))
		req.CreatedAt = DateTime(parseTime(createdAt))
		requests = append(requests, req)
	}
	return requests, nil
}

// 指定リクエストIDのリクエストを取得
func (db *Sqlite3DB) GetRequestByID(id int) (Request, error) {
	var req Request
	row := db.Conn.QueryRow("SELECT id, creator_id, start_date, end_date, deadline, created_at FROM requests WHERE id = ?", id)
	var startDate, endDate, deadline, createdAt string
	err := row.Scan(&req.ID, &req.CreatorID, &startDate, &endDate, &deadline, &createdAt)
	if err != nil {
		return Request{}, err
	}
	req.StartDate = DateOnly(parseTime(startDate))
	req.EndDate = DateOnly(parseTime(endDate))
	req.Deadline = DateTime(parseTime(deadline))
	req.CreatedAt = DateTime(parseTime(createdAt))
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
		var date string
		err := rows.Scan(&entry.ID, &entry.RequestID, &entry.UserID, &date, &entry.Hour)
		if err != nil {
			return nil, err
		}
		entry.Date = DateOnly(parseTime(date))
		entries = append(entries, entry)
	}
	return entries, nil
}

// 新しいシフトリクエストを作成
func (db *Sqlite3DB) CreateRequest(creatorID int, startDate DateOnly, endDate DateOnly, deadline DateTime) (int, error) {
	res, err := db.Conn.Exec(
		"INSERT INTO requests (creator_id, start_date, end_date, deadline) VALUES (?, ?, ?, ?)",
		creatorID, formatTime(startDate), formatTime(endDate), formatTime(deadline),
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
func (db *Sqlite3DB) createEntry(requestID int, userID int, date DateOnly, hour int) (int, error) {
	res, err := db.Conn.Exec(
		"INSERT INTO entries (request_id, user_id, date, hour) VALUES (?, ?, ?, ?)",
		requestID, userID, formatTime(date), hour,
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
