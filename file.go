package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	log "github.com/tengfei-xy/go-log"
	"github.com/tengfei-xy/go-tools"
)

func unzip(dest, zipFile string) (err error) {
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
		err = os.MkdirAll(filepath.Dir(filename), get_folder_permission())
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

func check_theme_file(dir string) error {
	update_theme_file(dir, "Odyssey")
	update_theme_file(dir, "Savor")

	return nil
}
func update_theme_file(dir string, theme string) error {
	filename := filepath.Join(dir, fmt.Sprintf("appearance/themes/%s/theme.css", theme))
	if !tools.FileExist(filename) {
		return nil
	}

	log.Infof("修改主题 %s", theme)

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Error(err)
		return nil
	}
	s := string(content)
	s = strings.ReplaceAll(s, fmt.Sprintf("/appearance/themes/%s/", theme), "")
	os.WriteFile(filename, []byte(s), 0644)

	filename = filepath.Join(dir, fmt.Sprintf("appearance/themes/%s/style/custom/link-icon.css", theme))

	content, err = os.ReadFile(filename)
	if err != nil {
		log.Error(err)
		return nil
	}
	s = string(content)
	s = strings.ReplaceAll(s, fmt.Sprintf("/appearance/themes/%s/", theme), "../../")
	os.WriteFile(filename, []byte(s), 0644)

	return nil
}

// 描述: 基于保存目录创建appid和docid文件夹
// 返回: 返回创建的文件夹路径
// 返回: 错误
func mkdir_all(app_id, doc_id string) (string, error) {
	// 创建目录
	f := filepath.Join(app.Basic.SavePath, app_id, doc_id)
	err := os.MkdirAll(f, get_folder_permission())
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
func rmdir_all(app_id, doc_id string) error {
	// 删除目录
	f := filepath.Join(app.Basic.SavePath, app_id, doc_id)
	err := os.RemoveAll(f)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func get_file_permission() fs.FileMode {
	if app.docker {
		return 0666
	} else {
		return 0660
	}
}
func get_folder_permission() fs.FileMode {
	if app.docker {
		return 0777
	} else {
		return 0755
	}
}
