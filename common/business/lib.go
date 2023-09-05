package business

import (
	"fmt"
	"go-admin/config"
	"go-admin/global"
	"net/url"
	"path"
	"strings"
)

func GetGoodsPathFirst(uid interface{},value string) string  {

	imagList:=strings.Split(value,",")[0]

	return GetDomainGoodPathName(uid,imagList,false)
}
//return:goods/站点ID/WechatIMG950.jpg
//这个路径是存在服务端的
func GetSiteGoodsPath(uid interface{},image string) string  {

	goodsImagePath := path.Join(global.GoodsPath, fmt.Sprintf("%v", uid),image)

	return goodsImagePath
}
//return:image/goods/1/WechatIMG950.jpg
func GetGoodPathName(uid interface{}) string {
	goodsImagePath := path.Join(config.ExtConfig.ImageBase, global.GoodsPath,
		fmt.Sprintf("%v", uid)) + "/"

	return goodsImagePath
}
//返回照片的域名
//直接返回腾讯COS的对象存储API地址
func GetDomainGoodPathName(uid interface{}, image string,local bool) string {
	if local{
		//旧逻辑,直接返回本地的文件路径
		goodsImagePath := config.ExtConfig.ImageUrl +  path.Join(config.ExtConfig.ImageBase,
			global.GoodsPath, fmt.Sprintf("%v", uid)) + "/" + image
		return goodsImagePath
	}
	//文件路径进行编码
	encodeUrl:=url.QueryEscape(GetSiteGoodsPath(uid,image))
	//cos域名桶的地址,防止后期更换导致的图片查不到
	imageUrl:=config.ExtConfig.ImageUrl
	//返回的图片路径
	//示例地址: "https://dcy-1318497773.cos.ap-nanjing.myqcloud.com/goods%2F1%2FWechatIMG954.jpg"
	domainUrl:=fmt.Sprintf("%v%v",imageUrl,encodeUrl)
	return domainUrl
}
//return: goods/1/38ab9d1e.jpg
func GetDomainSplitFilePath(url string)  string {
	//https://dcy-1318497773.cos.ap-nanjing.myqcloud.com/goods/1/38ab9d1e.jpg

	c:=strings.Replace(url,config.ExtConfig.ImageUrl,"",-1)
	return c

}