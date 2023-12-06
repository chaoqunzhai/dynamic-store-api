package business

import (
	"fmt"
	"go-admin/config"
	"go-admin/global"
	"net/url"
	"path"
	"strings"
)

func GetGoodsPathFirst(uid interface{},value,imageConst string) string  {

	imagList:=strings.Split(value,",")[0]

	return GetDomainCosEncodePathName(imageConst,uid,imagList,false)
}
//return:大B用户ID/常量定义{goods,ads}/WechatIMG950.jpg
//这个路径是存在服务端的

func GetSiteCosPath(uid interface{},imageConst,image string) string  {

	goodsImagePath := path.Join( fmt.Sprintf("%v", uid),imageConst,image)

	return goodsImagePath
}


//return:image/大B用户ID/goods/WechatIMG950.jpg
func GetGoodPathName(uid interface{}) string {
	goodsImagePath := path.Join(config.ExtConfig.ImageBase,fmt.Sprintf("%v", uid), global.GoodsPath) + "/"

	return goodsImagePath
}
//返回照片的域名
//直接返回的对象存储API地址
func GetDomainCosEncodePathName(imageConst string,uid interface{}, image string,local bool) string {
	//imageType: 是const常量
	if image == ""{
		return ""
	}
	if local{
		//旧逻辑,直接返回本地的文件路径
		goodsImagePath := config.ExtConfig.ImageUrl +  path.Join(config.ExtConfig.ImageBase,
			imageConst, fmt.Sprintf("%v", uid)) + "/" + image
		return goodsImagePath
	}

	//文件路径进行编码
	encodeUrl:=url.QueryEscape(GetSiteCosPath(uid,imageConst,image))
	//cos域名桶的地址,防止后期更换导致的图片查不到
	imageUrl:=config.ExtConfig.ImageUrl
	//返回的图片路径
	//示例地址: "https://dcy-1318497773.cos.ap-nanjing.myqcloud.com/goods%2F1%2FWechatIMG954.jpg"
	domainUrl:=fmt.Sprintf("%v%v",imageUrl,encodeUrl)
	return domainUrl
}
//return: 用户ID/goods/38ab9d1e.jpg
func GetDomainSplitFilePath(url string)  string {
	//https://dcy-1318497773.cos.ap-nanjing.myqcloud.com/goods/1/38ab9d1e.jpg

	c:=strings.Replace(url,config.ExtConfig.ImageUrl,"",-1)
	return c

}