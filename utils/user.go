package utils

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type User struct {
	Username        string
	Key             string
	BaseFilePath    string
	ConfigPath      string
	Objects         map[string]FObject
	ClientInstances map[string]ClientInstance
}

func (u *User) SetUp(c ConfigCloudStore) {
	base_dir := filepath.Join(c.BasePath, u.Username, "files")
	config_path := filepath.Join(c.BasePath, u.Username, "config.txt")
	u.BaseFilePath = base_dir

	logrus.Printf("Setting up Dirs for User: %v path:%v \n", u.Username, u.BaseFilePath)
	if _, e := os.Lstat(base_dir); os.IsNotExist(e) {
		os.MkdirAll(base_dir, 0777)
		os.Create(config_path)
	}
}

func (u *User) CreateUserDir(c ConfigCloudStore) {
	os.Mkdir(filepath.Join(c.BasePath, u.Username, "files"), os.ModeDir)
}

func (u *User) CreateFile(c ConfigCloudStore, filename string) {
	os.Create(filepath.Join(c.BasePath, u.Username, "file", filename))
}

func (u *User) SaveObject(f *FObject, data []byte) {
	to_save_path := filepath.Join(u.BaseFilePath, f.Location)
	to_save_dir := filepath.Dir(to_save_path)
	if _, e := os.Lstat(to_save_dir); os.IsNotExist(e) {
		os.MkdirAll(to_save_dir, 0777)
	}

	os.Create(to_save_path)

	file, err := os.OpenFile(to_save_path, os.O_APPEND, 0755)
	defer file.Close()

	// buf := make([]byte, 1024)
	if err != nil {
		if l, err := file.Write(data); err != nil {
			logrus.Printf("%v Bytes written to file %v", l, f.Name)
		} else {
			logrus.Printf("error in writing to file %v", err)
		}
	}
}
