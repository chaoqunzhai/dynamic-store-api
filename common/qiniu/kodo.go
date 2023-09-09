package qiniu

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"go-admin/config"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"go.uber.org/zap"
)


func SizeFile(filePath string) string {
	// 打开要压缩的图片文件
	file, err := os.Open(filePath)
	if err != nil {

		return filePath
	}

	fileDir,_:=path.Split(filePath)
	// 读取图片
	img, format, err := image.Decode(file)
	if err != nil {
		return filePath
	}
	// 关闭文件
	defer func() {
		_=file.Close()
	}()
	// 设置压缩后的宽度和高度，这里是压缩为原图宽度和高度的 1/4

	newWidth := uint(img.Bounds().Dx() )
	newHeight :=  uint(img.Bounds().Dy() )


	if newHeight < 750 {
		return filePath
	}
	resizeHeight:=750

	h:=(uint(resizeHeight) * newWidth) / newHeight

	// 压缩图片

	resizedImg := resize.Resize(h, uint(resizeHeight), img, resize.Lanczos3)

	// 创建输出文件

	//文件名

	uuidName:=strings.Split(uuid.New().String(), "-")[0]

	newFileName:=  fmt.Sprintf("%v.%v",uuidName,format)

	newFilePath	:=path.Join(fileDir,newFileName)


	outFile, err := os.Create(newFilePath)
	if err != nil {
		return filePath
	}
	defer func() {
		_=outFile.Close()
	}()
	// 根据原图格式进行输出
	if format == "jpeg" {
		jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: 80})
	} else if format == "png" {
		png.Encode(outFile, resizedImg)
	}
	fmt.Println("图片压缩成功！！！",newFilePath)
	return newFilePath

}


type QinUi struct {
	CId int `json:"c_id"` //站点ID
	Cfg storage.Config
	BucketName string `json:"bucket_name"`
	BucketManager *storage.BucketManager
	Token string `json:"token"`
}
func (q *QinUi)InitClient()  {

	accessKey := config.ExtConfig.Qiniu.AccessKey
	secretKey := config.ExtConfig.Qiniu.SecretKey
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuadongZheJiang2
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	//
	q.Cfg = cfg

	//统一创建dcy-用户ID的桶
	q.BucketName = fmt.Sprintf("dcy-%v",q.CId)

	putPolicy:=storage.PutPolicy{
		Scope: q.BucketName,
	}

	mac := qbox.NewMac(accessKey, secretKey)

	q.Token = putPolicy.UploadToken(mac)

	q.BucketManager = storage.NewBucketManager(mac, &cfg)

}
//创建存储桶
func (q *QinUi)CreateBucket()  {

}

//上传文件
func (q *QinUi)PostFile(filePath  string) (name string,err  error)  {
	//filePath := "/Users/zhaichaoqun/workespace/goProjects/src/test/70e3f85b.jpg"

	//压缩下文件
	sizeFilePath :=SizeFile(filePath)
	//sizeFilePath:是压缩后的文件
	_,fileName := path.Split(sizeFilePath)


	//对文件进行压缩

	formUploader := storage.NewFormUploader(&q.Cfg)

	ret := storage.PutRet{}

	//绝对路径

	putExtra := storage.PutExtra{}

	//保留全路径 会在七牛云上创建目录
	err = formUploader.PutFile(context.Background(), &ret, q.Token, sizeFilePath, sizeFilePath, &putExtra)
	if err != nil {
		zap.S().Infof("七牛云图片上传失败:%v",err)
		return "", err
	}
	//上传成功后,删除这个压缩的文件
	os.Remove(sizeFilePath)
	return fileName, err
}
func (q *QinUi)MakeUrl(fileName string) string  {


	return ""
}
//删除文件
func (q *QinUi)RemoveFile(FileName string)  {

	err:=q.BucketManager.Delete(q.BucketName,FileName)

	fmt.Printf("删除图片:%v 成功,返回:%v",FileName,err)
}