package business

import (
	"fmt"
	"go-admin/config"
	"go-admin/global"
	"path"
)

func GetGoodPathName(uid interface{}) string {
	goodsImagePath := path.Join(config.ExtConfig.ImageBase, global.GoodsPath,
		fmt.Sprintf("%v", uid)) + "/"

	return goodsImagePath
}
