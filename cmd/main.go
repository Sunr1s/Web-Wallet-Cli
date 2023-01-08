package main

import (
	s_user "github.com/Sunr1s/webclient"
	"github.com/Sunr1s/webclient/pkg/handler"
	"github.com/Sunr1s/webclient/pkg/repository"
	"github.com/Sunr1s/webclient/pkg/service"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	if err := initConfig(); err != nil {
		log.Fatalf("error initialization configs: %s", err.Error())
	}
	nodes := viper.GetStringSlice("node.port")

	db, err := repository.NewSqlite3DB(repository.Config{
		Host:   viper.GetString("db.host"),
		DBName: viper.GetString("db.dbname"),
	})

	if err != nil {
		log.Fatal("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos, nodes)
	handlers := handler.NewHandler(services)

	log.Println("Server is running ...")

	srv := new(s_user.Server)
	err = srv.Run(viper.GetString("port"), handlers.InitRouter())
	if err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
