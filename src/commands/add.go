package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/jsndz/neat/src/utils"
)

func Add(filename string) {
	entries := make(utils.Path)
	indexContent, err := os.ReadFile(".neat/index")
	if err == nil {
		entries = ReadIndex(indexContent)
	}
	rel, _ := filepath.Rel(".", filename)
	filename = rel

	blob := utils.CreateBlob(filename)
	shaBin, sha := utils.Sha1Hash(blob)
	utils.WriteToObjects(sha, blob)
	fmt.Println("Added object:", sha)
	entries[filename] = shaBin
	WriteIndex(entries)

}

func WriteIndex(files utils.Path) {
	var buff bytes.Buffer

	buff.Write(IndexHeader(len(files)))
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		sha := files[path]
		entry, err := entryForFile(path, sha)
		if err != nil {
			fmt.Println(err)
			return
		}
		buff.Write(entry)
	}

	if err := utils.EnsureDir(".neat"); err != nil {
		fmt.Println("Initialize neat:", err)
		return
	}
	sum, _ := utils.Sha1Hash(buff.Bytes())
	buff.Write(sum)
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

	// flag is the len of file name and has 2 bytes = 16 bits
	binary.Write(&buff, binary.BigEndian, flags)

	// ---- path ----
	buff.Write([]byte(filename))

	// ---- padding ----
	// adding a few \x00 (zero) bytes after the path
	// so the total size of one entry becomes a multiple of 8 bytes.
	entryLen := 62 + len(filename)
	padding := (8 - (entryLen % 8)) % 8
	buff.Write(make([]byte, padding))

	return buff.Bytes(), nil
}

func AddAll() {
	indexContent, err := os.ReadFile(".neat/index")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	filesInIndex := ReadIndex(indexContent)
	filesToIndex := filesInIndex
	filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if strings.HasPrefix(path, ".neat") {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(".", path)
		path = rel
		blob := utils.CreateBlob(path)
		shaB, sha := utils.Sha1Hash(blob)

		shaExisting, exist := filesToIndex[path]

		if !exist || !bytes.Equal(shaExisting, shaB) {
			utils.WriteToObjects(sha, blob)
			filesToIndex[path] = shaB
		} else {
			filesToIndex[path] = shaB
		}
		return err
	})
	for p := range filesToIndex {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			delete(filesToIndex, p)
		}
	}
	WriteIndex(filesToIndex)
}

func ReadIndex(indexContent []byte) utils.Path {
	// this func gives the path->sha
	// for this you need the exact bits of flag(filename length)

	filesInIndex := make(utils.Path)
	countOfEntries := binary.BigEndian.Uint32(indexContent[8:12])
	offset := 12
	for i := 0; i < int(countOfEntries); i++ {
		entryStart := offset

		sha := make([]byte, 20)
		copy(sha, indexContent[entryStart+40:entryStart+60])

		flag := binary.BigEndian.Uint16(indexContent[entryStart+60 : entryStart+62])
		pathLen := int(flag & 0x0FFF)
		path := indexContent[entryStart+62 : entryStart+62+pathLen]
		filesInIndex[string(path)] = sha
		entryLen := 62 + pathLen
		offset = entryStart + entryLen + (8-(entryLen%8))%8
	}
	return filesInIndex

}
