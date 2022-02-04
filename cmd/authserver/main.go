package main

import (
	"github.com/morozvol/AuthService/internal/app/apiserver"
	"github.com/morozvol/AuthService/internal/app/config"
	"github.com/morozvol/AuthService/internal/app/store/sqlstore"
	"github.com/morozvol/AuthService/internal/app/store/sqlstore/db"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	container.Provide(config.Init)
	container.Provide(logrus.StandardLogger)
	container.Provide(db.New)
	container.Provide(sqlstore.New)
	container.Provide(apiserver.New)
	return container
}

func main() {

	container := BuildContainer()
	err := container.Invoke(func(config *config.Config, srv *apiserver.Server) {
		apiserver.Start(config, srv)
	})

	if err != nil {
		panic(err)
	}
}
