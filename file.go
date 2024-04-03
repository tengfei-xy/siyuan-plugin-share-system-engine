package main

import (
	"archive/zip"
	"io"
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
		err = os.MkdirAll(filepath.Dir(filename), 0755)
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

func check_theme_file(f string) error {
	filename := filepath.Join(f, "appearance/themes/Odyssey/theme.css")
	if tools.FileExist(filename) {
		log.Info("修改主题 Odyssey")
		return update_theme_file(filename)
	}
	filename = filepath.Join(f, "appearance/themes/Savor/theme.css")
	if tools.FileExist(filename) {
		log.Info("修改主题 Savor")
		return update_theme_file(filename)
	}
	return nil
}
func update_theme_file(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Error(err)
		return nil
	}
	s := string(content)
	s = strings.ReplaceAll(s, "/appearance/themes/Odyssey/", "")
	os.WriteFile(filename, []byte(s), 0644)
	return nil
}
