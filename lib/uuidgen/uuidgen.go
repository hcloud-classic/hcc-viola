package uuidgen

import (
	uuid "github.com/nu7hatch/gouuid"
	"hcc/viola/lib/logger"
)

// Uuidgen : Generate uuid
func Uuidgen() (string, error) {
	out, err := uuid.NewV4()
	if err != nil {
		logger.Logger.Println(err)
		return "", err
	}

	return out.String(), nil
}
