package model

/*
	apiの仕様に沿ったレスポンスをJSON文字列で返す
*/

import (
	"backend/context"
	"encoding/json"
)

func GetRequests(ctx *context.AppContext) (string, error) {
	type ResponseRequest struct {
		ID      int `json:"id"`
		Creator struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"creator"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Deadline  string `json:"deadline"`
		CreatedAt string `json:"created_at"`
	}
	var response []ResponseRequest

	requests, err := ctx.DB.GetRequests()
	if err != nil {
		return "", err
	}

	for _, request := range requests {
		user, err := ctx.DB.GetUserByID(request.CreatorID)
		if err != nil {
			return "", err
		}
		response = append(response, ResponseRequest{
			ID: request.ID,
			Creator: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{user.ID, user.Name},
			StartDate: request.StartDate.Format("2006-01-02"),
			EndDate:   request.EndDate.Format("2006-01-02"),
			Deadline:  request.Deadline.Format("2006-01-02"),
			CreatedAt: request.CreatedAt.Format("2006-01-02"),
		})
	}

	json, _ := json.Marshal(response)
	return string(json), nil
}
