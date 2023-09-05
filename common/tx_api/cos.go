/**
@Author: chaoqun
* @Date: 2023/9/4 00:56
*/
package tx_api

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"

	"net/http"
	"net/url"

)

const  (
	SecretID = "AKIDuc92sDKFwPiwymMisoVxCF5h5FBtMbcO"
	SecretKey = "byAQSFpAHApGSQAiVBdrreDZhaFm9LI3"
)

type TxCos struct {
	Client *cos.Client
}

func (t *TxCos)InitClient()  {
	u, _ := url.Parse("https://dcy-1318497773.cos.ap-nanjing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  SecretID,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
			SecretKey: SecretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
		},
	})
	t.Client = c
}
func (t *TxCos)PostFile(name  string) (url string,err  error)  {
	//name:="goods/1/WechatIMG955.jpg"

	UploadResult,_, err := t.Client.Object.Upload(context.Background(), name, name, nil)
	if err != nil {
		fmt.Println("照片上传失败",err)
		return "", err
	}
	fmt.Println("照片上传成功")
	return UploadResult.Location,nil
}

func (t *TxCos)RemoveFile(name string)  {
	_, _ = t.Client.Object.Delete(context.Background(), name)

	fmt.Println("删除成功")

}
//func main()  {
//	name1:="goods/1/WechatIMG950.jpg"
//	name2:="goods/1/WechatIMG954.jpg"
//	name3:="goods/1/WechatIMG955.jpg"
//	t:=TxCos{}
//	t.init()
//	t.PostFile(name1)
//	t.PostFile(name2)
//	t.PostFile(name3)
//	t.RemoveFile(name2)
//
//
//	encodeUrl:=url.QueryEscape(name1)
//
//	fmt.Println("-->",fmt.Sprintf("https://dcy-1318497773.cos.ap-nanjing.myqcloud.com/%v",encodeUrl))
//
//
//}

