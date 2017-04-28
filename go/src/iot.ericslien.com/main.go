package main

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"

	"iot.ericslien.com/config"
	"iot.ericslien.com/database"
	"iot.ericslien.com/logging"
	"iot.ericslien.com/model"
)

func main() {
	appconf := config.Get()
	logging.SetupLogging(appconf.Log.Name, appconf.Log.Level)
	initDB(appconf)
	server := gin.Default()

	model.SetupRoutes(server)

	server.Run(fmt.Sprintf("0.0.0.0:%d", appconf.Port))
}

var migrate sync.Once

func initDB(conf config.Conf) {
	database.SetConfig(conf.Sqldb)
	var err error
	migrate.Do(func() {
		d := model.Device{}
		err = d.Migrate()
		if err != nil {
			panic(fmt.Sprintf("Could not init or migrate database because: %s", err))
		}
		u := model.User{}
		err = u.Migrate()
		if err != nil {
			panic(fmt.Sprintf("Could not init or migrate database because: %s", err))
		}
		fu := model.FirmwareUpdate{}
		err = fu.Migrate()
		if err != nil {
			panic(fmt.Sprintf("Could not init or migrate database because: %s", err))
		}
	})
}
