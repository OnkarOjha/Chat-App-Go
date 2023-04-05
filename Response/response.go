package Response

type Response struct {
	Status  string      `json:"status"`
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}