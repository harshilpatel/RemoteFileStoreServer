package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/harshilkumar/cloud-store-server/utils"
	"github.com/spf13/viper"
)

type Args struct {
	A, B int
}

type Arth struct {
	Quo, Rem int
}

func (a *Arth) Mul(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func main() {

	// arth := new(Arth)
	// rpc.Register(arth)
	// rpc.HandleHTTP()

	viper.SetConfigName("server_config.json")
	viper.SetConfigType("json")

	c := utils.GetConfiguration()
	viper.AddConfigPath(c.ServerBasePath)

	server := utils.CreateServer(c)
	s := utils.Storage{server}
	server.Register(&s)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Could not open the config file %s", err))
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Printf("Shutting down Server")

	viper.Set("server", server)
	viper.WriteConfig()
}
