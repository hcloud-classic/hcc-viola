package init

import (
	"hcc/viola/lib/config"
	"hcc/viola/lib/controlcli"
	"hcc/viola/lib/logger"
	"os"
)

// MainInit : Main initialization function
func MainInit() error {
	err := syscheckInit()
	if err != nil {
		return err
	}

	err = loggerInit()
	if err != nil {
		return err
	}

	err = prepareENVSetting()
	if err != nil {
		logger.Logger.Println(err)
		os.Exit(0)
	}
	config.Parser()
	status, telegraferr := controlcli.TelegrafCheck()
	logger.Logger.Println(status, telegraferr)

	err = rabbitmqInit()
	if err != nil {
		return err
	}

	return nil
}
