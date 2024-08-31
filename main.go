package main

import (
	"RescueSupport.sv/config"
	"RescueSupport.sv/serverutil"
)

var c config.Config

func main() {
	c = config.ImportConfig(config.OSSource{})
	store, client := serverutil.SetUpDatabase(c.DatabaseURL, c.DatabaseName)
	pro := serverutil.SetUpKakifyHandler(c.KafkaBrokers)
	handler := serverutil.SetUpHandler(store, pro)
	server := serverutil.SetUpServer(handler)
	router := serverutil.SetupRouter(&server)
	serverutil.StartServer(router, client)
}
