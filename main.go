package main

import (
	"log"
	"os"

	"github.com/lazybark/go-jwt/api"
	"github.com/lazybark/go-jwt/config"
	"github.com/lazybark/go-jwt/storage/redis"
	"github.com/lazybark/lazyevent/v2/events"
	"github.com/lazybark/lazyevent/v2/logger"
	"github.com/lazybark/lazyevent/v2/lproc"
)

func main() {
	//Create logger according to LazyEvent
	cli := logger.NewCLI(events.Any)
	plt, err := logger.NewPlaintext("api.log", false, events.Any)
	if err != nil {
		log.Fatal(err)
	}
	//New LogProcessor to rule them all
	p := lproc.New("", make(chan error), false, cli, plt)

	conf, err := config.GetConfig()
	if err != nil {
		p.FatalInCaseErr(events.Error(err.Error()).Red())
		os.Exit(1)
	}

	rdb, err := redis.NewRedisStorage("localhost:6379", "", false, 5, p)
	if err != nil {
		p.FatalInCaseErr(events.Error(err.Error()).Red())
		os.Exit(1)
	}

	server := api.New(rdb, conf, p)

	if conf.FlushDB {
		err = server.StorageFlush()
		if err != nil {
			p.FatalInCaseErr(events.Error(err.Error()).Red())
			os.Exit(1)
		}
	}

	if conf.InitDB {
		err = server.StorgageInit()
		if err != nil {
			p.FatalInCaseErr(events.Error(err.Error()).Red())
			os.Exit(1)
		}
	}

	server.Start()
}
