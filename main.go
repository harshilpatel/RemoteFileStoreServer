package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/harshilkumar/cloud-store-server/utils"
	"github.com/spf13/viper"
)

func main() {

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)

	logrus.Info("Starting Server")

	viper.SetDefault("server.address", "localhost")
	viper.SetDefault("server.port", "4533")
	viper.SetDefault("server.Config.BasePath", "/Users/harshilpatel/Projects/test/cloud-store-server")

	viper.SetConfigName("server_config.json")
	viper.SetConfigType("json")

	config := utils.GetConfiguration()
	viper.AddConfigPath(config.BasePath)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Could not open the config file %s", err))
	}

	storage := utils.LoadOrCreateStorage(config)
	server := utils.CreateServer(config, storage)

	// NOTE: GlobalServer does not server any purpose. This was written to ensure storage receives server. Time constraint :).
	utils.GlobalServer = server

	// sampleUser := utils.User{
	// 	Username: "1234",
	// 	Key:      "1234",
	// }

	// server.RegisterUser(sampleUser)

	// sampleUser1 := utils.User{
	// 	Username: "12345",
	// 	Key:      "12345",
	// }

	// server.RegisterUser(sampleUser1)

	server.Housekeeping()

	server.Register(&storage)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	server.Listen()

	<-sigs // receive shutdown signal

	server.InitiateWatchers()
	server.Watcher.Close()
	server.Watcher = nil

	viper.Set("server", server)
	viper.WriteConfig()

	server.ActiveServer.Close()

	logrus.Printf("Config saved safely")
	logrus.Printf("Shutting down Server")
}
