package meta

// FileMeta 文件元信息结构
type FileMeta struct {
	FileSha1 string // 文件哈希
	FileName string // 文件名
	FileSize int64  // 文件大小
	Location string // 文件位置
	UploadAt string // 上传时间
}

// 保存文件元信息映射
var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta) // 初始化
}

// UpdateFileMeta 更新文件元信息
func UpdateFileMeta(fileMeta FileMeta) {
	fileMetas[fileMeta.FileSha1] = fileMeta
}

// GetFileMeta 获取文件元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// RemoveFileMeta 删除文件元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
