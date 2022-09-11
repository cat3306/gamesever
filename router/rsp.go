package router

import "encoding/json"

const (
	ResponseOK  = 200
	ResponseErr = -1
)

type GameJsonResponse struct {
	Code        int    `json:"code"`
	Description string `json:"Description"`
	Data        string `json:"Data"`
}

func JsonRspOK(data interface{}) *GameJsonResponse {
	var raw string
	switch data.(type) {
	default:
		b, _ := json.Marshal(data)
		raw = string(b)
	case string:
		raw = data.(string)
	}
	return &GameJsonResponse{
		Code:        ResponseOK,
		Data:        raw,
		Description: "",
	}
}
func JsonRspErr(desc string) *GameJsonResponse {
	return &GameJsonResponse{
		Code:        ResponseErr,
		Data:        "",
		Description: desc,
	}
}
