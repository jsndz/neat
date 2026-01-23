package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Sha1Hash(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func Compress(blob []byte) ([]byte, error) {
	var buff bytes.Buffer

	zw := zlib.NewWriter(&buff)
	_, err := zw.Write(blob)
	if err != nil {
		return nil, fmt.Errorf("Couldn't compress the file: %v", err)
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func Decompress(compressed []byte) ([]byte, error) {
	zr, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, fmt.Errorf("Error decompressing file: %v", err)
	}
	defer zr.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, zr)
	if err != nil {
		return nil, fmt.Errorf("Error decompressing file: %v", err)
	}
	return out.Bytes(), nil
}

func RepoPath(parts ...string) string {
	all := append([]string{".neat"}, parts...)
	return filepath.Join(all...)
}

func ObjectPath(sha string) (string, string) {
	dir := sha[:2]
	file := sha[2:]
	return file, filepath.Join(RepoPath("objects"), dir)
}

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
func ReadObject(sha string) ([]byte, error) {
	file, objDir := ObjectPath(sha)
	path := filepath.Join(objDir, file)

	compressed, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Decompress(compressed)
}
