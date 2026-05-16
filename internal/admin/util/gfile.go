package util

import (
	"fmt"
	"mime/multipart"

	"gotribe/internal/core/constant"

	"github.com/h2non/filetype"
)

// File 文件工具类，用于处理文件相关操作
type File struct{}

// GetFileType 获取文件类型
func (f *File) GetFileType(header *multipart.FileHeader) (int, error) {
	file, err := header.Open()
	if err != nil {
		return constant.FILE_TYPE_UNKNOWN, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	head := make([]byte, 261)
	_, err = file.Read(head)
	if err != nil {
		return constant.FILE_TYPE_UNKNOWN, fmt.Errorf("无法读取文件头: %v", err)
	}

	// 检查文件类型
	var fileType int
	switch {
	case filetype.IsImage(head):
		fileType = constant.FILE_TYPE_IMAGE
	case filetype.IsAudio(head):
		fileType = constant.FILE_TYPE_AUDIO
	case filetype.IsApplication(head):
		fileType = constant.FILE_TYPE_APP
	case filetype.IsVideo(head):
		fileType = constant.FILE_TYPE_VIDEO
	case filetype.IsArchive(head):
		fileType = constant.FILE_TYPE_ARCHIVE
	case filetype.IsDocument(head):
		fileType = constant.FILE_TYPE_DOCUMENT
	case filetype.IsFont(head):
		fileType = constant.FILE_TYPE_FONT
	default:
		fileType = constant.FILE_TYPE_UNKNOWN
	}
	return fileType, nil
}

// 全局实例
var FileUtil = &File{}
