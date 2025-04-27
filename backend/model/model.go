package model

/*
	apiの仕様に沿ったレスポンスをJSON文字列で返す
*/

import (
	"backend/context"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Request struct {
	ID        int    `json:"id"`
	Creator   User   `json:"creator"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Deadline  string `json:"deadline"`
	CreatedAt string `json:"created_at"`
}

type GetRequestsResponse []Request

func GetRequests(ctx *context.AppContext) (GetRequestsResponse, error) {
	var response GetRequestsResponse

	requests, err := ctx.DB.GetRequests()
	if err != nil {
		return nil, err
	}

	for _, request := range requests {
		user, err := ctx.DB.GetUserByID(request.CreatorID)
		if err != nil {
			return nil, err
		}
		response = append(response, struct {
			ID        int    `json:"id"`
			Creator   User   `json:"creator"`
			StartDate string `json:"start_date"`
			EndDate   string `json:"end_date"`
			Deadline  string `json:"deadline"`
			CreatedAt string `json:"created_at"`
		}{
			ID:        request.ID,
			Creator:   User{user.ID, user.Name},
			StartDate: request.StartDate.Format("2006-01-02"),
			EndDate:   request.EndDate.Format("2006-01-02"),
			Deadline:  request.Deadline.Format("2006-01-02"),
			CreatedAt: request.CreatedAt.Format("2006-01-02"),
		})
	}

	return response, nil
}

type Entry struct {
	ID   int    `json:"id"`
	User User   `json:"user"`
	Date string `json:"date"`
	Hour int    `json:"hour"`
}

type GetEntriesResponse struct {
	ID      int     `json:"id"`
	Entries []Entry `json:"entries"`
}

func GetEntries(ctx *context.AppContext, requestID int) (GetEntriesResponse, error) {
	response := GetEntriesResponse{
		ID: requestID,
	}

	entries, err := ctx.DB.GetEntriesByRequestID(requestID)
	if err != nil {
		return GetEntriesResponse{}, err
	}

	for _, entry := range entries {
		user, err := ctx.DB.GetUserByID(entry.UserID)
		if err != nil {
			return GetEntriesResponse{}, err
		}
		response.Entries = append(response.Entries, Entry{
			ID:   entry.ID,
			User: User{user.ID, user.Name},
			Date: entry.Date.Format("2006-01-02"),
			Hour: entry.Hour,
		})
	}

	return response, nil
}
