package main

import "cais/internal/app"

var opts = app.InstanceOpts{
	Port:      ":8080",
	CfgPgVars: []string{"HOST", "PORT", "USER", "PASSWD", "DBNAME"},
}

func main() {
	app.RunServiceInstance(opts)
}
