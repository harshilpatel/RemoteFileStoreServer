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

	viper.SetConfigName("server_config.json")
	viper.SetConfigType("json")

	config := utils.GetConfiguration()
	viper.AddConfigPath(config.ServerBasePath)

	storage := utils.LoadOrCreateStorage()
	server := utils.CreateServer(config, storage)

	sampleUser := utils.User{
		Username: "1234",
		Key:      "1234",
	}

	server.RegisterUser(sampleUser)
	server.Register(&storage)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Could not open the config file %s", err))
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	server.Listen()

	<-sigs // receive shutdown signal
	log.Printf("Shutting down Server")

	viper.Set("server", server)
	viper.WriteConfig()
}
