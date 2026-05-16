package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
)

// FileSize 读取 multipart 文件的字节大小。
func FileSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)
	return len(content), err
}

// FileExt 获取文件扩展名。
func FileExt(fileName string) string {
	return path.Ext(fileName)
}

// FileExist 检查文件是否存在，存在时返回 os.FileInfo。
func FileExist(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	return info, err
}

// IsNotExistMkDir 目录不存在时创建。
func IsNotExistMkDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

// CopyFile 复制文件。
func CopyFile(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer df.Close()

	if _, err := io.Copy(df, sf); err != nil {
		return err
	}
	if info, err := os.Stat(src); err == nil {
		return os.Chmod(dst, info.Mode())
	}
	return nil
}

// CopyDir 递归复制目录。
func CopyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}
