package handler

import (
	"backend/auth"
	"backend/context"
	"backend/handler/dto"
	"backend/model"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

/*
	APIエンドポイントに対応するハンドラー関数
*/

var ErrNotLoggedIn = errors.New("user not logged in")

func LoginRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	var loginReq dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		return NewAppError(err, "リクエストボディのデコードに失敗しました", http.StatusBadRequest)
	}

	err := auth.Login(ctx, w, r, loginReq.LoginID, loginReq.Password)
	if err != nil {
		if errors.Is(err, auth.ErrIncorrectAuth) {
			return NewAppError(err, "ログインIDまたはパスワードが間違っています", http.StatusUnauthorized)
		}
		return NewAppError(err, "ログインに失敗しました", http.StatusInternalServerError)
	}
	return nil
}

func GetSessionRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	// ログインユーザのみ認可
	userID, isLoggedIn := auth.GetUserID(ctx, r)
	if !isLoggedIn {
		return NewAppError(ErrNotLoggedIn, "ログインしていません", http.StatusUnauthorized)
	}

	// ユーザー情報を取得
	var usr model.User
	user, err := usr.FindByID(ctx, userID)
	if err != nil {
		return NewAppError(err, "セッションの取得に失敗しました", http.StatusInternalServerError)
	}

	// レスポンスDTOを作成
	roles := []string{}
	// TODO: 抽象化できていない
	if (user.Role & auth.RoleEmployee) != 0 {
		roles = append(roles, "employee")
	}
	if (user.Role & auth.RoleManager) != 0 {
		roles = append(roles, "manager")
	}

	sessionResponse := dto.SessionResponse{
		User: dto.UserSessionInfo{
			ID:        user.ID,
			Name:      user.Name,
			Roles:     roles,
			CreatedAt: user.CreatedAt.Format(),
		},
	}

	json.NewEncoder(w).Encode(sessionResponse)
	return nil
}

func LogoutRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	err := auth.Logout(ctx, w, r)
	if err != nil {
		return NewAppError(err, "ログアウトに失敗しました", http.StatusInternalServerError)
	}
	return nil
}

func GetRequestsRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	// ログインユーザのみ認可
	if _, isLoggedIn := auth.GetUserID(ctx, r); !isLoggedIn {
		return NewAppError(ErrNotLoggedIn, "ログインしていません", http.StatusUnauthorized)
	}

	var req model.Request
	requests, err := req.FindAll(ctx)
	if err != nil {
		return NewAppError(err, "シフトリクエストの取得に失敗しました", http.StatusInternalServerError)
	}

	// モデルをDTOに変換
	var requestsResponse dto.RequestsResponse
	for _, req := range requests {
		requestInfo := dto.RequestInfo{
			ID: req.ID,
			Creator: dto.UserInfo{
				ID:   req.Creator.ID,
				Name: req.Creator.Name,
			},
			StartDate: req.StartDate.Format(),
			EndDate:   req.EndDate.Format(),
			Deadline:  req.Deadline.Format(),
			CreatedAt: req.CreatedAt.Format(),
		}
		requestsResponse = append(requestsResponse, requestInfo)
	}

	json.NewEncoder(w).Encode(requestsResponse)
	return nil
}

func GetRequestRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	// ログインユーザのみ認可
	if _, isLoggedIn := auth.GetUserID(ctx, r); !isLoggedIn {
		return NewAppError(ErrNotLoggedIn, "ログインしていません", http.StatusUnauthorized)
	}

	requestId := r.PathValue("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return NewAppError(err, "requestIdが整数ではありません", http.StatusBadRequest)
	}

	// リクエスト情報を取得
	var req model.Request
	request, err := req.FindByID(ctx, requestIdInt)
	if err != nil {
		// TODO: errcode修正
		return NewAppError(err, "リクエストの取得に失敗しました", http.StatusInternalServerError)
	}

	// 提出情報を取得
	var sub model.Submission
	submissions, err := sub.FindByRequestID(ctx, requestIdInt)
	if err != nil {
		return NewAppError(err, "提出情報の取得に失敗しました", http.StatusInternalServerError)
	}

	// シフト提出情報をDTOに変換
	var submissionsInfo []dto.SubmissionInfo
	var entriesInfo []dto.EntryInfo

	for _, submission := range submissions {
		// シフト提出情報のみを処理
		submissionInfo := dto.SubmissionInfo{
			Submitter: dto.UserInfo{
				ID:   submission.Submitter.ID,
				Name: submission.Submitter.Name,
			},
		}
		submissionsInfo = append(submissionsInfo, submissionInfo)

		// エントリー情報を処理
		for _, entry := range submission.Entries {
			entryInfo := dto.EntryInfo{
				ID: entry.ID,
				User: dto.UserInfo{
					ID:   submission.Submitter.ID,
					Name: submission.Submitter.Name,
				},
				Date: entry.Date.Format(),
				Hour: entry.Hour,
			}
			entriesInfo = append(entriesInfo, entryInfo)
		}

	}

	// レスポンスDTOを作成
	response := dto.RequestDetailResponse{
		ID: request.ID,
		Creator: dto.UserInfo{
			ID:   request.Creator.ID,
			Name: request.Creator.Name,
		},
		StartDate:   request.StartDate.Format(),
		EndDate:     request.EndDate.Format(),
		Deadline:    request.Deadline.Format(),
		CreatedAt:   request.CreatedAt.Format(),
		Submissions: submissionsInfo,
		Entries:     entriesInfo,
	}

	json.NewEncoder(w).Encode(response)
	return nil
}

func PostRequestsRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	// ログインしているユーザーのIDを取得する
	userID, isLoggedIn := auth.GetUserID(ctx, r)
	if !isLoggedIn {
		return NewAppError(ErrNotLoggedIn, "ログインしていません", http.StatusUnauthorized)
	}

	// リクエストボディのデコード
	var createReq dto.CreateRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		return NewAppError(err, "リクエストボディのデコードに失敗しました", http.StatusBadRequest)
	}

	// DTOからモデルに変換
	// 文字列の日付をモデルの型に変換
	startDate, err := model.NewDateOnly(createReq.StartDate)
	if err != nil {
		return NewAppError(err, "開始日のフォーマットが不正です", http.StatusBadRequest)
	}

	endDate, err := model.NewDateOnly(createReq.EndDate)
	if err != nil {
		return NewAppError(err, "終了日のフォーマットが不正です", http.StatusBadRequest)
	}

	deadline, err := model.NewDateTime(createReq.Deadline)
	if err != nil {
		return NewAppError(err, "期限日のフォーマットが不正です", http.StatusBadRequest)
	}

	// 新しいシフトリクエストを作成する
	var req model.Request
	requestID, err := req.Create(ctx, model.NewRequest{
		CreatorID: userID,
		StartDate: startDate,
		EndDate:   endDate,
		Deadline:  deadline,
	})
	if err != nil {
		if errors.Is(err, model.ErrForbidden) {
			return NewAppError(err, "権限がありません", http.StatusForbidden)
		}

		var inputErr model.InputError
		if errors.As(err, &inputErr) {
			return NewAppError(inputErr, inputErr.Message(), http.StatusBadRequest)
		}

		return NewAppError(err, "シフトリクエストの作成に失敗しました", http.StatusInternalServerError)
	}

	// レスポンスDTOを作成
	response := dto.CreateRequestResponse{
		ID: requestID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	return nil
}

func PostSubmissionsRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	// ログインしているユーザーのIDを取得する
	userID, isLoggedIn := auth.GetUserID(ctx, r)
	if !isLoggedIn {
		return NewAppError(ErrNotLoggedIn, "ログインしていません", http.StatusUnauthorized)
	}

	// シフトリクエストのIDを取得する
	// 整数ではない場合はエラーを返す
	requestId := r.PathValue("id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return NewAppError(err, "requestidが整数ではありません", http.StatusBadRequest)
	}

	// リクエストボディのデコード
	var entryRequests []dto.CreateEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&entryRequests); err != nil {
		return NewAppError(err, "リクエストボディのデコードに失敗しました", http.StatusBadRequest)
	}

	// モデルに渡す形に変換する
	newSubmission := model.NewSubmission{
		RequestID:   requestIdInt,
		SubmitterID: userID,
		NewEntries:  []model.NewEntry{},
	}

	// DTOからモデルに変換
	for _, entry := range entryRequests {
		// 日付文字列をモデルの型に変換
		dateOnly, err := model.NewDateOnly(entry.Date)
		if err != nil {
			return NewAppError(err, "日付のフォーマットが不正です", http.StatusBadRequest)
		}

		newSubmission.NewEntries = append(newSubmission.NewEntries, model.NewEntry{
			Date: dateOnly,
			Hour: entry.Hour,
		})
	}

	// 新しい提出を作成
	var sub model.Submission
	submissionID, err := sub.Create(ctx, newSubmission)
	if err != nil {
		if errors.Is(err, model.ErrForbidden) {
			return NewAppError(err, "権限がありません", http.StatusForbidden)
		}

		var inputErr model.InputError
		if errors.As(err, &inputErr) {
			return NewAppError(inputErr, inputErr.Message(), http.StatusBadRequest)
		}

		return NewAppError(err, "エントリーの作成に失敗しました", http.StatusInternalServerError)
	}

	// レスポンスDTOを作成
	response := struct {
		ID int `json:"id"`
	}{submissionID}

	// レスポンスを返す
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	return nil
}

func GetMySubmissionRequest(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) *AppError {
	// ログインしているユーザーのIDを取得する
	userID, isLoggedIn := auth.GetUserID(ctx, r)
	if !isLoggedIn {
		return NewAppError(ErrNotLoggedIn, "ログインしていません", http.StatusUnauthorized)
	}

	// シフトリクエストのIDを取得する
	requestId := r.PathValue("request_id")
	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		return NewAppError(err, "request_idが整数ではありません", http.StatusBadRequest)
	}

	// 自分の提出を取得
	var sub model.Submission
	submission, err := sub.FindByRequestIDAndSubmitterID(ctx, requestIdInt, userID)
	if err != nil {
		return NewAppError(err, "提出情報の取得に失敗しました", http.StatusInternalServerError)
	}

	// 提出がない場合は空のレスポンスを返す
	if submission == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"submission": nil,
		})
		return nil
	}

	// エントリー情報をDTOに変換
	type EntryDTO struct {
		ID   int    `json:"id"`
		Date string `json:"date"`
		Hour int    `json:"hour"`
	}

	var entriesInfo []EntryDTO
	for _, entry := range submission.Entries {
		entryInfo := EntryDTO{
			ID:   entry.ID,
			Date: entry.Date.Format(),
			Hour: entry.Hour,
		}
		entriesInfo = append(entriesInfo, entryInfo)
	}

	// レスポンスDTOを作成
	type SubmissionDTO struct {
		ID      int          `json:"id"`
		User    dto.UserInfo `json:"user"`
		Entries []EntryDTO   `json:"entries"`
	}

	response := struct {
		Submission *SubmissionDTO `json:"submission"`
	}{
		Submission: &SubmissionDTO{
			ID: submission.ID,
			User: dto.UserInfo{
				ID:   submission.Submitter.ID,
				Name: submission.Submitter.Name,
			},
			Entries: entriesInfo,
		},
	}

	json.NewEncoder(w).Encode(response)
	return nil
}
