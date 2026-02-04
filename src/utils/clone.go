package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Download(url, dest string) {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Couldn't get the code")
		return
	}

	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		fmt.Println("couldn't write to destination")
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("couldn't copy to destination")
	}
}

func Unzip(src, dest string) {
	r, err := zip.OpenReader(src)
	if err != nil {
		fmt.Println("couldn't unzip")
	}
	for _, f := range r.File {
		fp := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fp, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(fp), 0755)
		rc, _ := f.Open()
		out, _ := os.Create(fp)
		io.Copy(out, rc)

		rc.Close()
		out.Close()
	}
}
