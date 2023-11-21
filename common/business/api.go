package business
type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Extend interface{} `json:"extend"`
	Data interface{}`json:"data"`

}

