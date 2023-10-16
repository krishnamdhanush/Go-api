package main

import (
	"github.com/krishnamdhanush/web-dawn/configs"
	"github.com/krishnamdhanush/web-dawn/router"
)

func main() {
	//connect to db
	configs.ConnectDB()
	configs.CreateRedisClient()

	_, err := router.Run()
	if err != nil {
		panic(err)
	}
}
