package qiniu

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"go-admin/config"
	"go-admin/global"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"go.uber.org/zap"
)


func CompressImg(source string, hight uint) (newFilePath string,err error) {
	fileDir,_:=path.Split(source)


	var file *os.File
	//reg, _ := regexp.Compile(`^.*\.((png)|(jpg))$`)
	//if !reg.MatchString(source) {
	//	err = errors.New("%s is not a .png or .jpg file")
	//
	//	return "", err
	//}
	if file, err = os.Open(source); err != nil {

		return "", err
	}
	defer func() {
		_=file.Close()
	}()
	name := file.Name()
	var img image.Image
	switch {
	case strings.HasSuffix(name, ".png"):
		if img, err = png.Decode(file); err != nil {

			return "", err
		}
	case strings.HasSuffix(name, ".jpg"):
		if img, err = jpeg.Decode(file); err != nil {

			return "", err
		}
	default:
		return "",errors.New("不支持的类型")
	}
	resizeImg := resize.Resize(hight, 0, img, resize.Lanczos3)

	newFileName := newName(source, int(hight))



	newFilePath	=path.Join(fileDir,newFileName)

	if outFile, createErr := os.Create(newFilePath); createErr != nil {
		return "", createErr
	} else {
		defer outFile.Close()
		err = jpeg.Encode(outFile, resizeImg, nil)
		if err != nil {

			return "", err
		}
	}

	return newFilePath,nil
}
func newName(name string, size int) string {
	_, file := filepath.Split(name)
	return fmt.Sprintf("%d%s", size, file)
}
func SizeFile(filePath string,resizeHeight int) string {
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
	//if newHeight < 400 {
	//	return filePath
	//}

	h:=(uint(resizeHeight) * newWidth) / newHeight

	// 压缩图片
	resizedImg := resize.Resize(h, uint(resizeHeight), img, resize.Lanczos3)
	// 创建输出文件
	//文件名
	uuidName:=strings.Split(uuid.New().String(), "-")[0]

	var newFileName string
	if resizeHeight > 400 {
		newFileName =  fmt.Sprintf("%v.%v",uuidName,format)
	}else {
		newFileName =  fmt.Sprintf("%v.%v",uuidName,format)
	}

	newFilePath	:=path.Join(fileDir,newFileName)


	outFile, err := os.Create(newFilePath)
	if err != nil {
		return filePath
	}
	defer func() {
		_=outFile.Close()
	}()
	// 根据原图格式进行输出
	err = jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: 80})
	return newFilePath

}


type QinUi struct {
	CId interface{} `json:"c_id"` //站点ID
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
	if config.ExtConfig.Qiniu.Region == "ZoneHuadong"{
		cfg.Region = &storage.ZoneHuadong
	}else {
		cfg.Region = &storage.ZoneHuadongZheJiang2
	}
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	//
	q.Cfg = cfg

	//存储图片的桶
	//q.BucketName = "dcy-goods"
	q.BucketName = config.ExtConfig.Qiniu.BucketName
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

func (q *QinUi)ClearCacheImageName(filepath string) string  {

	if strings.HasPrefix(filepath,global.CacheImage){


		fileList :=strings.Split(filepath,global.CacheImage)
		if len(fileList) >0{
			return fileList[1]
		}
		return filepath
	}
	return filepath



}
//上传文件
//不同的大B做不同的桶,
//不设置过期时间
func (q *QinUi)PostImageFile(filePath  string) (name string,err  error)  {
	//filePath := "/Users/zhaichaoqun/workespace/goProjects/src/test/70e3f85b.jpg"

	//压缩下文件
	//细致压缩失败,那就用第二种
	sizeFilePath :=SizeFile(filePath,700)
	//fmt.Println("sizeFilePath",sizeFilePath)
	//sizeFilePath:是压缩后的文件
	_,fileName := path.Split(sizeFilePath)


	//对文件进行压缩

	formUploader := storage.NewFormUploader(&q.Cfg)

	ret := storage.PutRet{}

	//绝对路径

	putExtra := storage.PutExtra{}

	targetPath :=q.ClearCacheImageName(sizeFilePath)

	//保留全路径 会在七牛云上创建目录
	//fmt.Printf("本地路径: %v 对象存储路径:%v\n",sizeFilePath,targetPath)

	err = formUploader.PutFile(context.Background(), &ret, q.Token, targetPath, sizeFilePath, &putExtra)
	if err != nil {
		zap.S().Errorf("BackName:%v 七牛云图片上传文件：%v 失败:%v",q.BucketName,sizeFilePath,err,)
		return "", errors.New(fmt.Sprintf("图片上传失败:%v",err))
	}
	//上传成功后,删除这个压缩的文件
	defer func() {
		_=os.RemoveAll(sizeFilePath)
	}()
	return fileName, err
}

//文件都设置了过期时间

func (q *QinUi)PostFile(filePath  string) (name string,err  error)  {

	//对文件进行压缩

	formUploader := storage.NewFormUploader(&q.Cfg)
	_,fileName := path.Split(filePath)

	ret := storage.PutRet{}

	//绝对路径

	putExtra := storage.PutExtra{}
	//创建全路径 会在七牛云上创建目录
	pathValue :=path.Join(fmt.Sprintf("%v",q.CId),global.CloudExportOrderFilePath,filePath)
	err = formUploader.PutFile(context.Background(), &ret, q.Token, pathValue, filePath, &putExtra)
	if err != nil {
		zap.S().Infof("BackName:%v 七牛云上传文件：%v 失败:%v",q.BucketName,filePath,err,)
		return "", errors.New(fmt.Sprintf("上传失败:%v",err))
	}

	_=q.BucketManager.DeleteAfterDays(q.BucketName,filePath,config.ExtConfig.ExportDay + 1)
	return fileName, err
}
func (q *QinUi)MakeUrl(fileName string) string  {


	return ""
}
//删除文件
func (q *QinUi)RemoveFile(FileName string)  {

	err:=q.BucketManager.Delete(q.BucketName,FileName)

	zap.S().Infof("删除文件 %v 成功,返回:%v\n",FileName,err)
}