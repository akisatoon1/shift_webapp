package dto

// httpレスポンスDTO

// SessionResponse はセッション情報のレスポンス構造体です
type SessionResponse struct {
	User UserSessionInfo `json:"user"`
}

// RequestsResponse はリクエスト一覧のレスポンス構造体です
type RequestsResponse []RequestInfo

// CreateRequestResponse はリクエスト作成レスポンスの構造体です
type CreateRequestResponse struct {
	ID int `json:"id"`
}

// RequestDetailResponse はリクエスト詳細レスポンスの構造体です
type RequestDetailResponse struct {
	ID        int         `json:"id"`
	Creator   UserInfo    `json:"creator"`
	StartDate string      `json:"start_date"`
	EndDate   string      `json:"end_date"`
	Deadline  string      `json:"deadline"`
	CreatedAt string      `json:"created_at"`
	Entries   []EntryInfo `json:"entries"`
}

// CreateEntriesResponse はエントリー作成レスポンスの構造体です
type CreateEntriesResponse struct {
	ID      int           `json:"id"`
	Entries []EntryIDInfo `json:"entries"`
}
