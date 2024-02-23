package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const (
	SMSURL = "http://smsapi.weiqucloud.com/sms/httpSmsInterface2"
)

type SmsPostBody struct {
	UserId   string `json:"userId"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Mobile   string `json:"mobile"`
	Content  string `json:"content"`
	SendTime string `json:"sendTime"`
	Action   string `json:"action"`
	Custom   string `json:"custom"`
}
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Status        string `json:"status"`
		Message       string `json:"message"`
		RemainPoint   string `json:"remainPoint"`
		TaskID        string `json:"taskID"`
		SuccessCounts string `json:"successCounts"`
	} `json:"data"`
}

var (
	SmsClient http.Client
)

func SendAliEms(Mobile, code string) (status int) {
	status = -1
	SmsClient = http.Client{
		Timeout: 10 * time.Second,
	}
	content := fmt.Sprintf("【动创云】您的验证码为：%v（10分钟内有效），为了保证您的帐户安全，请勿向任何人提供此验证码。", code)
	smsBody := SmsPostBody{
		UserId:   "dongchuangyun",
		Account:  "dongchuangyun",
		Password: "dongchuangyun",
		Action:   "sendhy",
		Mobile:   Mobile,
		Content:  content,
	}

	jsonUser, err := json.Marshal(smsBody)
	if err != nil {
		zap.S().Errorf("发送短信接口:%v,JSON.Marshal失败 原因:%v 手机号:%v", SMSURL, err.Error(), Mobile)
		return
	}
	req, err := http.NewRequest("POST", SMSURL, bytes.NewBuffer(jsonUser))
	if err != nil {
		zap.S().Errorf("发送短信接口:%v,NewRequest失败 原因:%v 手机号:%v", SMSURL, err.Error(), Mobile)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := SmsClient.Do(req)
	if err != nil {

		zap.S().Errorf("发送短信接口:%v,失败 原因:%v 手机号:%v", SMSURL, err.Error(), Mobile)
		return
	}
	defer res.Body.Close()

	responseBytes, _ := ioutil.ReadAll(res.Body)

	ResponseObj := Response{}
	unmarkErr := json.Unmarshal(responseBytes, &ResponseObj)
	if unmarkErr != nil {
		zap.S().Errorf("反序列短信发送返回数据失败 err:%v,接口返回数据:%v 手机号:%v", unmarkErr.Error(), string(responseBytes))
		return
	}

	if ResponseObj.Data.Status == "Success" {
		zap.S().Infof("发送短信接口:%v,成功！返回:%v 手机号:%v", SMSURL, ResponseObj, Mobile)
	} else {
		zap.S().Errorf("发送短信接口:%v,失败 返回:%v 手机号:%v", SMSURL, ResponseObj, Mobile)
	}
	//fmt.Printf("response:%v",ResponseObj)
	return 1
}

// 发生短信验证码时,需要查询大B是否有足够的短信条数
func SendSms(source, phone string, cid int, orm *gorm.DB) (code string, err error) {

	var emsQuotaCnf models2.CompanySmsQuotaCnf
	Available := global.CompanySmsNumber
	RecordTag := global.CompanySmsRecordTag
	orm.Model(&emsQuotaCnf).Select("available,id,record").Where("c_id = ?", cid).Limit(1).Find(&emsQuotaCnf)
	if emsQuotaCnf.Id > 0 {
		Available = emsQuotaCnf.Available
		if Available <= 0 {
			return "", errors.New("短信条数使用完毕")
		}
		RecordTag = emsQuotaCnf.Record
	} else {
		//
		emsQuotaCnf = models2.CompanySmsQuotaCnf{
			Available: Available - 1,
		}
		emsQuotaCnf.CId = cid
		emsQuotaCnf.Layer = 0
		emsQuotaCnf.Enable = true
		emsQuotaCnf.Record = RecordTag
		orm.Create(&emsQuotaCnf)
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	vCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))

	go func(vCode string) {
		status := SendAliEms(phone, vCode)
		//发送成功.减一条
		if status == 1 {
			orm.Model(&emsQuotaCnf).Where("c_id = ?", cid).Updates(map[string]interface{}{
				"available": Available - 1,
			})
			//如果开启了短信记录就记录在DB中
			if RecordTag {
				orm.Create(&models2.CompanySmsRecordLog{
					CId:    cid,
					Phone:  phone,
					Source: source,
					Code:   vCode,
				})
			}

		}
		//是否大B开启短信记录功能,
		//如果开启短信记录,就写入DB
	}(vCode)
	return vCode, nil
}
