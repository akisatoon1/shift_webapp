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

// ユーザーIDでユーザーを取得
func (db *Sqlite3DB) GetUserByID(id int) (User, error) {
	var user User
	row := db.Conn.QueryRow("SELECT id, login_id, password, name, role, created_at FROM users WHERE id = ?", id)
	var createdAt string
	err := row.Scan(&user.ID, &user.LoginID, &user.Password, &user.Name, &user.Role, &createdAt)
	if err != nil {
		return User{}, err
	}
	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
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
		req.StartDate, _ = time.Parse("2006-01-02 15:04:05", startDate)
		req.EndDate, _ = time.Parse("2006-01-02 15:04:05", endDate)
		req.Deadline, _ = time.Parse("2006-01-02 15:04:05", deadline)
		req.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		requests = append(requests, req)
	}
	return requests, nil
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
		entry.Date, _ = time.Parse("2006-01-02 15:04:05", date)
		entries = append(entries, entry)
	}
	return entries, nil
}
