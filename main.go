package main

import (
	"hcc/viola/lib/logger"
)

func main() {
	if !syscheck.CheckRoot() {
		return
	}

	if !logger.Prepare() {
		return
	}
	defer logger.FpLog.Close()

	// err := mysql.Prepare()
	// if err != nil {
	// 	return
	// }
	// defer mysql.Db.Close()

	// http.Handle("/graphql", graphql.GraphqlHandler)

	// logger.Logger.Println("Server is running on port " + config.HTTPPort)
	// err = http.ListenAndServe(":"+config.HTTPPort, nil)
	// if err != nil {
	// 	logger.Logger.Println("Failed to prepare http server!")
	// }




	
}
