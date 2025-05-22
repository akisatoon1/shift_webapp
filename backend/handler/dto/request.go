package dto

// httpリクエストDTO

// LoginRequest はログインリクエストの構造体です
type LoginRequest struct {
	LoginID  string `json:"login_id"`
	Password string `json:"password"`
}

// CreateRequestRequest はシフトリクエスト作成リクエストの構造体です
type CreateRequestRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Deadline  string `json:"deadline"`
}

// CreateEntryRequest はエントリー作成リクエストの構造体です
type CreateEntryRequest struct {
	Date string `json:"date"`
	Hour int    `json:"hour"`
}
