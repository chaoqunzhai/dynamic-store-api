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
	FyPayClient FyPayClient   `json:"fyPayClient"`
	WxLeader    WxLeaderLogin `json:"wxLeader"`
	WxUser      WxUserLogin   `json:"wxUser"`
	WxOfficial  WxOfficial    `json:"wxOfficial"`
	Domain      string        `json:"domain"`
	ImageBase string        `json:"imageBase"`
	CityAdv     string        `json:"cityAdv"`
	Compose     string        `json:"compose"`
	Influx      Influx        `json:"influx"`
	Callback    string        `json:"callback"`
	Harbor      Harbor        `json:"harbor"`
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
type Harbor struct {
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Password string `json:"password"`
	Callback string `json:"callback"`
}
type Influx struct {
	Host     string `json:"host"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Port     int    `json:"port"`
}
type AMap struct {
	Key string
}
