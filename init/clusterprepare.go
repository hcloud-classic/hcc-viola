package init

import (
	"errors"
	"hcc/viola/lib/controlcli"
)

func prepareENVSetting() error {
	if !controlcli.VerifyClusrter() {
		controlcli.WriteNFSConfigObject()
		controlcli.WriteHostsConfigObject()
		return errors.New("Please Run In Kerrighed\n But Prepare to Cluster environment Success\nHasta La Vista Baby!!")
	}
	return nil
}
