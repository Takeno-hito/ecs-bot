package main

import (
	"github.com/Takeno-hito/ecs-bot/bot"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
		PadLevelText:    true,
	})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)

	if err := godotenv.Load(); err != nil {
		logrus.Warn("cannot load env file")
	}

	bot.Init()

	logrus.Info("application has started successfully")

	bot.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-quit

	logrus.Warn("Shutting down...")

	bot.Close()
	return
}
