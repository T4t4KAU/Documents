package main

import (
	"fmt"
	"net/http"
	"store/handler"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/delete", handler.FileDelHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %s", err.Error())
	}
}
