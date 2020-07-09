package utils

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type Server struct {
	Address      string
	Port         string
	ActiveServer net.Listener
	// Users   map[string]User
	Config  ConfigCloudStore
	Storage Storage
}

// CreateServer creates a new server
func CreateServer(c ConfigCloudStore, s Storage) Server {

	log.Printf("Creating Server")
	return Server{
		Address: viper.GetString("server.Address"),
		Port:    viper.GetString("server.Port"),
		Config:  c,
		Storage: s,
	}
}

func (s *Server) Register(obj interface{}) {
	if err := rpc.Register(obj); err != nil {
		log.Printf("error registering service %v \n", err)
	}
}

func (s *Server) RegisterUser(u User) {
	s.Storage.Users[u.Username] = u
	viper.WriteConfig()
}

func (s *Server) UnRegisterUser(u User) {
	delete(s.Storage.Users, u.Username)
}

func (s *Server) UpdateHashForObjects() {
	for _, u := range s.Storage.Users {
		log.WithFields(log.Fields{
			"User": u.Username,
		}).Infof("Parsing User for hash user objects")
		for relativePath, obj := range u.Objects {
			obj.UpdateHashForObject(u, s.Config)
			obj.UpdateHashForObjectBlocks(u, s.Config)

			u.Objects[relativePath] = obj
		}
	}

}

func (s *Server) Housekeeping() {

	// NOTE :parses disk and create in memory objects
	if l, e := ioutil.ReadDir(s.Config.BasePath); e == nil {
		for _, item := range l {
			if item.IsDir() {
				log.WithFields(log.Fields{
					"User": item.Name(),
				}).Infof("Found User While parsing Storage")
				if _, ok := s.Storage.Users[item.Name()]; !ok {
					// NOTE :creates a NEW user object
					u := User{
						Username: item.Name(),
						Objects:  make(map[string]FObject),
					}
					userDir := filepath.Join(s.Config.BasePath, item.Name(), "files")
					if _, err := os.Lstat(userDir); os.IsNotExist(err) {
						os.Mkdir(userDir, 0777)
					}

					// NOTE :parses objects on disk
					filepath.Walk(userDir, func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return nil
						}

						if info.IsDir() {
							log.WithFields(log.Fields{
								"DIR": path,
							}).Errorln("Skiping DIR check")
							return nil
						}

						relativePath := strings.TrimPrefix(path, userDir)
						if _, ok := u.Objects[relativePath]; !ok {

							obj := FObject{
								Name:         info.Name(),
								Relativepath: relativePath,
								LastWritten:  info.ModTime().UTC(),
								LastPushed:   info.ModTime().UTC(),
								LastPulled:   info.ModTime().UTC(),
								Size:         info.Size(),
								Version:      0,
							}

							u.Objects[relativePath] = obj
						}

						return nil
					})

					log.Printf("created objects %v \n", len(u.Objects))
					s.Storage.Users[item.Name()] = u
					// u.SetUp(s.Config)
				}
			}
		}
	}

	log.Debugln("Created Users and Objects from DISK")

	// NOTE :Parse Config. imported as raw and type casted as per expectations
	// 		 The purpose to laod all previous KNOWLEDGE on server Restart
	configRUserList := viper.GetStringMap("server.Storage.Users")
	log.Infof("Parsing %v users from config\n", len(configRUserList))
	for configUsername, configUser := range configRUserList {
		log.WithFields(log.Fields{
			"User": configUsername,
		}).Infof("Found User While parsing Config")

		configUserMap := configUser.(map[string]interface{})
		if configUserMap["objects"] != nil {
			configUserObjectsMap := configUserMap["objects"].(map[string]interface{})
			if user, ok := s.Storage.Users[configUsername]; ok {
				user.Key = configUserMap["key"].(string)
				s.Storage.Users[configUsername] = user

				for configFileRelativePath, configObj := range configUserObjectsMap {
					configObjMap := configObj.(map[string]interface{})

					lastWritten := configObjMap["lastwritten"].(string)
					lastPushed := configObjMap["lastpushed"].(string)
					lastPulled := configObjMap["lastpulled"].(string)
					version := configObjMap["version"].(float64)

					if userObject, ok := user.Objects[configFileRelativePath]; ok {
						if LastPulled, err := time.Parse(time.RFC3339Nano, lastPulled); err == nil {
							userObject.LastPulled = LastPulled
						}
						if LastWritten, err := time.Parse(time.RFC3339Nano, lastWritten); err == nil {
							userObject.LastWritten = LastWritten
						}
						if LastPushed, err := time.Parse(time.RFC3339Nano, lastPushed); err == nil {
							userObject.LastPushed = LastPushed
						}

						userObject.Version = int64(version)
						user.Objects[configFileRelativePath] = userObject
					}

				}
			}

		}
	}

	log.Debugln("Created Users and Objects from CONFIG")

	s.UpdateHashForObjects()

}

func (s *Server) Listen() {
	rpc.HandleHTTP()

	log.Printf("Listening at %v:%v", s.Address, s.Port)
	l, e := net.Listen("tcp", s.Address+":"+s.Port)
	s.ActiveServer = l
	if e != nil {
		log.Fatal("listen error")
	}

	go http.Serve(l, nil)
}
