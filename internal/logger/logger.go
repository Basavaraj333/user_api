package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Init(dev bool) {
	if dev {
		Log = zap.Must(zap.NewDevelopment())
	} else {
		Log = zap.Must(zap.NewProduction())
	}
}

func Sync() { _ = Log.Sync() }
