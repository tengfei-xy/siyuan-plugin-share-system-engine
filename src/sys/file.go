package sys

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	log "github.com/tengfei-xy/go-log"
)

func Unzip(dest, zipFile string) (err error) {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		file.Name = strings.TrimPrefix(file.Name, "resources/share-note/")
		filename := filepath.Join(dest, file.Name)
		err = os.MkdirAll(filepath.Dir(filename), GetFolderPermission())
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

// 描述: 基于保存目录创建appid和docid文件夹
// 返回: 返回创建的文件夹路径
// 返回: 错误
func MkdirAll(save_path, app_id, doc_id string) (string, error) {
	f := filepath.Join(save_path, app_id, doc_id)
	err := os.MkdirAll(f, GetFolderPermission())
	if err != nil {
		if err != os.ErrExist {
			log.Error(err)
			return f, err
		} else {
			return "", err
		}
	}
	return f, nil
}

// 删除目录
func RmdirAll(save_path, app_id, doc_id string) error {
	f := filepath.Join(save_path, app_id, doc_id)
	err := os.RemoveAll(f)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func GetFilePermission() fs.FileMode {
	return 0660
}
func GetFolderPermission() fs.FileMode {
	return 0755
}
