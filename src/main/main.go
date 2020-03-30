package main

// macos 编译 GOOS=darwin go build -o ../../publish/darwin/upload_download
// windows 编译 GOOS=windows go build -o ../../publish/windows/upload_download.exe
//
import (
	"download"
	"fmt"
	"html/template"
	"io"
	"log"
	"path/filepath"
	"strings"

	"net/http"
	"os"
	"utils"
)

func main() {
	fmt.Println("系统开始")
	httpServer()
	fmt.Println("系统结束")
}

func httpServer() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)

	// 启动静态文件服务
	downloadDir := filepath.Join(utils.GetCurrentDirectory(), "res", "download")
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(downloadDir))))

	log.Println("http 启动 - ", utils.GetIp(), ":10000")
	if err := http.ListenAndServe(":10000", nil); err != nil {
		log.Fatal("错误结束", "\n", err)
	}
}

type UrlModel struct {
	Url  string
	Name string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	fmt.Println("name = ", name)

	if name == "" {
		curDir := utils.GetCurrentDirectory()
		downloadDir := filepath.Join(curDir, "res", "download")
		os.MkdirAll(downloadDir, os.ModePerm)
		names, err := download.GetFileList(downloadDir)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("文件夹打开失败 - %v", err)))
			return
		}

		urls := make([]UrlModel, 0, 0)
		for _, name := range names {
			if strings.Index(name, ".") == 0 {
				continue
			}
			urls = append(urls, UrlModel{
				Url:  "download/" + name,
				Name: name,
			})
		}

		if len(urls) == 0 {
			io.WriteString(w, "空")
		} else {
			htmlpath := filepath.Join(curDir, "gt", "dirlist.gtml")
			t, err := template.ParseFiles(htmlpath)
			if err != nil {
				io.WriteString(w, fmt.Sprintf("网页错误 = %v", err))
			} else {
				fmt.Println(t.Execute(w, urls))
			}
		}
	} else {
		name = filepath.Join("res", "download", name)
		name, _ = filepath.Abs(name)
		fmt.Println("file path = ", name)
		download.Download(name, w)
	}

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("uploadHandler")
	if r.Method == "GET" {

		htmlpath := utils.GetCurrentDirectory()
		htmlpath = filepath.Join(htmlpath, "gt", "upload.gtml")
		t, err := template.ParseFiles(htmlpath)
		if err != nil {
			io.WriteString(w, fmt.Sprintf("网页错误 = %v", err))
		} else {
			fmt.Println(t.Execute(w, nil))
		}
	} else {
		r.ParseForm()

		fmt.Println("form = ", r.Form)

		r.ParseMultipartForm(128)

		fmt.Println("multipartForm = ", r.MultipartForm)

		f, header, err := r.FormFile("file")
		if err != nil {
			fmt.Println("上传错误:", err)
			io.WriteString(w, fmt.Sprintf("%v", err))
			return
		}

		cur := utils.GetCurrentDirectory()
		updir := filepath.Join(cur, "res", "upload")
		os.MkdirAll(updir, os.ModePerm)

		fileName := header.Filename
		fmt.Println("filename = ", fileName)
		fileName = filepath.Join(updir, fileName)
		fmt.Println("filepath = ", fileName)

		if err := utils.SaveFile(fileName, f); err != nil {
			fmt.Println("上传失败：", err)
			w.Write([]byte(fmt.Sprintf("上传失败：%v", err)))
		} else {
			fmt.Println("上传成功！")
			w.Write([]byte("上传成功！"))
		}

	}

}
