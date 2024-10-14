package main

import (
	"RescueSupport.sv/config"
	"RescueSupport.sv/serverutil"
)

var c config.Config

func main() {
	c = config.ImportConfig(config.OSSource{})
	store, client := serverutil.SetUpDatabase(c.DatabaseURL, c.DatabaseName)
	//pro := serverutil.SetUpKakifyHandler(c.KafkaBrokers)
	handler := serverutil.SetUpHandler(store)
	companyHandler := serverutil.SetUpCompanyHandler(store)
	server := serverutil.SetUpServer(handler)
	companyServer := serverutil.SetUpCompanyServer(companyHandler)
	router := serverutil.SetupRouter(&server, &companyServer)
	serverutil.StartServer(router, client)
}
