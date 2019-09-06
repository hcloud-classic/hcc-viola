package main

import (
	"GraphQL_viola/violacheckroot"
	"GraphQL_viola/violaconfig"
	"GraphQL_viola/violagraphql"
	"GraphQL_viola/violalogger"
	"GraphQL_viola/violamysql"
	"net/http"
)

func main() {
	if !violacheckroot.CheckRoot() {
		return
	}

	if !violalogger.Prepare() {
		return
	}
	defer violalogger.FpLog.Close()

	err := violamysql.Prepare()
	if err != nil {
		return
	}
	defer violamysql.Db.Close()

	http.Handle("/graphql", violagraphql.GraphqlHandler)

	violalogger.Logger.Println("Server is running on port " + violaconfig.HTTPPort)
	err = http.ListenAndServe(":"+violaconfig.HTTPPort, nil)
	if err != nil {
		violalogger.Logger.Println("Failed to prepare http server!")
	}
}
