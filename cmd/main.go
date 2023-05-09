package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"wallet-app/pkg/handler"
	wallet "wallet-app/pkg/http_server"
	"wallet-app/pkg/repository"
	"wallet-app/pkg/service"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logger.Errorf("error initializating configs %s", err.Error())
		return
	}

	if err := godotenv.Load(); err != nil {
		logger.Errorf("error loading env variables %s", err.Error())
		return
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logger.Errorf("fail to initialize db: %s", err.Error())
		return
	}

	repos := repository.NewRepository(db, logger)
	services := service.NewService(repos, logger)
	handlers := handler.NewHandler(services, logger)

	srv := new(wallet.Server)

	if err := srv.Run(viper.GetString("port"), handlers.InitRouters()); err != nil {
		logger.Errorf("error occured while running http http_server: %s", err.Error())
		return
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
