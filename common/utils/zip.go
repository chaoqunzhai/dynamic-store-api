/**
@Author: chaoqun
* @Date: 2023/12/29 11:47
*/
package utils

import (
	"archive/zip"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
)

func ZipFile(zipName string,files []string) (zipOutFile string,err error) {

	// 创建输出文件
	outputFile, err := os.Create(zipName)
	if err != nil {
		return "", err
	}
	defer func() {
		_=outputFile.Close()
	}()

	// 创建zip写入器
	zipWriter := zip.NewWriter(outputFile)
	defer func() {
		_=zipWriter.Close()
	}()

	// 遍历文件列表，将每个文件添加到zip压缩包中
	for _, file := range files {
		// 打开要压缩的文件
		srcFile, openErr := os.Open(file)
		if openErr != nil {
			zap.S().Errorf("压缩文件时出错, %v 无法打开文件 err:%v",file,openErr)
			_=srcFile.Close()
			continue
		}

		// 获取文件信息，包括文件名和文件大小
		fileInfo, StatErr := srcFile.Stat()
		if StatErr != nil {
			zap.S().Errorf("压缩文件时出错, %v 无法获取文件信息 err:%v",file,StatErr)
			_=srcFile.Close()
			continue
		}

		// 创建zip文件条目，并将文件内容写入zip压缩包中
		zipEntry, zipErr := zipWriter.Create(fileInfo.Name())
		if zipErr != nil {
			zap.S().Errorf("压缩文件时出错, %v 无法创建zip条目 err:%v",file,StatErr)
			_=srcFile.Close()
			continue
		}
		if _,copyErr:=io.Copy(zipEntry, srcFile);copyErr!=nil{
			zap.S().Errorf("压缩文件时出错, %v io.Copy err:%v",file,StatErr)
			_=srcFile.Close()
			continue
		}
		_=srcFile.Close()
	}

	fmt.Println("压缩完成！")
	return zipName,nil
}
