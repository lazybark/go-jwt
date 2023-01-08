package config

import "github.com/alexflint/go-arg"

type Config struct {
	Host    string `arg:"--host, env:HOST" help:"Host app listens to" default:":80"`
	DBName  string `arg:"--db, env:DBNAME" help:"Database name" default:"auth_data.db"`
	FlushDB bool   `arg:"--flush, env:FLUSHDB" help:"Flush database" default:"false"`
	InitDB  bool   `arg:"--init, env:INITDB" help:"Init database" default:"false"`
	Secret  string `arg:"--secret, env:AUTH_SECRET" help:"Secret to auth JWT" default:"sample_secret_this_is"`
}

func GetConfig() (Config, error) {
	var c Config
	return c, arg.Parse(&c)
}
