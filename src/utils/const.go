package utils

type Path map[string][]byte

type TreeEntry struct {
	Name string
	Sha  []byte
	Mode string
}

type Folders map[string][]TreeEntry
