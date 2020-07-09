package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/google/uuid"
)

type UserRequestPackage struct {
	ClientUser User
	Obj        FObject
	Operation  string
	Data       []byte
}

type ClientInstance struct {
	ClientCode string
	UserName   string
}

type Storage struct {
	Users  map[string]User
	Config ConfigCloudStore
}

func LoadOrCreateStorage(c ConfigCloudStore) Storage {
	s := Storage{
		Users:  make(map[string]User),
		Config: c,
	}
	return s
}

func (s *Storage) SaveObject(r *UserRequestPackage, data *[]byte) error {
	logrus.Printf("Received request SaveObject %v %v\n", r.ClientUser.Username, r.Obj.Relativepath)
	if u, ok := s.Users[r.ClientUser.Username]; ok {
		if u.Key != r.ClientUser.Key {
			return errors.New("key doesn't match")
		}

		realPath := filepath.Join(s.Config.BasePath, u.Username, "files", r.Obj.Relativepath)
		os.MkdirAll(filepath.Dir(realPath), 0777)

		switch r.Operation {
		case "Append":
			file, err := os.OpenFile(realPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
			defer file.Close()

			if err == nil {
				defer file.Close()

				if _, e := file.Write(r.Data); e != nil {
					return errors.New(fmt.Sprintf("Could not save data %v", e))
				}
				// logrus.Printf("created new file")
				return nil
			}
		case "Create":
			file, err := os.OpenFile(realPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
			defer file.Close()

			if err == nil {
				defer file.Close()

				if _, e := file.Write(r.Data); e != nil {
					return errors.New(fmt.Sprintf("Could not save data %v", e))
				}
				// logrus.Debugln("created new file")
				return nil
			}
		}
	} else {
		return errors.New("User Not Found")
	}

	return nil

}

func (s *Storage) DownloadObject(inputPack UserRequestPackage, pack *UserRequestPackage) error {
	logrus.WithFields(logrus.Fields{
		"File": inputPack.Obj.Relativepath,
		"Name": inputPack.Obj.Name,
		"User": inputPack.ClientUser.Username,
	}).Debugf("Received request DownloadObject")

	if localUser, ok := s.Users[inputPack.ClientUser.Username]; ok {
		if localUser.Key == inputPack.ClientUser.Key {
			obj := localUser.Objects[inputPack.Obj.Relativepath]
			realPath := obj.GetRealPath(localUser, s.Config)

			if data, err := ioutil.ReadFile(realPath); err == nil {
				pack.Data = data
				pack.Obj.Version = obj.Version
				return nil
			} else {
				logrus.Errorln(err)
			}
		}
	}

	return nil
}

func (s *Storage) VerifyObject(pack *UserRequestPackage, response *int) error {
	logrus.WithFields(logrus.Fields{
		"File": pack.Obj.Relativepath,
		"User": pack.ClientUser.Username,
	}).Debugf("Received request VerifyObject")

	if user, ok := s.Users[pack.ClientUser.Username]; ok {
		if serverObject, okay := user.Objects[pack.Obj.Relativepath]; okay {
			*response = 0
			if !bytes.Equal(serverObject.HashOfFile, pack.Obj.HashOfFile) || serverObject.Version != pack.Obj.Version {
				logrus.WithFields(logrus.Fields{"Reason": "Hash Mismatch or Version Mismatch"}).Infof("Request to Alter")
				*response = 2
				if serverObject.LastWritten.After(pack.Obj.LastWritten) || serverObject.Version > pack.Obj.Version {
					*response = 1
					logrus.WithFields(logrus.Fields{"Reason": "Local Copy is New"}).Infof("Request to Push to Client")
				}
			}
		} else {
			*response = 2
			logrus.WithFields(logrus.Fields{
				"Reason": "Local Copy unavailable",
				"File":   pack.Obj.Relativepath,
			}).Infof("Request to Pull from Client")
		}

		return nil
	}

	return errors.New("User Not Found")
}

func (s *Storage) VerifyUser(u *string, p *string) error {
	logrus.WithFields(logrus.Fields{
		"User": *u,
	}).Debugf("Received request VerifyUser")

	if user, ok := s.Users[*u]; ok {
		if user.Key == *p {
			return nil
		}
		return errors.New("Key not Matched")
	}

	return errors.New("User not Matched")
}

func (s *Storage) RegisterUser(username string, p *string) error {

	if _, ok := s.Users[username]; ok {
		logrus.WithField("Request", "Register Request").Error("Already Exists")
		return errors.New("User Already Exists")
	}

	user := User{
		Username: username,
		Key:      uuid.New().String(),
	}

	*p = user.Key

	user.SetUp(s.Config)
	s.Users[username] = user
	viper.WriteConfig()

	return nil
}
