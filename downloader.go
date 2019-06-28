package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const (
	TEMP_ZIP_DIR  = "tmp_zip"
	TEMP_ZIP_NAME = "ken_all.zip"
	KEN_ALL_NAME  = "KEN_ALL.CSV"
)

func cleanup(p string) error {
	// 仮zipを削除
	if err := os.RemoveAll(TEMP_ZIP_DIR); err != nil {
		return err
	}

	// さらば昔のケン
	if err := os.RemoveAll(p); err != nil {
		return err
	}
	return nil
}

func downloadKenAll(url, p string) error {
	// 一応仮zipフォルダと指定されたKEN_ALLのファイルを削除
	if err := cleanup(p); err != nil {
		return err
	}

	// 仮zipの保存先を作成
	if err := os.MkdirAll(TEMP_ZIP_DIR, 0777); err != nil {
		return err
	}

	temp_zip_path := filepath.Join(TEMP_ZIP_DIR, TEMP_ZIP_NAME)
	// zipをダウンロード
	if err := downloadZip(url, temp_zip_path); err != nil {
		return err
	}

	// 解凍して配置
	if err := unzip(temp_zip_path, p); err != nil {
		return err
	}

	// 仮zipを削除
	if err := os.RemoveAll(TEMP_ZIP_DIR); err != nil {
		return err
	}

	return nil
}

func downloadZip(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func unzip(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != KEN_ALL_NAME {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		buf := make([]byte, f.UncompressedSize)
		_, err = io.ReadFull(rc, buf)
		if err != nil {
			return err
		}

		if err = ioutil.WriteFile(dst, buf, f.Mode()); err != nil {
			return err
		}
	}

	return nil
}
