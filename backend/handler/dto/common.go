package dto

// UserSessionInfo はセッション内のユーザー情報の構造体です
type UserSessionInfo struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Roles     []string `json:"roles"`
	CreatedAt string   `json:"created_at"`
}

// UserInfo はユーザー情報の構造体です
type UserInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// RequestInfo はリクエスト一覧内の個別リクエストの構造体です
type RequestInfo struct {
	ID        int      `json:"id"`
	Creator   UserInfo `json:"creator"`
	StartDate string   `json:"start_date"`
	EndDate   string   `json:"end_date"`
	Deadline  string   `json:"deadline"`
	CreatedAt string   `json:"created_at"`
}

// EntryInfo はエントリー情報の構造体です
type EntryInfo struct {
	ID   int      `json:"id"`
	User UserInfo `json:"user"`
	Date string   `json:"date"`
	Hour int      `json:"hour"`
}

// EntryIDInfo はエントリーID情報の構造体です
type EntryIDInfo struct {
	ID int `json:"id"`
}
