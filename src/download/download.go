package download

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func Download(filePath string, w http.ResponseWriter) {
	fmt.Println("download", filePath)

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("文件打开失败", err)
		msg := fmt.Sprintf("文件打开失败，错误：%v", err.Error())
		w.Write([]byte(msg))
		return
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("读取文件失败", err)
		msg := fmt.Sprintf("读取文件失败，错误：%v", err.Error())
		w.Write([]byte(msg))
		return
	}

	len, err := w.Write(data)
	// len, err := io.Copy(w, f)
	if err != nil {
		fmt.Println("下载文件失败", err)
		return
	}

	fmt.Println("下载文件成功", len)
}

func GetFileList(dirName string) ([]string, error) {
	dirpath, err := filepath.Abs(dirName)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return nil, err
	}

	pathList := make([]string, 0, 0)
	for _, f := range files {
		pathList = append(pathList, f.Name())
	}

	return pathList, nil
}
