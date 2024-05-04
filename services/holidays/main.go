package main

import (
	"flag"
	"log/slog"

	"github.com/usk81/ashihara/services/holidays/interface/transport/http/server"
	"github.com/usk81/ashihara/shared/utils/logger"
)

func main() {

	var listenAddr string
	flag.StringVar(&listenAddr, "listen-addr", ":80", "server listen address")
	flag.Parse()

	// logLv := os.Getenv("LOG_LEVEL")
	// if logLv == "" {
	// 	logLv = "INFO"
	// }

	conf := server.Config{
		Address: listenAddr,
		Logger:  logger.NewWithLevel(slog.LevelDebug),
	}

	server.Run(conf)
}
