package main

import (
	"chat/conf"
	"chat/router"
	"chat/service"
	"fmt"
)

func main()  {
	conf.Init()

	fmt.Println("tttt", conf.MongoDBClient)
	go service.Manager.Start()

	r := router.NewRouter()
	_ = r.Run(conf.HttpPort)
}
