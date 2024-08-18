package model

type ServerResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Object *Users `json:"data"`
	Error  any    `json:"error"`
	Token  string `json:"token"`
}
