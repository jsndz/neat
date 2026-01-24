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
	"syscall"
)

func Sha1Hash(data []byte) ([]byte, string) {
	h := sha1.New()
	h.Write(data)
	return (h.Sum(nil)), hex.EncodeToString(h.Sum(nil))
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

func GetFileInfo(filename string) error {

	fi, err := os.Stat(filename)
	if err != nil {
		fmt.Println("Write error:", err)
		return err
	}

	// Low level stats
	// copied code from internet
	stat := fi.Sys().(*syscall.Stat_t)

	// --- timestamps ---
	ctimeSec := stat.Ctim.Sec
	ctimeNsec := stat.Ctim.Nsec

	mtimeSec := stat.Mtim.Sec
	mtimeNsec := stat.Mtim.Nsec

	// --- filesystem identity ---
	dev := stat.Dev
	ino := stat.Ino

	// --- permissions / mode ---
	mode := fi.Mode().Perm()

	// Git uses full mode like 100644 / 100755
	var gitMode uint32 = 0100000 | uint32(mode)

	// --- ownership ---
	uid := stat.Uid
	gid := stat.Gid

	// --- size ---
	size := fi.Size()

	// --- flags ---
	nameLen := len(filename)
	stage := 0 // normal files are stage 0
	flags := uint16(nameLen) | uint16(stage<<12)

	// --- output ---
	fmt.Printf("Path:  %s\n", filename)
	fmt.Printf("ctime: %d.%d\n", ctimeSec, ctimeNsec)
	fmt.Printf("mtime: %d.%d\n", mtimeSec, mtimeNsec)
	fmt.Printf("dev:   %d\n", dev)
	fmt.Printf("ino:   %d\n", ino)
	fmt.Printf("mode:  %o\n", gitMode)
	fmt.Printf("uid:   %d\n", uid)
	fmt.Printf("gid:   %d\n", gid)
	fmt.Printf("size:  %d\n", size)
	fmt.Printf("flags: %d\n", flags)

	return nil

}
