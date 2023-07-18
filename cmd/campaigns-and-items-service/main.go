package main

import (
	"cais/internal/app"
	"cais/pkg/caching"
)

var opts = app.InstanceOpts{
	Port: ":8080",
	CfgPgVars: []string{
		"PG_HOST",
		"PG_PORT",
		"PG_USER",
		"PG_PASSWD",
		"PG_DBNAME",
	},
	RedisVar: caching.RedisEnvars{
		Addr: "REDIS_ADDR",
		Pass: "REDIS_PASS",
		Db:   "REDIS_DB",
	},
}

func main() {
	app.RunServiceInstance(opts)
}
