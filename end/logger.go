package end

import "hcc/viola/lib/logger"

func loggerEnd() {
	_ = logger.FpLog.Close()
}
