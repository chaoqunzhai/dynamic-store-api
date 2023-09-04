/**
@Author: chaoqun
* @Date: 2023/5/28 14:26
*/
package apis

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/utils"
	"github.com/google/uuid"
	"go-admin/config"
	"io/ioutil"
	"path"
)
type Tools struct {
	api.Api
}


func (e Tools)ShowImage (c *gin.Context) {
	mode := c.Param("t")
	imageName := c.Param("name")
	pathFile := path.Join(config.ExtConfig.ImageBase,mode,imageName)

	file, _ := ioutil.ReadFile(pathFile)
	_, _ = c.Writer.WriteString(string(file))


}

func (e Tools)SaveImage(c *gin.Context)  {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	reqMode := c.Param("mode")
	files, err := c.FormFile("file")

	if err != nil {

		e.Error(-1, errors.New(""), "图片不能为空")
		return
	}
	// 上传文件至指定目录
	guid := uuid.New().String()

	fileName := guid + utils.GetExt(files.Filename)
	pathName := path.Join(config.ExtConfig.ImageBase, reqMode)
	err = utils.IsNotExistMkDir(pathName)
	if err != nil {
		e.Error(-1, nil, "初始化文件路径失败")
		return
	}

	singleFile := pathName + "/" + fileName
	_ = c.SaveUploadedFile(files, singleFile)
	//只返回重命名的文件名字就行
	e.OK(map[string]string{
		"path": fileName,
	}, "successful")
	return
}