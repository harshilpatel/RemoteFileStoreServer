package utils

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/spf13/viper"
)

type Server struct {
	Address string
	Port    string
	Users   map[string]User
	Config  ConfigCloudStore
}

// CreateServer creates a new server
func CreateServer(c ConfigCloudStore) Server {

	log.Printf("Creating Server at localhost:1234")
	return Server{
		Address: viper.GetString("server.Address"),
		Port:    viper.GetString("server.Port"),
		Users:   make(map[string]User),
		Config:  c,
	}
}

func (s *Server) Register(obj interface{}) {
	if err := rpc.Register(obj); err != nil {
		log.Printf("error registering service %v \n", err)
	}

}

func (s *Server) RegisterUser(u User) {
	s.Users[u.Username] = u
}

func (s *Server) UnRegisterUser(u User) {
	delete(s.Users, u.Username)
}

func (s *Server) Listen() {
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", "localhost:1234")
	if e != nil {
		log.Fatal("listen error")
	}

	http.Serve(l, nil)
}
