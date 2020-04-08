package api

import (
	"fmt"
	"go-blog-api/config"
	"log"
	"net/http"
)

func Start() error {
	appConfig := config.LoadConfig()

	api := NewAPI(appConfig)
	api.initialize()

	//db, err := services.NewDB(appConfig)
	//if err != nil {
	//	return err
	//}
	//
	//err = db.CreateTables(dal.User{})
	//if err != nil {
	//	return err
	//}
	//
	//defer db.Close()

	log.Printf("Server started at %s...", appConfig.AppHost)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", appConfig.AppHost), api.r); err != nil {
		return err
	}

	return nil
}
