// @title           Subscriptions API
// @version         1.0
// @description     REST-сервис для агрегации данных об онлайн подписках пользователей
// @host            localhost:8080
// @BasePath        /
package main

import (
	"log/slog"
	"os"
	"testTask/internal/config"
	"testTask/internal/database"
	"testTask/internal/handlers"
	"testTask/internal/repository"
	"testTask/internal/server"
)

func main() {
	Conf, err := config.Load()
	if err != nil {
		slog.Error("unable to load config:", "err", err)
		os.Exit(1)
	}

	db, err := database.Connect(Conf.DBConn)
	if err != nil {
		slog.Error("unable to connect to database: ", "err", err)
		os.Exit(1)
	}

	repo := repository.New(db)
	h := handlers.New(repo)

	srv := server.New(Conf.Port, h)
	if err := srv.Start(); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}

}
