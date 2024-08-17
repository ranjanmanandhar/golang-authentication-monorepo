package logger

import (
	"flag"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func Logger() log.Logger {

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	flag.Parse()

	logger = log.With(logger,
		"service", "nettv-auth-consumer",
		"time:", log.Timestamp(func() time.Time {
			loc, err := time.LoadLocation("Asia/Kathmandu")
			if err != nil {
				return time.Now()
			}

			return time.Now().In(loc)
		}),
		"caller", log.DefaultCaller,
	)

	defer level.Info(logger).Log("msg", "service ended")

	return logger

}
