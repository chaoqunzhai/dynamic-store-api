package config

var ExtConfig Extend

// Extend 扩展配置
//
//	extend:
//	  demo:
//	    name: demo-name
//
// 使用方法： config.ExtConfig......即可！！
type Extend struct {
	Work        string `json:"work"`
	AMap        AMap
	Redis       Redis  `json:"redis"`
	FyPayClient FyPayClient   `json:"fyPayClient"`
	WxLeader    WxLeaderLogin `json:"wxLeader"`
	WxUser      WxUserLogin   `json:"wxUser"`
	WxOfficial  WxOfficial    `json:"wxOfficial"`
	CloudObsUrl string  `json:"cloudObsUrl"` //云对象存储
	H5Url string `json:"h5Url"`
	ExportDay int `json:"exportDay"`
	ImageBase string        `json:"imageBase"`
	PromotionCode string `json:"promotionCode"`

	Qiniu Qiniu `json:"qiniu"`
}
type Qiniu struct {
	AccessKey string `json:"AccessKey"`
	SecretKey string `json:"SecretKey"`
}
type Tx struct {
	CosSecretID string `json:"cosSecretID"`
	CosSecretKey string `json:"cosSecretKey"`
}
type Redis struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Password string `json:"password"`
}
type FyPayClient struct {
	Inscd       string `json:"inscd"`
	PayDomain   string `json:"payDomain"`
	PaySearch   string `json:"paySearch"`
	Mchntcd     string `json:"mchntcd"`
	OrderRandom string `json:"orderRandom"`
	OpenId      string `json:"openId"`
}
type WxOfficial struct {
	Appid  string `json:"appid"`
	Secret string `json:"secret"`
}
type WxLeaderLogin struct {
	Appid  string `json:"appid"`
	Secret string `json:"secret"`
}
type WxUserLogin struct {
	Appid  string `json:"appid"`
	Secret string `json:"secret"`
}
type AMap struct {
	Key string
}
