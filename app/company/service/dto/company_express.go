package dto

type CompanyExpressCnfReq struct {
	Store struct {
		Enable  bool      `json:"enable"`
		Address []Address `json:"address"`
	} `json:"store"`
	Local struct {
		Enable       bool `json:"enable"`
		StartMoney   int  `json:"start_money"`
		QuotaMoney   int  `json:"quota_money"`
		FreightMoney int  `json:"freight_money"`
	} `json:"local"`
}
type Address struct {
	Id int `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Start   string `json:"start"`
	End     string `json:"end"`
}
