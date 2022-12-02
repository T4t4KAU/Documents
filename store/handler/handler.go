package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"store/meta"
	"store/utils"
	"time"
)

const storePath = "./tmp/"

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// HTML 上传页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			_, _ = io.WriteString(w, "internal server error")
			return
		}
		_, _ = io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 接收文件流
		file, head, err := r.FormFile("file")
		if err != nil {
			log.Printf("Failed to get data: %s\n", err.Error())
			_, _ = io.WriteString(w, "")
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: storePath + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			log.Printf("Failed to create file: %s\n", err.Error())
			return
		}
		defer newFile.Close()
		_, err = io.Copy(newFile, file)
		if err != nil {
			log.Printf("Failed to save data into file: %s\n", err.Error())
			return
		}

		_, _ = newFile.Seek(0, 0)
		fileMeta.FileSha1 = utils.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Upload successfully")
}

// GetFileMetaHandler 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	fileHash := r.Form["filehash"][0]
	fileMeta := meta.GetFileMeta(fileHash)
	data, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	fileHash := r.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(fileHash)
	file, err := os.Open(fileMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fileMeta.FileName+"\"")
	_, _ = w.Write(data)
}

func FileDelHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	fileHash := r.Form.Get("fileHash")
	fileMeta := meta.GetFileMeta(fileHash)
	_ = os.Remove(fileMeta.Location)
	meta.RemoveFileMeta(fileHash)
	w.WriteHeader(http.StatusOK)
}

func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	_ = r.ParseForm()
	opType := r.Form.Get("op")
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	fileHash := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")
	curFileMeta := meta.GetFileMeta(fileHash)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
