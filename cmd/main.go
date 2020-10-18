package main

import (
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/routing"
	"github.com/sirupsen/logrus"
)

func main() {
	log := initLogger()
	log.Info("starting photolum")

	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	debug.SetGCPercent(25600)

	plData, err := config.InitPhotolumData(log)
	if err != nil {
		log.WithError(err).Fatal("cannot initialize photolum data")
		os.Exit(1)
	}

	log.Info("starting API server")
	routing.ListenAndServe(plData, log)

	log.Info("blocking until signalled to shutdown")
	// make channel for interrupt signal
	c := make(chan os.Signal, 1)
	// tell os to send to chan when signal received
	signal.Notify(c, os.Interrupt)
	// wait for signal
	<-c

	log.Info("shutting down")
	os.Exit(0)
}

func initLogger() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	switch strings.ToUpper(os.Getenv("PHOTOLUM_LOG_LEVEL")) {
	case "TRACE":
		log.SetLevel(logrus.TraceLevel)
	case "DEBUG":
		log.SetLevel(logrus.DebugLevel)
	case "INFO":
		log.SetLevel(logrus.InfoLevel)
	case "WARN":
		log.SetLevel(logrus.WarnLevel)
	case "ERROR":
		log.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		log.SetLevel(logrus.FatalLevel)
	case "PANIC":
		log.SetLevel(logrus.PanicLevel)
	default:
		log.SetLevel(logrus.WarnLevel)
	}
	return log
}
