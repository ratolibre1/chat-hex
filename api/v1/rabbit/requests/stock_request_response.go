package requests

type StockRequestResponse struct {
	StockCode string `json:"stockCode"`
	Chatroom  string `json:"chatroom"`
}