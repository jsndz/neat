package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/jsndz/neat/src/utils"
)

func Add(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	header := []byte(fmt.Sprintf("blob %d\x00", len(content)))
	blob := append(header, content...)

	shaBin, sha := utils.Sha1Hash(blob)

	file, objDir := utils.ObjectPath(sha)
	if err := utils.EnsureDir(objDir); err != nil {
		fmt.Println("Error creating object dir:", err)
		return
	}

	blobPath := filepath.Join(objDir, file)

	if _, err := os.Stat(blobPath); err == nil {
		fmt.Println("Object already exists:", sha)
		return
	}

	compressed, err := utils.Compress(blob)
	if err != nil {
		fmt.Println("Compression error:", err)
		return
	}

	if err := os.WriteFile(blobPath, compressed, 0644); err != nil {
		fmt.Println("Write error:", err)
		return
	}

	fmt.Println("Added object:", sha)
	Index(filename, shaBin)
	fmt.Println("Index Created")

}

func Index(filename string, sha []byte) {
	var buff bytes.Buffer

	buff.Write(IndexHeader(1))
	entry, err := entryForFile(filename, sha)
	if err != nil {
		fmt.Println(err)
		return
	}
	buff.Write(entry)
	if err := utils.EnsureDir(".neat"); err != nil {
		fmt.Println("Initialize neat:", err)
		return
	}
	if err := os.WriteFile(".neat/index", buff.Bytes(), 0644); err != nil {
		fmt.Println("Write error:", err)
		return
	}
}

func IndexHeader(files int) []byte {

	// Git header format is,
	// 4 bytes for each so 4 bytes  =  uint32 -> 2^4 -> 4bytes

	var buff bytes.Buffer
	buff.Write([]byte("DIRC"))
	// the bytes are written in BigEndian (left->right / Most  significant bit at the beginning)
	// Endian is the way you write bytes into address 00 00 00 10 can be written to memory as 00 00 00 10 or 10 00 00 00
	binary.Write(&buff, binary.BigEndian, uint32(2))
	binary.Write(&buff, binary.BigEndian, uint32(files))
	// binary is package to write and read from buffer/address in binary
	return buff.Bytes()
}

func entryForFile(filename string, sha []byte) ([]byte, error) {
	// entry file consist of 62 bytes of file related data and file path with extra padding
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	stat := fi.Sys().(*syscall.Stat_t)

	// ---- times (must be uint32) ----
	ctimeSec := uint32(stat.Ctim.Sec)
	ctimeNsec := uint32(stat.Ctim.Nsec)
	mtimeSec := uint32(stat.Mtim.Sec)
	mtimeNsec := uint32(stat.Mtim.Nsec)

	dev := uint32(stat.Dev)
	ino := uint32(stat.Ino)

	// ---- git mode ----
	var gitMode uint32 = 0100644
	if fi.Mode()&0111 != 0 {
		gitMode = 0100755
	}

	uid := uint32(stat.Uid)
	gid := uint32(stat.Gid)
	size := uint32(fi.Size())

	// ---- flags ----
	flags := uint16(len(filename))

	var buff bytes.Buffer

	// fixed 62 bytes
	// every number fixed 4 bytes and sha of 20 bytes
	binary.Write(&buff, binary.BigEndian, ctimeSec)
	binary.Write(&buff, binary.BigEndian, ctimeNsec)
	binary.Write(&buff, binary.BigEndian, mtimeSec)
	binary.Write(&buff, binary.BigEndian, mtimeNsec)
	binary.Write(&buff, binary.BigEndian, dev)
	binary.Write(&buff, binary.BigEndian, ino)
	binary.Write(&buff, binary.BigEndian, gitMode)
	binary.Write(&buff, binary.BigEndian, uid)
	binary.Write(&buff, binary.BigEndian, gid)
	binary.Write(&buff, binary.BigEndian, size)

	buff.Write(sha) // raw 20 bytes

	binary.Write(&buff, binary.BigEndian, flags)

	// ---- path ----
	buff.Write([]byte(filename))

	// ---- padding ----
	//adding a few \x00 (zero) bytes after the path
	//so the total size of one entry becomes a multiple of 8 bytes.
	entryLen := 62 + len(filename)
	padding := (8 - (entryLen % 8)) % 8
	buff.Write(make([]byte, padding))

	return buff.Bytes(), nil
}
