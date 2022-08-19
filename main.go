package main

import (
	"fmt"
	"os"

	"github.com/lazybark/go-helpers/cli/clf"
	"github.com/lazybark/go-jwt/api"
	"github.com/lazybark/go-jwt/config"
	"github.com/lazybark/go-jwt/storage/redis"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		fmt.Println(clf.Red(err))
		os.Exit(1)
	}

	rdb, err := redis.NewRedisStorage("localhost:6379", "", false, 5)
	if err != nil {
		fmt.Println(clf.Red(err))
		os.Exit(1)
	}

	server := api.New(rdb, conf)
	if conf.FlushDB {
		err = server.StorageFlush()
		if err != nil {
			fmt.Println(clf.Red(err))
			os.Exit(1)
		}
	}

	if conf.InitDB {
		err = server.StorgageInit()
		if err != nil {
			fmt.Println(clf.Red(err))
			os.Exit(1)
		}
	}

	server.Start()
}
