package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type FObject struct {
	Name         string
	Location     string
	Relativepath string
	IsDir        bool
	IsBinary     bool
	LastWritten  time.Time
	LastPulled   time.Time
	LastPushed   time.Time
	Size         int64
	RequiresPush bool
	RequiresPull bool
	HashOfFile   []byte
	Hash         [][]byte
	Version      int64
}

func (f *FObject) GetRealPath(u User, c ConfigCloudStore) string {
	return filepath.Join(c.BasePath, u.Username, "files", f.Relativepath)
}

func (f *FObject) CreateHashForObject(u User, c ConfigCloudStore) ([]byte, error) {
	realPath := f.GetRealPath(u, c)
	if data, err := ioutil.ReadFile(realPath); err == nil {
		h := sha256.New()
		h.Write(data)
		return h.Sum(nil), nil
	}

	return nil, errors.New("Could not find the file " + realPath)
}

func (f *FObject) CreateHashForObjectBlocks(u User, c ConfigCloudStore) ([][]byte, error) {
	realPath := f.GetRealPath(u, c)

	hash := make([][]byte, 0)
	buf := make([]byte, 500)
	if file, err := os.Open(realPath); err == nil {
		defer file.Close()
		for {
			if n, e := file.Read(buf); e == nil {
				if n > 0 {
					h := sha256.New()
					h.Write(buf)
					res := h.Sum(nil)
					hash = append(hash, res)
				}
			} else if e == io.EOF {
				return hash, nil
			}
		}

	}

	return nil, errors.New("Something went wrong")
}

func (f *FObject) UpdateHashForObject(u User, c ConfigCloudStore) {
	if h, err := f.CreateHashForObject(u, c); err == nil {
		f.HashOfFile = h
	} else {
		fmt.Println(err)
	}
}

func (f *FObject) UpdateHashForObjectBlocks(u User, c ConfigCloudStore) {
	if h, err := f.CreateHashForObjectBlocks(u, c); err == nil {
		f.Hash = h
	} else {
		fmt.Println(err)
	}
}
